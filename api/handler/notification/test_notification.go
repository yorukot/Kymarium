package notification

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	notificationcore "github.com/yorukot/kymarium/core/notification"
	"github.com/yorukot/kymarium/models"
	authutil "github.com/yorukot/kymarium/utils/auth"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

// TestNotification godoc
// @Summary Send a test notification
// @Description Sends a test message to verify a notification configuration (owner/admin only)
// @Tags notifications
// @Produce json
// @Param teamID path string true "Team ID"
// @Param id path string true "Notification ID"
// @Success 200 {object} response.SuccessResponse "Test notification sent successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid team ID or notification ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Notification not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/notifications/{id}/test [post]
func (h *Handler) TestNotification(c echo.Context) error {
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

	if member.Role != models.MemberRoleOwner && member.Role != models.MemberRoleAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to send test notifications for this team")
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

	title := "Kymarium notification test"
	description := fmt.Sprintf("Test notification for team %d and channel %q", teamID, notification.Name)

	if err := notificationcore.Send(c.Request().Context(), *notification, title, description, models.PingStatusSuccessful); err != nil {
		zap.L().Error("Failed to send test notification", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send test notification")
	}

	return c.JSON(http.StatusOK, response.SuccessMessage("Test notification sent successfully"))
}
