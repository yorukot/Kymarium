package monitor

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

// ListMonitors godoc
// @Summary List monitors
// @Description Lists monitors for a team the user belongs to
// @Tags monitors
// @Produce json
// @Param teamID path string true "Team ID"
// @Success 200 {object} response.SuccessResponse "Monitors retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid team ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Team not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/monitors [get]
func (h *MonitorHandler) ListMonitors(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
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
	defer h.Repo.DeferRollback(tx, c.Request().Context())

	member, err := h.Repo.GetTeamMemberByUserID(c.Request().Context(), tx, teamID, *userID)
	if err != nil {
		zap.L().Error("Failed to get team membership", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get team membership")
	}

	if member == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Team not found")
	}

	monitors, err := h.Repo.ListMonitorsByTeamID(c.Request().Context(), tx, teamID)
	if err != nil {
		zap.L().Error("Failed to list monitors", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list monitors")
	}

	if err := h.Repo.CommitTransaction(tx, c.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Monitors retrieved successfully", newMonitorResponses(monitors)))
}
