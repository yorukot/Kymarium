package monitor

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/utils"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

type updateMonitorRequest struct {
	Name              string             `json:"name" validate:"required,min=1,max=255"`
	Type              models.MonitorType `json:"type" validate:"required,oneof=http ping"`
	Interval          int                `json:"interval" validate:"required,min=30,max=2592000"`
	Config            json.RawMessage    `json:"config" validate:"required"`
	FailureThreshold  int16              `json:"failure_threshold" validate:"required,gt=0"`
	RecoveryThreshold int16              `json:"recovery_threshold" validate:"required,gt=0"`
	Regions           regionIDList       `json:"regions" validate:"required,min=1"`
	NotificationIDs   notificationIDList `json:"notification"`
}

// UpdateMonitor godoc
// @Summary Update a monitor
// @Description Updates a monitor for a team (owner/admin only)
// @Tags monitors
// @Accept json
// @Produce json
// @Param teamID path string true "Team ID"
// @Param id path string true "Monitor ID"
// @Param request body updateMonitorRequest true "Monitor update request"
// @Success 200 {object} response.SuccessResponse "Monitor updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or IDs"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Monitor not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/monitors/{id} [put]
func (h *Handler) UpdateMonitor(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	monitorID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid monitor ID")
	}

	var req updateMonitorRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if len(req.Config) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Monitor config is required")
	}

	if len(req.Regions) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "At least one region is required")
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

	if member.Role != models.MemberRoleOwner && member.Role != models.MemberRoleAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to update monitors for this team")
	}

	existing, err := h.Repo.GetMonitorByID(c.Request().Context(), tx, teamID, monitorID)
	if err != nil {
		zap.L().Error("Failed to get monitor", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get monitor")
	}

	if existing == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Monitor not found")
	}

	now := time.Now()
	regionIDs := utils.UniqueInt64s(req.Regions.Int64s())
	regions, err := h.Repo.ListRegionsByIDs(c.Request().Context(), tx, regionIDs)
	if err != nil {
		zap.L().Error("Failed to load regions", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load regions")
	}

	if len(regions) != len(regionIDs) {
		return echo.NewHTTPError(http.StatusBadRequest, "One or more regions do not exist")
	}

	notificationIDs := req.NotificationIDs.Int64s()
	monitor := models.Monitor{
		ID:                monitorID,
		TeamID:            teamID,
		Name:              req.Name,
		Type:              req.Type,
		Status:            existing.Status, // preserve current status when updating config
		Interval:          req.Interval,
		Config:            req.Config,
		LastChecked:       existing.LastChecked,
		NextCheck:         now.Add(time.Duration(req.Interval) * time.Second),
		FailureThreshold:  req.FailureThreshold,
		RecoveryThreshold: req.RecoveryThreshold,
		RegionIDs:         regionIDs,
		NotificationIDs:   notificationIDs,
		UpdatedAt:         now,
		CreatedAt:         existing.CreatedAt,
	}

	updated, err := h.Repo.UpdateMonitor(c.Request().Context(), tx, monitor)
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Monitor not found")
		}

		zap.L().Error("Failed to update monitor", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update monitor")
	}

	if updated == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Monitor not found")
	}

	// Update monitor-notification associations
	// First, delete existing associations
	if err := h.Repo.DeleteMonitorNotifications(c.Request().Context(), tx, monitorID); err != nil {
		zap.L().Error("Failed to delete monitor notifications", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete monitor notifications")
	}

	if err := h.Repo.DeleteMonitorRegions(c.Request().Context(), tx, monitorID); err != nil {
		zap.L().Error("Failed to delete monitor regions", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete monitor regions")
	}

	// Then create new associations
	if len(notificationIDs) > 0 {
		if err := h.Repo.CreateMonitorNotifications(c.Request().Context(), tx, monitorID, notificationIDs); err != nil {
			zap.L().Error("Failed to create monitor notifications", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create monitor notifications")
		}
	}

	if err := h.Repo.CreateMonitorRegions(c.Request().Context(), tx, monitorID, regions); err != nil {
		zap.L().Error("Failed to create monitor regions", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create monitor regions")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	updated.NotificationIDs = notificationIDs
	updated.RegionIDs = regionIDs

	return c.JSON(http.StatusOK, response.Success("Monitor updated successfully", newMonitorResponse(*updated)))
}
