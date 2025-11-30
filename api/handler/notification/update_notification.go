package notification

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/repository"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

type updateNotificationRequest struct {
	Type   *models.NotificationType `json:"type" validate:"omitempty,oneof=discord telegram email"`
	Name   *string                  `json:"name" validate:"omitempty,min=1,max=255"`
	Config *json.RawMessage         `json:"config"`
}

// UpdateNotification godoc
// @Summary Update a notification
// @Description Updates a notification for a team (owner/admin only)
// @Tags notifications
// @Accept json
// @Produce json
// @Param teamID path string true "Team ID"
// @Param id path string true "Notification ID"
// @Param request body updateNotificationRequest true "Notification update request"
// @Success 200 {object} response.SuccessResponse "Notification updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or IDs"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Notification not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/notifications/{id} [patch]
func (h *NotificationHandler) UpdateNotification(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	notificationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid notification ID")
	}

	var req updateNotificationRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Type == nil && req.Name == nil && req.Config == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "At least one field must be provided to update")
	}

	if err := validator.New().Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Config != nil && len(*req.Config) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Notification config cannot be empty")
	}

	userID, err := authutil.GetUserIDFromContext(c)
	if err != nil {
		zap.L().Error("Failed to parse user ID from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	if userID == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	tx, err := repository.StartTransaction(h.DB, c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer repository.DeferRollback(tx, c.Request().Context())

	member, err := repository.GetTeamMemberByUserID(c.Request().Context(), tx, teamID, *userID)
	if err != nil {
		zap.L().Error("Failed to get team membership", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get team membership")
	}

	if member == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Notification not found")
	}

	if member.Role != models.MemberRoleOwner && member.Role != models.MemberRoleAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to update this notification")
	}

	existing, err := repository.GetNotificationByID(c.Request().Context(), tx, teamID, notificationID)
	if err != nil {
		zap.L().Error("Failed to get notification", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get notification")
	}

	if existing == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Notification not found")
	}

	if req.Type != nil {
		existing.Type = *req.Type
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}

	if req.Config != nil {
		existing.Config = *req.Config
	}

	existing.UpdatedAt = time.Now()

	notification, err := repository.UpdateNotification(c.Request().Context(), tx, *existing)
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Notification not found")
		}

		zap.L().Error("Failed to update notification", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update notification")
	}

	if notification == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Notification not found")
	}

	if err := repository.CommitTransaction(tx, c.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Notification updated successfully", notification))
}
