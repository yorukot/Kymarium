package notification

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/repository"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/id"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

type createNotificationRequest struct {
	Type   models.NotificationType `json:"type" validate:"required,oneof=discord telegram email"`
	Name   string                  `json:"name" validate:"required,min=1,max=255"`
	Config json.RawMessage         `json:"config" validate:"required"`
}

// New godoc
// @Summary Create a notification
// @Description Creates a notification channel for the given team (owner/admin only)
// @Tags notifications
// @Accept json
// @Produce json
// @Param teamID path string true "Team ID"
// @Param request body createNotificationRequest true "Notification create request"
// @Success 200 {object} response.SuccessResponse "Notification created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or team ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Team not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/notifications [post]
func (h *NotificationHandler) New(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	var req createNotificationRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if len(req.Config) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Notification config is required")
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
		return echo.NewHTTPError(http.StatusNotFound, "Team not found")
	}

	if member.Role != models.MemberRoleOwner && member.Role != models.MemberRoleAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to create notifications for this team")
	}

	notificationID, err := id.GetID()
	if err != nil {
		zap.L().Error("Failed to generate notification ID", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate notification ID")
	}

	now := time.Now()
	notification := models.Notification{
		ID:        notificationID,
		TeamID:    teamID,
		Type:      req.Type,
		Name:      req.Name,
		Config:    req.Config,
		UpdatedAt: now,
		CreatedAt: now,
	}

	if err := repository.CreateNotification(c.Request().Context(), tx, notification); err != nil {
		zap.L().Error("Failed to create notification", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create notification")
	}

	if err := repository.CommitTransaction(tx, c.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Notification created successfully", notification))
}
