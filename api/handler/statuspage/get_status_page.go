package statuspage

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

// GetStatusPage godoc
// @Summary Get a status page
// @Description Fetches a status page with groups and monitors for a team the user belongs to
// @Tags status_pages
// @Produce json
// @Param teamID path string true "Team ID"
// @Param id path string true "Status Page ID"
// @Success 200 {object} response.SuccessResponse "Status page retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid team ID or status page ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Status page not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/status-pages/{id} [get]
func (h *Handler) GetStatusPage(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	statusPageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid status page ID")
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

	page, err := h.Repo.GetStatusPageByID(c.Request().Context(), tx, teamID, statusPageID)
	if err != nil {
		zap.L().Error("Failed to get status page", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get status page")
	}

	if page == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Status page not found")
	}

	groups, err := h.Repo.ListStatusPageGroupsByStatusPageID(c.Request().Context(), tx, page.ID)
	if err != nil {
		zap.L().Error("Failed to list status page groups", zap.Error(err), zap.Int64("status_page_id", page.ID))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list status page groups")
	}
	if groups == nil {
		groups = []models.StatusPageGroup{}
	}

	monitors, err := h.Repo.ListStatusPageMonitorsByStatusPageID(c.Request().Context(), tx, page.ID)
	if err != nil {
		zap.L().Error("Failed to list status page monitors", zap.Error(err), zap.Int64("status_page_id", page.ID))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list status page monitors")
	}
	if monitors == nil {
		monitors = []models.StatusPageMonitor{}
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	elements := buildStatusPageElementResponses(groups, monitors)

	resp := statusPageResponse{
		StatusPage: *page,
		Elements:   elements,
	}

	return c.JSON(http.StatusOK, response.Success("Status page retrieved successfully", resp))
}
