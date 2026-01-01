package monitor

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

// GetMonitor godoc
// @Summary Get a monitor
// @Description Retrieves a monitor for a team the user belongs to
// @Tags monitors
// @Produce json
// @Param teamID path string true "Team ID"
// @Param id path string true "Monitor ID"
// @Success 200 {object} response.SuccessResponse "Monitor retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid team ID or monitor ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Monitor not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/monitors/{id} [get]
func (h *Handler) GetMonitor(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	monitorID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid monitor ID")
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
		return echo.NewHTTPError(http.StatusNotFound, "Monitor not found")
	}

	monitor, err := h.Repo.GetMonitorByID(c.Request().Context(), tx, teamID, monitorID)
	if err != nil {
		zap.L().Error("Failed to get monitor", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get monitor")
	}

	if monitor == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Monitor not found")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Monitor retrieved successfully", newMonitorResponse(*monitor)))
}
