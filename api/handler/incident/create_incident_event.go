package incident

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

type createIncidentEventRequest struct {
	Message   string           `json:"message" validate:"required,min=1"`
	EventType models.EventType `json:"event_type" validate:"omitempty,oneof=detected notification_sent manually_resolved auto_resolved unpublished published investigating identified update monitoring"`
}

// CreateIncidentEvent godoc
// @Summary Create an incident event
// @Description Adds a new event to an incident timeline
// @Tags incidents
// @Accept json
// @Produce json
// @Param teamID path string true "Team ID"
// @Param incidentID path string true "Incident ID"
// @Param request body createIncidentEventRequest true "Incident event payload"
// @Success 200 {object} response.SuccessResponse "Incident event created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Incident not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/incidents/{incidentID}/events [post]
func (h *Handler) CreateIncidentEvent(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	incidentID, err := strconv.ParseInt(c.Param("incidentID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid incident ID")
	}

	var req createIncidentEventRequest
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

	incident, err := h.Repo.GetIncidentByIDForTeam(ctx, tx, teamID, incidentID)
	if err != nil {
		zap.L().Error("Failed to get incident", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get incident")
	}

	if incident == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Incident not found")
	}

	eventType := req.EventType
	if eventType == "" {
		eventType = models.IncidentEventTypeUpdate
	}

	now := time.Now().UTC()
	event := models.EventTimeline{
		IncidentID: incident.ID,
		CreatedBy:  userID,
		Message:    req.Message,
		EventType:  eventType,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := h.Repo.CreateEventTimeline(ctx, tx, event); err != nil {
		zap.L().Error("Failed to create incident event", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create incident event")
	}

	if err := h.Repo.CommitTransaction(ctx, tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Incident event created successfully", event))
}
