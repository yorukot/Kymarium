package incident

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	authutil "github.com/yorukot/kymarium/utils/auth"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

type updateIncidentEventRequest struct {
	Message *string `json:"message" validate:"omitempty,min=1,max=1000"`
	Public  *bool   `json:"public"`
}

// UpdateIncidentEvent godoc
// @Summary Update an incident event
// @Description Updates an incident event's visibility between public and private
// @Tags incidents
// @Accept json
// @Produce json
// @Param teamID path string true "Team ID"
// @Param incidentID path string true "Incident ID"
// @Param eventID path string true "Event ID"
// @Param request body updateIncidentEventRequest true "Incident event update payload"
// @Success 200 {object} response.SuccessResponse "Incident event updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Incident or event not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/incidents/{incidentID}/events/{eventID} [patch]
func (h *Handler) UpdateIncidentEvent(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	incidentID, err := strconv.ParseInt(c.Param("incidentID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid incident ID")
	}

	eventID, err := strconv.ParseInt(c.Param("eventID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid event ID")
	}

	var req updateIncidentEventRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Public == nil && req.Message == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "At least one field (public or message) must be provided")
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

	existingEvent, err := h.Repo.GetEventTimelineByID(ctx, tx, incident.ID, eventID)
	if err != nil {
		zap.L().Error("Failed to get incident event", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get incident event")
	}

	if existingEvent == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Incident event not found")
	}

	var newMessage *string
	if req.Message != nil {
		trimmed := strings.TrimSpace(*req.Message)
		newMessage = &trimmed
	}

	now := time.Now().UTC()
	updatedEvent, err := h.Repo.UpdateEventTimeline(ctx, tx, incident.ID, existingEvent.ID, newMessage, req.Public, now)
	if err != nil {
		zap.L().Error("Failed to update incident event", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update incident event")
	}

	if updatedEvent == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Incident event not found")
	}

	if err := h.Repo.CommitTransaction(ctx, tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Incident event updated successfully", updatedEvent))
}
