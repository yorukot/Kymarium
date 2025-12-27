package incident

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

type updateIncidentRequest struct {
	Title       *string `json:"title"`
	Public      *bool   `json:"public"`
	AutoResolve *bool   `json:"auto_resolve"`
}

// UpdateIncident godoc
// @Summary Update incident settings
// @Description Updates an incident's visibility and auto-resolve settings
// @Tags incidents
// @Accept json
// @Produce json
// @Param teamID path string true "Team ID"
// @Param incidentID path string true "Incident ID"
// @Param request body updateIncidentRequest true "Incident update payload"
// @Success 200 {object} response.SuccessResponse "Incident updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Incident not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/incidents/{incidentID} [patch]
func (h *IncidentHandler) UpdateIncident(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	incidentID, err := strconv.ParseInt(c.Param("incidentID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid incident ID")
	}

	var req updateIncidentRequest
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := json.Unmarshal(body, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	_, titleProvided := raw["title"]
	_, publicProvided := raw["public"]
	_, autoResolveProvided := raw["auto_resolve"]

	if publicProvided && req.Public == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Public must be a boolean")
	}

	if autoResolveProvided && req.AutoResolve == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Auto-resolve must be a boolean")
	}

	if !titleProvided && !publicProvided && !autoResolveProvided {
		return echo.NewHTTPError(http.StatusBadRequest, "At least one field must be provided to update")
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

	isPublic := existing.IsPublic
	if req.Public != nil {
		isPublic = *req.Public
	}

	autoResolve := existing.AutoResolve
	if autoResolveProvided && req.AutoResolve != nil {
		autoResolve = *req.AutoResolve
	}

	title := existing.Title
	if titleProvided {
		if req.Title == nil {
			title = nil
		} else {
			trimmed := strings.TrimSpace(*req.Title)
			if len(trimmed) > 255 {
				return echo.NewHTTPError(http.StatusBadRequest, "Title must be 255 characters or fewer")
			}
			if trimmed == "" {
				title = nil
			} else {
				title = &trimmed
			}
		}
	}

	now := time.Now().UTC()
	updated, err := h.Repo.UpdateIncidentSettings(ctx, tx, existing.ID, isPublic, autoResolve, title, now)
	if err != nil {
		zap.L().Error("Failed to update incident settings", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update incident")
	}

	if updated == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Incident not found")
	}

	if err := h.Repo.CommitTransaction(tx, ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Incident updated successfully", updated))
}
