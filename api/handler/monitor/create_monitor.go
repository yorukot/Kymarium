package monitor

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/models/monitorm"
	"github.com/yorukot/kymarium/utils"
	authutil "github.com/yorukot/kymarium/utils/auth"
	"github.com/yorukot/kymarium/utils/id"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

type createMonitorRequest struct {
	Name              string             `json:"name" validate:"required,min=1,max=255"`
	Type              models.MonitorType `json:"type" validate:"required,oneof=http ping"`
	Interval          int                `json:"interval" validate:"required,min=30,max=2592000"`
	Config            json.RawMessage    `json:"config" validate:"required"`
	FailureThreshold  int16              `json:"failure_threshold" validate:"required,gt=0"`
	RecoveryThreshold int16              `json:"recovery_threshold" validate:"required,gt=0"`
	Regions           regionIDList       `json:"regions" validate:"required,min=1"`
	NotificationIDs   notificationIDList `json:"notification"`
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
func (h *Handler) CreateMonitor(c echo.Context) error {
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

	if len(req.Regions) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "At least one region is required")
	}

	if err := validateMonitorConfig(c, req.Type, req.Config); err != nil {
		return err
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
		Status:            models.MonitorStatusUp, // newly created monitors start in healthy state
		Interval:          req.Interval,
		Config:            req.Config,
		LastChecked:       now,
		NextCheck:         now.Add(time.Duration(req.Interval) * time.Second),
		FailureThreshold:  req.FailureThreshold,
		RecoveryThreshold: req.RecoveryThreshold,
		RegionIDs:         regionIDs,
		NotificationIDs:   notificationIDs,
		UpdatedAt:         now,
		CreatedAt:         now,
	}

	if err := h.Repo.CreateMonitor(c.Request().Context(), tx, monitor); err != nil {
		zap.L().Error("Failed to create monitor", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create monitor")
	}

	// Create monitor-notification associations
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

	monitor.NotificationIDs = notificationIDs
	monitor.RegionIDs = regionIDs

	return c.JSON(http.StatusOK, response.Success("Monitor created successfully", newMonitorResponse(monitor)))
}

func validateMonitorConfig(c echo.Context, monitorType models.MonitorType, configRaw json.RawMessage) error {
	switch monitorType {
	case models.MonitorTypeHTTP:
		var config monitorm.HTTPMonitorConfig
		if err := json.Unmarshal(configRaw, &config); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid monitor config")
		}
		if err := validator.New().Struct(config); err != nil {
			if errResponse, ok := response.ValidationErrorResponse(err, "Invalid request body", "VALIDATION_ERROR"); ok {
				return c.JSON(http.StatusBadRequest, errResponse)
			}

			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}
	case models.MonitorTypePing:
		var config monitorm.PingMonitorConfig
		if err := json.Unmarshal(configRaw, &config); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid monitor config")
		}
		if err := validator.New().Struct(config); err != nil {
			if errResponse, ok := response.ValidationErrorResponse(err, "Invalid request body", "VALIDATION_ERROR"); ok {
				return c.JSON(http.StatusBadRequest, errResponse)
			}

			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "Unsupported monitor type")
	}

	return nil
}
