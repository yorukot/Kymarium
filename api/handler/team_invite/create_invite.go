package teaminvite

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/models"
	authutil "github.com/yorukot/kymarium/utils/auth"
	"github.com/yorukot/kymarium/utils/config"
	"github.com/yorukot/kymarium/utils/encrypt"
	"github.com/yorukot/kymarium/utils/id"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

type createTeamInviteRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
	Role  string `json:"role" validate:"omitempty,oneof=admin member viewer"`
}

// CreateInvite godoc
// @Summary Create a team invite
// @Description Sends an invite to a user email (owner/admin only)
// @Tags team-invites
// @Accept json
// @Produce json
// @Param teamID path string true "Team ID"
// @Param request body createTeamInviteRequest true "Invite request"
// @Success 200 {object} response.SuccessResponse "Invite created"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or team ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Team or user not found"
// @Failure 409 {object} response.ErrorResponse "Invite already exists"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/invites [post]
func (h *Handler) CreateInvite(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	var req createTeamInviteRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	role, err := normalizeInviteRole(req.Role)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid role")
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
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to invite users to this team")
	}

	invitedUserID, err := h.Repo.GetUserIDByEmail(c.Request().Context(), tx, req.Email)
	if err != nil {
		zap.L().Error("Failed to lookup invitee", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to lookup invitee")
	}

	if invitedUserID == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	existingMember, err := h.Repo.GetTeamMemberByUserID(c.Request().Context(), tx, teamID, *invitedUserID)
	if err != nil {
		zap.L().Error("Failed to get team membership for invitee", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check team membership")
	}

	if existingMember != nil {
		return echo.NewHTTPError(http.StatusConflict, "User is already a team member")
	}

	pendingInvite, err := h.Repo.GetPendingTeamInviteByTeamAndUser(c.Request().Context(), tx, teamID, *invitedUserID)
	if err != nil {
		zap.L().Error("Failed to check existing invite", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check existing invite")
	}

	if pendingInvite != nil {
		return echo.NewHTTPError(http.StatusConflict, "Invite already exists")
	}

	inviteID, err := id.GetID()
	if err != nil {
		zap.L().Error("Failed to generate invite ID", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate invite ID")
	}

	token, err := encrypt.GenerateRandomString(32)
	if err != nil {
		zap.L().Error("Failed to generate invite token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate invite token")
	}

	expiresAt := time.Now().UTC().Add(7 * 24 * time.Hour)
	now := time.Now().UTC()

	invite := models.TeamInvite{
		ID:           inviteID,
		TeamID:       teamID,
		InvitedBy:    *userID,
		InvitedTo:    *invitedUserID,
		InvitedEmail: req.Email,
		Role:         role,
		Status:       models.InviteStatusPending,
		Token:        &token,
		ExpiresAt:    expiresAt,
		UpdatedAt:    now,
		CreatedAt:    now,
	}

	if err := h.Repo.CreateTeamInvite(c.Request().Context(), tx, invite); err != nil {
		zap.L().Error("Failed to create team invite", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create team invite")
	}

	if config.Env().SMTPEnabled {
		team, err := h.Repo.GetTeamForUser(c.Request().Context(), tx, teamID, *userID)
		if err != nil {
			zap.L().Error("Failed to fetch team", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch team")
		}

		if team == nil {
			return echo.NewHTTPError(http.StatusNotFound, "Team not found")
		}

		acceptURL, err := buildInviteAcceptURL(token)
		if err != nil {
			zap.L().Error("Failed to build invite URL", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to build invite URL")
		}

		subject := "You have been invited to join a team"
		body := buildInviteEmailBody(team.Name, acceptURL, expiresAt)
		if err := config.SendEmail(req.Email, nil, nil, subject, body); err != nil {
			zap.L().Error("Failed to send invite email", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send invite email")
		}
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Invite created", invite))
}
