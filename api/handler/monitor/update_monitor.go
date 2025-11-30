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
	"github.com/yorukot/knocker/repository"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

type updateMonitorRequest struct {
	Name            string             `json:"name" validate:"required,min=1,max=255"`
	Type            models.MonitorType `json:"type" validate:"required,oneof=http"`
	Interval        int                `json:"interval" validate:"required,gt=0"`
	Config          json.RawMessage    `json:"config" validate:"required"`
	NotificationIDs []int64            `json:"notification"`
	GroupID         *int64             `json:"group,omitempty"`
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
func (h *MonitorHandler) UpdateMonitor(c echo.Context) error {
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
		return echo.NewHTTPError(http.StatusNotFound, "Monitor not found")
	}

	if member.Role != models.MemberRoleOwner && member.Role != models.MemberRoleAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to update monitors for this team")
	}

	existing, err := repository.GetMonitorByID(c.Request().Context(), tx, teamID, monitorID)
	if err != nil {
		zap.L().Error("Failed to get monitor", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get monitor")
	}

	if existing == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Monitor not found")
	}

	now := time.Now()
	monitor := models.Monitor{
		ID:              monitorID,
		TeamID:          teamID,
		Name:            req.Name,
		Type:            req.Type,
		Interval:        req.Interval,
		Config:          req.Config,
		LastChecked:     existing.LastChecked,
		NextCheck:       now.Add(time.Duration(req.Interval) * time.Second),
		NotificationIDs: req.NotificationIDs,
		UpdatedAt:       now,
		CreatedAt:       existing.CreatedAt,
		GroupID:         req.GroupID,
	}

	updated, err := repository.UpdateMonitor(c.Request().Context(), tx, monitor)
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

	if err := repository.CommitTransaction(tx, c.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Monitor updated successfully", updated))
}
