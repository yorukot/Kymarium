package incident

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/utils"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/id"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

type monitorIDList = utils.IDList

type createIncidentRequest struct {
	Status     models.IncidentStatus `json:"status" validate:"omitempty,oneof=detected investigating identified monitoring resolved"`
	Severity   models.IncidentSeverity `json:"severity" validate:"omitempty,oneof=emergency critical major minor info"`
	Message    string                `json:"message" validate:"omitempty,min=1"`
	StartedAt  *time.Time            `json:"started_at,omitempty"`
	Public     *bool                 `json:"public"`
	AutoResolve *bool                `json:"auto_resolve"`
	MonitorIDs monitorIDList         `json:"monitor_ids" validate:"required,min=1"`
}

// CreateIncident godoc
// @Summary Create a new incident
// @Description Manually creates an incident for a monitor the user has access to
// @Tags incidents
// @Accept json
// @Produce json
// @Param teamID path string true "Team ID"
// @Param request body createIncidentRequest true "Incident create payload"
// @Success 200 {object} response.SuccessResponse "Incident created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Monitor not found"
// @Failure 409 {object} response.ErrorResponse "An open incident already exists"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/incidents [post]
func (h *IncidentHandler) CreateIncident(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	var req createIncidentRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	userID, err := authutil.GetUserIDFromContext(c)
	if err != nil {
		zap.L().Error("Failed to parse user ID from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	if userID == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	ctx := c.Request().Context()
	tx, err := h.Repo.StartTransaction(ctx)
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer h.Repo.DeferRollback(tx, ctx)

	member, err := h.Repo.GetTeamMemberByUserID(ctx, tx, teamID, *userID)
	if err != nil {
		zap.L().Error("Failed to get team membership", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get team membership")
	}

	if member == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Monitor not found")
	}

	monitorIDs := utils.UniqueInt64s(req.MonitorIDs.Int64s())
	if len(monitorIDs) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "At least one monitor_id is required")
	}

	for _, mid := range monitorIDs {
		monitor, err := h.Repo.GetMonitorByID(ctx, tx, teamID, mid)
		if err != nil {
			zap.L().Error("Failed to get monitor", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get monitor")
		}
		if monitor == nil {
			return echo.NewHTTPError(http.StatusNotFound, "Monitor not found")
		}

		openIncident, err := h.Repo.GetOpenIncidentByMonitorID(ctx, tx, mid)
		if err != nil {
			zap.L().Error("Failed to check open incidents", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check open incidents")
		}

		if openIncident != nil && openIncident.Status != models.IncidentStatusResolved {
			return echo.NewHTTPError(http.StatusConflict, "An open incident already exists for one of the monitors")
		}
	}

	status := req.Status
	if status == "" {
		status = models.IncidentStatusDetected
	}
	severity := req.Severity
	if severity == "" {
		severity = models.IncidentSeverityMajor
	}

	isPublic := true
	if req.Public != nil {
		isPublic = *req.Public
	}

	autoResolve := false
	if req.AutoResolve != nil {
		autoResolve = *req.AutoResolve
	}

	now := time.Now().UTC()
	startedAt := req.StartedAt
	if startedAt == nil {
		startedAt = &now
	}

	var resolvedAt *time.Time
	if status == models.IncidentStatusResolved {
		resolvedAt = &now
	}

	incidentID, err := id.GetID()
	if err != nil {
		zap.L().Error("Failed to generate incident ID", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create incident")
	}

	incident := models.Incident{
		ID:          incidentID,
		Status:      status,
		Severity:    severity,
		IsPublic:    isPublic,
		AutoResolve: autoResolve,
		StartedAt:   *startedAt,
		ResolvedAt:  resolvedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.Repo.CreateIncident(ctx, tx, incident); err != nil {
		zap.L().Error("Failed to create incident", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create incident")
	}

	for _, mid := range monitorIDs {
		if err := h.Repo.CreateIncidentMonitor(ctx, tx, incident.ID, mid); err != nil {
			zap.L().Error("Failed to link incident to monitor", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to link incident to monitor")
		}
	}

	msg := req.Message
	if msg == "" {
		msg = string(status)
	}

	event := models.EventTimeline{
		IncidentID: incidentID,
		CreatedBy:  userID,
		Message:    msg,
		EventType:  eventTypeFromStatus(status),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := h.Repo.CreateEventTimeline(ctx, tx, event); err != nil {
		zap.L().Error("Failed to create incident event", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create incident event")
	}

	if err := h.Repo.CommitTransaction(tx, ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	resp := struct {
		Incident models.Incident      `json:"incident"`
		Event    models.EventTimeline `json:"event"`
	}{
		Incident: incident,
		Event:    event,
	}

	return c.JSON(http.StatusOK, response.Success("Incident created successfully", resp))
}
