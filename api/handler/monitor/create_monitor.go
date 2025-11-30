package monitor

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

type createMonitorRequest struct {
	Name            string             `json:"name" validate:"required,min=1,max=255"`
	Type            models.MonitorType `json:"type" validate:"required,oneof=http"`
	Interval        int                `json:"interval" validate:"required,gt=0"`
	Config          json.RawMessage    `json:"config" validate:"required"`
	NotificationIDs []int64            `json:"notification"`
	GroupID         *int64             `json:"group,omitempty"`
}

// CreateMonitor godocit 
// @Summary Create a monitor
// @Description Creates a monitor for the given team (owner/admin only)
// @Tags monitors
// @Accept json
// @Produce json
// @Param teamID path string true "Team ID"
// @Param request body createMonitorRequest true "Monitor create request"
// @Success 200 {object} response.SuccessResponse "Monitor created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or team ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Team not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/monitors [post]
func (h *MonitorHandler) CreateMonitor(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	var req createMonitorRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if len(req.Config) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Monitor config is required")
	}

	if req.NotificationIDs == nil {
		req.NotificationIDs = []int64{}
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
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to create monitors for this team")
	}

	monitorID, err := id.GetID()
	if err != nil {
		zap.L().Error("Failed to generate monitor ID", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate monitor ID")
	}

	now := time.Now()
	monitor := models.Monitor{
		ID:              monitorID,
		TeamID:          teamID,
		Name:            req.Name,
		Type:            req.Type,
		Interval:        req.Interval,
		Config:          req.Config,
		LastChecked:     now,
		NextCheck:       now.Add(time.Duration(req.Interval) * time.Second),
		NotificationIDs: req.NotificationIDs,
		UpdatedAt:       now,
		CreatedAt:       now,
		GroupID:         req.GroupID,
	}

	if err := repository.CreateMonitor(c.Request().Context(), tx, monitor); err != nil {
		zap.L().Error("Failed to create monitor", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create monitor")
	}

	if err := repository.CommitTransaction(tx, c.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Monitor created successfully", monitor))
}
