package invite_token

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/utils/id"
	"go.uber.org/zap"
)

// GetInviteToken godoc
// @Summary Accept a team invite by token
// @Description Accepts a team invite and redirects to the frontend
// @Tags team-invites
// @Produce json
// @Param token path string true "Invite token"
// @Success 302 {string} string "Redirect"
// @Router /invite-tokens/{token} [get]
func (h *InviteTokenHandler) GetInviteToken(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		return c.Redirect(http.StatusFound, buildFrontendInviteURL("invalid", 0, ""))
	}

	tx, err := h.Repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return c.Redirect(http.StatusFound, buildFrontendInviteURL("error", 0, ""))
	}
	defer h.Repo.DeferRollback(tx, c.Request().Context())

	invite, err := h.Repo.GetTeamInviteByToken(c.Request().Context(), tx, token)
	if err != nil {
		zap.L().Error("Failed to fetch invite by token", zap.Error(err))
		return c.Redirect(http.StatusFound, buildFrontendInviteURL("error", 0, ""))
	}

	if invite == nil {
		return c.Redirect(http.StatusFound, buildFrontendInviteURL("invalid", 0, ""))
	}

	if invite.Status == models.InviteStatusAccepted {
		return c.Redirect(http.StatusFound, buildFrontendInviteURL("already", invite.TeamID, invite.InvitedEmail))
	}

	if invite.Status != models.InviteStatusPending {
		return c.Redirect(http.StatusFound, buildFrontendInviteURL("invalid", invite.TeamID, invite.InvitedEmail))
	}

	if time.Now().UTC().After(invite.ExpiresAt) {
		return c.Redirect(http.StatusFound, buildFrontendInviteURL("expired", invite.TeamID, invite.InvitedEmail))
	}

	member, err := h.Repo.GetTeamMemberByUserID(c.Request().Context(), tx, invite.TeamID, invite.InvitedTo)
	if err != nil {
		zap.L().Error("Failed to check team membership", zap.Error(err))
		return c.Redirect(http.StatusFound, buildFrontendInviteURL("error", invite.TeamID, invite.InvitedEmail))
	}

	now := time.Now().UTC()
	if member == nil {
		memberID, err := id.GetID()
		if err != nil {
			zap.L().Error("Failed to generate team member ID", zap.Error(err))
			return c.Redirect(http.StatusFound, buildFrontendInviteURL("error", invite.TeamID, invite.InvitedEmail))
		}

		newMember := models.TeamMember{
			ID:        memberID,
			TeamID:    invite.TeamID,
			UserID:    invite.InvitedTo,
			Role:      invite.Role,
			UpdatedAt: now,
			CreatedAt: now,
		}

		if err := h.Repo.CreateTeamMember(c.Request().Context(), tx, newMember); err != nil {
			zap.L().Error("Failed to create team member", zap.Error(err))
			return c.Redirect(http.StatusFound, buildFrontendInviteURL("error", invite.TeamID, invite.InvitedEmail))
		}
	}

	updated, err := h.Repo.UpdateTeamInviteStatus(c.Request().Context(), tx, invite.ID, models.InviteStatusAccepted, now, &now, nil, nil)
	if err != nil {
		zap.L().Error("Failed to accept invite", zap.Error(err))
		return c.Redirect(http.StatusFound, buildFrontendInviteURL("error", invite.TeamID, invite.InvitedEmail))
	}

	if updated == nil {
		return c.Redirect(http.StatusFound, buildFrontendInviteURL("invalid", invite.TeamID, invite.InvitedEmail))
	}

	if err := h.Repo.CommitTransaction(tx, c.Request().Context()); err != nil {
		zap.L().Error("Failed to commit transaction", zap.Error(err))
		return c.Redirect(http.StatusFound, buildFrontendInviteURL("error", invite.TeamID, invite.InvitedEmail))
	}

	return c.Redirect(http.StatusFound, buildFrontendInviteURL("success", invite.TeamID, invite.InvitedEmail))
}
