package team

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

type updateMemberRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin member viewer"`
}

// UpdateMemberRole godoc
// @Summary Update a team member role
// @Description Updates a user's role in a team (owner/admin only, cannot update self, cannot modify owner)
// @Tags teams
// @Accept json
// @Produce json
// @Param teamID path string true "Team ID"
// @Param userID path string true "User ID"
// @Param request body updateMemberRoleRequest true "Update role request"
// @Success 200 {object} response.SuccessResponse "Member role updated"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or IDs"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Member not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/members/{userID} [patch]
func (h *Handler) UpdateMemberRole(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	targetUserID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	var req updateMemberRoleRequest
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

	if *userID == targetUserID {
		return echo.NewHTTPError(http.StatusForbidden, "You cannot update your own role")
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
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to manage team members")
	}

	targetMember, err := h.Repo.GetTeamMemberByUserID(c.Request().Context(), tx, teamID, targetUserID)
	if err != nil {
		zap.L().Error("Failed to get target membership", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get team membership")
	}

	if targetMember == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Member not found")
	}

	if targetMember.Role == models.MemberRoleOwner {
		return echo.NewHTTPError(http.StatusForbidden, "You cannot modify the team owner")
	}

	if member.Role == models.MemberRoleAdmin && targetMember.Role == models.MemberRoleAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "You cannot modify other admins")
	}

	role := models.MemberRole(req.Role)
	now := time.Now().UTC()

	updated, err := h.Repo.UpdateTeamMemberRoleByUserID(c.Request().Context(), tx, teamID, targetUserID, role, now)
	if err != nil {
		zap.L().Error("Failed to update member role", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update member role")
	}

	if updated == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Member not found")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Member role updated", updated))
}
