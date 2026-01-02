package notification

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	authutil "github.com/yorukot/kymarium/utils/auth"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

// GetNotification godoc
// @Summary Get a notification
// @Description Retrieves a notification for a team the user belongs to
// @Tags notifications
// @Produce json
// @Param teamID path string true "Team ID"
// @Param id path string true "Notification ID"
// @Success 200 {object} response.SuccessResponse "Notification retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid team ID or notification ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Notification not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/notifications/{id} [get]
func (h *Handler) GetNotification(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	notificationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid notification ID")
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
		return echo.NewHTTPError(http.StatusNotFound, "Notification not found")
	}

	notification, err := h.Repo.GetNotificationByID(c.Request().Context(), tx, teamID, notificationID)
	if err != nil {
		zap.L().Error("Failed to get notification", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get notification")
	}

	if notification == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Notification not found")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Notification retrieved successfully", notification))
}
