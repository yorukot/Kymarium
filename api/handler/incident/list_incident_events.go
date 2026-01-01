package incident

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

// ListIncidentEvents godoc
// @Summary List incident events
// @Description Lists events for an incident the user has access to
// @Tags incidents
// @Produce json
// @Param teamID path string true "Team ID"
// @Param incidentID path string true "Incident ID"
// @Success 200 {object} response.SuccessResponse "Incident events retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid IDs"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Incident not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/incidents/{incidentID}/events [get]
func (h *Handler) ListIncidentEvents(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	incidentID, err := strconv.ParseInt(c.Param("incidentID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid incident ID")
	}

	userID, err := authutil.GetUserIDFromContext(c)
	if err != nil {
		zap.L().Error("Failed to parse user ID from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	if userID == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	ctx := c.Request().Context()
	tx, err := h.Repo.StartTransaction(ctx)
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer h.Repo.DeferRollback(ctx, tx)

	member, err := h.Repo.GetTeamMemberByUserID(ctx, tx, teamID, *userID)
	if err != nil {
		zap.L().Error("Failed to get team membership", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get team membership")
	}

	if member == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Incident not found")
	}

	incident, err := h.Repo.GetIncidentByIDForTeam(ctx, tx, teamID, incidentID)
	if err != nil {
		zap.L().Error("Failed to get incident", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get incident")
	}

	if incident == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Incident not found")
	}

	events, err := h.Repo.ListEventTimelinesByIncidentID(ctx, tx, incidentID)
	if err != nil {
		zap.L().Error("Failed to list incident events", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list incident events")
	}

	if err := h.Repo.CommitTransaction(ctx, tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Incident events retrieved successfully", events))
}
