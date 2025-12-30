package team_invite

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/id"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

type updateInviteRequest struct {
	Status string `json:"status" validate:"required,oneof=accepted rejected"`
}

// UpdateInvite godoc
// @Summary Accept or reject a team invite
// @Description Accepts or rejects a pending invite
// @Tags team-invites
// @Accept json
// @Produce json
// @Param teamID path string true "Team ID"
// @Param inviteID path string true "Invite ID"
// @Param request body updateInviteRequest true "Invite update request"
// @Success 200 {object} response.SuccessResponse "Invite updated"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or IDs"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Invite not found"
// @Failure 409 {object} response.ErrorResponse "Invite not pending"
// @Failure 410 {object} response.ErrorResponse "Invite expired"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/invites/{inviteID} [patch]
func (h *TeamInviteHandler) UpdateInvite(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	inviteID, err := strconv.ParseInt(c.Param("inviteID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid invite ID")
	}

	var req updateInviteRequest
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

	tx, err := h.Repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer h.Repo.DeferRollback(tx, c.Request().Context())

	invite, err := h.Repo.GetTeamInviteByID(c.Request().Context(), tx, teamID, inviteID)
	if err != nil {
		zap.L().Error("Failed to get invite", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get invite")
	}

	if invite == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invite not found")
	}

	if invite.InvitedTo != *userID {
		return echo.NewHTTPError(http.StatusForbidden, "You are not allowed to update this invite")
	}

	if invite.Status != models.InviteStatusPending {
		return echo.NewHTTPError(http.StatusConflict, "Invite is not pending")
	}

	if time.Now().UTC().After(invite.ExpiresAt) {
		return echo.NewHTTPError(http.StatusGone, "Invite expired")
	}

	if req.Status == string(models.InviteStatusAccepted) {
		member, err := h.Repo.GetTeamMemberByUserID(c.Request().Context(), tx, teamID, *userID)
		if err != nil {
			zap.L().Error("Failed to check team membership", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check team membership")
		}

		if member != nil {
			return echo.NewHTTPError(http.StatusConflict, "User is already a team member")
		}

		memberID, err := id.GetID()
		if err != nil {
			zap.L().Error("Failed to generate team member ID", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate team member ID")
		}

		now := time.Now().UTC()
		newMember := models.TeamMember{
			ID:        memberID,
			TeamID:    teamID,
			UserID:    *userID,
			Role:      invite.Role,
			UpdatedAt: now,
			CreatedAt: now,
		}

		if err := h.Repo.CreateTeamMember(c.Request().Context(), tx, newMember); err != nil {
			zap.L().Error("Failed to create team member", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create team member")
		}

		updated, err := h.Repo.UpdateTeamInviteStatus(c.Request().Context(), tx, invite.ID, models.InviteStatusAccepted, now, &now, nil, nil)
		if err != nil {
			zap.L().Error("Failed to accept invite", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to accept invite")
		}

		if updated == nil {
			return echo.NewHTTPError(http.StatusNotFound, "Invite not found")
		}

		if err := h.Repo.CommitTransaction(tx, c.Request().Context()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
		}

		return c.JSON(http.StatusOK, response.Success("Invite accepted", updated))
	}

	now := time.Now().UTC()
	updated, err := h.Repo.UpdateTeamInviteStatus(c.Request().Context(), tx, invite.ID, models.InviteStatusRejected, now, nil, &now, nil)
	if err != nil {
		zap.L().Error("Failed to reject invite", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to reject invite")
	}

	if updated == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invite not found")
	}

	if err := h.Repo.CommitTransaction(tx, c.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Invite rejected", updated))
}
