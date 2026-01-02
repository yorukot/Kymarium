package incident

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/models"
	authutil "github.com/yorukot/kymarium/utils/auth"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

type updateIncidentStatusRequest struct {
	Status  models.IncidentStatus `json:"status" validate:"required,oneof=detected investigating identified monitoring resolved"`
	Message string                `json:"message"`
	Public  *bool                 `json:"public"`
}

// UpdateIncidentStatus godoc
// @Summary Update incident status
// @Description Updates an incident's status and records a timeline event
// @Tags incidents
// @Accept json
// @Produce json
// @Param teamID path string true "Team ID"
// @Param incidentID path string true "Incident ID"
// @Param request body updateIncidentStatusRequest true "Incident status payload"
// @Success 200 {object} response.SuccessResponse "Incident status updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Incident not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/incidents/{incidentID}/status [post]
func (h *Handler) UpdateIncidentStatus(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	incidentID, err := strconv.ParseInt(c.Param("incidentID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid incident ID")
	}

	var req updateIncidentStatusRequest
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
	defer h.Repo.DeferRollback(ctx, tx)

	member, err := h.Repo.GetTeamMemberByUserID(ctx, tx, teamID, *userID)
	if err != nil {
		zap.L().Error("Failed to get team membership", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get team membership")
	}

	if member == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Incident not found")
	}

	existing, err := h.Repo.GetIncidentByIDForTeam(ctx, tx, teamID, incidentID)
	if err != nil {
		zap.L().Error("Failed to get incident", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get incident")
	}

	if existing == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Incident not found")
	}

	now := time.Now().UTC()

	var resolvedAt *time.Time
	if req.Status == models.IncidentStatusResolved {
		resolvedAt = &now
	}

	updatedIncident, err := h.Repo.UpdateIncidentStatus(ctx, tx, existing.ID, req.Status, resolvedAt, now)
	if err != nil {
		zap.L().Error("Failed to update incident status", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update incident status")
	}

	if updatedIncident == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Incident not found")
	}

	eventType := eventTypeFromStatus(req.Status)

	msg := req.Message
	if msg == "" {
		msg = string(req.Status)
	}

	event := models.EventTimeline{
		IncidentID: updatedIncident.ID,
		CreatedBy:  userID,
		Message:    msg,
		EventType:  eventType,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := h.Repo.CreateEventTimeline(ctx, tx, event); err != nil {
		zap.L().Error("Failed to record incident status event", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to record incident status event")
	}

	if err := h.Repo.CommitTransaction(ctx, tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	resp := struct {
		Incident models.Incident      `json:"incident"`
		Event    models.EventTimeline `json:"event"`
	}{
		Incident: *updatedIncident,
		Event:    event,
	}

	return c.JSON(http.StatusOK, response.Success("Incident status updated successfully", resp))
}

func eventTypeFromStatus(status models.IncidentStatus) models.EventType {
	switch status {
	case models.IncidentStatusResolved:
		return models.IncidentEventTypeManuallyResolved
	case models.IncidentStatusInvestigating:
		return models.IncidentEventTypeInvestigating
	case models.IncidentStatusIdentified:
		return models.IncidentEventTypeIdentified
	case models.IncidentStatusMonitoring:
		return models.IncidentEventTypeMonitoring
	case models.IncidentStatusDetected:
		return models.IncidentEventTypeDetected
	default:
		return models.IncidentEventTypeUpdate
	}
}
