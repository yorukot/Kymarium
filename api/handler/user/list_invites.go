package user

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

// ListInvites godoc
// @Summary List pending invites
// @Description Lists pending team invites for the authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} response.SuccessResponse "Invites retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/me/invites [get]
func (h *Handler) ListInvites(c echo.Context) error {
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

	invites, err := h.Repo.ListPendingTeamInvitesByUserID(c.Request().Context(), tx, *userID, time.Now().UTC())
	if err != nil {
		zap.L().Error("Failed to list pending invites", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list invites")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Invites retrieved successfully", invites))
}
