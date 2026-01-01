package statuspage

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/id"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

// CreateStatusPage godoc
// @Summary Create a status page
// @Description Creates a status page with groups and monitors for the given team (owner/admin only)
// @Tags status_pages
// @Accept json
// @Produce json
// @Param teamID path string true "Team ID"
// @Param request body statusPageUpsertRequest true "Status page create request"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /teams/{teamID}/status-pages [post]
func (h *Handler) CreateStatusPage(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	var req statusPageUpsertRequest
	if err := c.Bind(&req); err != nil {
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
		return echo.NewHTTPError(http.StatusNotFound, "Team not found")
	}

	if member.Role != models.MemberRoleOwner && member.Role != models.MemberRoleAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to create status pages for this team")
	}

	// ensure slug uniqueness
	existingSlug, err := h.Repo.GetStatusPageBySlug(c.Request().Context(), tx, normalizedReq.Slug)
	if err != nil {
		zap.L().Error("Failed to check slug uniqueness", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check slug uniqueness")
	}
	if existingSlug != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Slug already exists")
	}

	now := time.Now()
	statusPageID, err := id.GetID()
	if err != nil {
		zap.L().Error("Failed to generate status page ID", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate status page ID")
	}

	page := models.StatusPage{
		ID:        statusPageID,
		TeamID:    teamID,
		Title:     normalizedReq.Title,
		Slug:      normalizedReq.Slug,
		Icon:      normalizedReq.Icon,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.Repo.CreateStatusPage(c.Request().Context(), tx, page); err != nil {
		zap.L().Error("Failed to create status page", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create status page")
	}

	groups, monitors, err := buildStatusPageElements(normalizedReq, page.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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
		StatusPage: page,
		Elements:   elements,
	}

	return c.JSON(http.StatusOK, response.Success("Status page created successfully", resp))
}

// buildStatusPageElements assigns IDs and fan-out relationships.
func buildStatusPageElements(req statusPageUpsertRequest, statusPageID int64) ([]models.StatusPageGroup, []models.StatusPageMonitor, error) {
	groupIDMap := make(map[int64]int64) // client temp ID -> generated ID
	groups := make([]models.StatusPageGroup, 0, len(req.Groups))
	for _, g := range req.Groups {
		gid, err := id.GetID()
		if err != nil {
			return nil, nil, err
		}
		groups = append(groups, models.StatusPageGroup{
			ID:           gid,
			StatusPageID: statusPageID,
			Name:         g.Name,
			Type:         g.Type,
			SortOrder:    g.SortOrder,
		})
		if g.ID != nil {
			groupIDMap[*g.ID] = gid
		}
	}

	monitors := make([]models.StatusPageMonitor, 0, len(req.Monitors))
	for _, m := range req.Monitors {
		mid, err := id.GetID()
		if err != nil {
			return nil, nil, err
		}
		var groupID *int64
		if m.GroupID != nil {
			if mapped, ok := groupIDMap[*m.GroupID]; ok {
				groupID = &mapped
			} else {
				groupID = m.GroupID
			}
		}
		monitors = append(monitors, models.StatusPageMonitor{
			ID:           mid,
			StatusPageID: statusPageID,
			MonitorID:    m.MonitorID,
			GroupID:      groupID,
			Name:         m.Name,
			Type:         m.Type,
			SortOrder:    m.SortOrder,
		})
	}

	return groups, monitors, nil
}
