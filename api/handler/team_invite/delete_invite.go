package teaminvite

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

// DeleteInvite godoc
// @Summary Cancel a team invite
// @Description Cancels a pending invite (owner/admin only)
// @Tags team-invites
// @Produce json
// @Param teamID path string true "Team ID"
// @Param inviteID path string true "Invite ID"
// @Success 200 {object} response.SuccessResponse "Invite canceled"
// @Failure 400 {object} response.ErrorResponse "Invalid team or invite ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Invite not found"
// @Failure 409 {object} response.ErrorResponse "Invite not pending"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/invites/{inviteID} [delete]
func (h *Handler) DeleteInvite(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	inviteID, err := strconv.ParseInt(c.Param("inviteID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid invite ID")
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
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to cancel invites for this team")
	}

	invite, err := h.Repo.GetTeamInviteByID(c.Request().Context(), tx, teamID, inviteID)
	if err != nil {
		zap.L().Error("Failed to get invite", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get invite")
	}

	if invite == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invite not found")
	}

	if invite.Status != models.InviteStatusPending {
		return echo.NewHTTPError(http.StatusConflict, "Invite is not pending")
	}

	now := time.Now().UTC()
	updated, err := h.Repo.UpdateTeamInviteStatus(c.Request().Context(), tx, invite.ID, models.InviteStatusCanceled, now, nil, nil, &now)
	if err != nil {
		zap.L().Error("Failed to cancel invite", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to cancel invite")
	}

	if updated == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invite not found")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Invite canceled", updated))
}
