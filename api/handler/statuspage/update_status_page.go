package statuspage

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/models"
	authutil "github.com/yorukot/kymarium/utils/auth"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

// UpdateStatusPage godoc
// @Summary Update a status page
// @Description Updates slug/icon/groups/monitors for the given status page (owner/admin only)
// @Tags status_pages
// @Accept json
// @Produce json
// @Param teamID path string true "Team ID"
// @Param id path string true "Status Page ID"
// @Param request body statusPageUpsertRequest true "Status page update request"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /teams/{teamID}/status-pages/{id} [put]
func (h *Handler) UpdateStatusPage(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	statusPageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid status page ID")
	}

	req, err := bindStatusPageUpsert(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	normalizedReq, err := normalizeStatusPageUpsert(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := validator.New().Struct(normalizedReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validateStatusPagePayload(normalizedReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, err := authutil.GetUserIDFromContext(c)
	if err != nil {
		zap.L().Error("Failed to parse user ID from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	if userID == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	tx, err := h.Repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer h.Repo.DeferRollback(c.Request().Context(), tx)

	member, err := h.Repo.GetTeamMemberByUserID(c.Request().Context(), tx, teamID, *userID)
	if err != nil {
		zap.L().Error("Failed to get team membership", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get team membership")
	}

	if member == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Status page not found")
	}

	if member.Role != models.MemberRoleOwner && member.Role != models.MemberRoleAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to update status pages for this team")
	}

	existing, err := h.Repo.GetStatusPageByID(c.Request().Context(), tx, teamID, statusPageID)
	if err != nil {
		zap.L().Error("Failed to get status page", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get status page")
	}

	if existing == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Status page not found")
	}

	// slug uniqueness (allow same slug for same record)
	if slugOwner, err := h.Repo.GetStatusPageBySlug(c.Request().Context(), tx, normalizedReq.Slug); err != nil {
		zap.L().Error("Failed to check slug uniqueness", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check slug uniqueness")
	} else if slugOwner != nil && slugOwner.ID != existing.ID {
		return echo.NewHTTPError(http.StatusBadRequest, "Slug already exists")
	}

	now := time.Now()
	updatedPage := models.StatusPage{
		ID:        existing.ID,
		TeamID:    existing.TeamID,
		Title:     normalizedReq.Title,
		Slug:      normalizedReq.Slug,
		Icon:      normalizedReq.Icon,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: now,
	}

	page, err := h.Repo.UpdateStatusPage(c.Request().Context(), tx, updatedPage)
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Status page not found")
		}
		zap.L().Error("Failed to update status page", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update status page")
	}

	groups, monitors, err := buildStatusPageElements(normalizedReq, page.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.Repo.DeleteStatusPageMonitorsByStatusPageID(c.Request().Context(), tx, page.ID); err != nil {
		zap.L().Error("Failed to clear status page monitors", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to clear status page monitors")
	}

	if err := h.Repo.DeleteStatusPageGroupsByStatusPageID(c.Request().Context(), tx, page.ID); err != nil {
		zap.L().Error("Failed to clear status page groups", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to clear status page groups")
	}

	if err := h.Repo.CreateStatusPageGroups(c.Request().Context(), tx, groups); err != nil {
		zap.L().Error("Failed to create status page groups", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create status page groups")
	}

	if err := h.Repo.CreateStatusPageMonitors(c.Request().Context(), tx, monitors); err != nil {
		zap.L().Error("Failed to create status page monitors", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create status page monitors")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	elements := buildStatusPageElementResponses(groups, monitors)

	resp := statusPageResponse{
		StatusPage: *page,
		Elements:   elements,
	}

	return c.JSON(http.StatusOK, response.Success("Status page updated successfully", resp))
}
