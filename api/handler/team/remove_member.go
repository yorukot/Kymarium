package team

import (
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

// RemoveMember godoc
// @Summary Remove a team member
// @Description Removes a user from a team (owner/admin only, cannot remove self)
// @Tags teams
// @Produce json
// @Param teamID path string true "Team ID"
// @Param userID path string true "User ID"
// @Success 200 {object} response.SuccessResponse "Member removed"
// @Failure 400 {object} response.ErrorResponse "Invalid team or user ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Member not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/members/{userID} [delete]
func (h *TeamHandler) RemoveMember(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	targetUserID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
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
		return echo.NewHTTPError(http.StatusForbidden, "You cannot remove yourself from the team")
	}

	tx, err := h.Repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer h.Repo.DeferRollback(tx, c.Request().Context())

	member, err := h.Repo.GetTeamMemberByUserID(c.Request().Context(), tx, teamID, *userID)
	if err != nil {
		zap.L().Error("Failed to get team membership", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get team membership")
	}

	if member == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Team not found")
	}

	if member.Role != models.MemberRoleOwner && member.Role != models.MemberRoleAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to remove team members")
	}

	targetMember, err := h.Repo.GetTeamMemberByUserID(c.Request().Context(), tx, teamID, targetUserID)
	if err != nil {
		zap.L().Error("Failed to get target membership", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get team membership")
	}

	if targetMember == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Member not found")
	}

	if err := h.Repo.DeleteTeamMemberByUserID(c.Request().Context(), tx, teamID, targetUserID); err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Member not found")
		}
		zap.L().Error("Failed to delete team member", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete team member")
	}

	if err := h.Repo.CommitTransaction(tx, c.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.SuccessMessage("Member removed"))
}
