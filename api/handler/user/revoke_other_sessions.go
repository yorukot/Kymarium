package user

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

// RevokeOtherSessions godoc
// @Summary Revoke other sessions
// @Description Revokes all active sessions except the current one
// @Tags users
// @Produce json
// @Success 200 {object} response.SuccessResponse "Other sessions revoked successfully"
// @Failure 400 {object} response.ErrorResponse "Refresh token not found"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/me/sessions/revoke-others [post]
func (h *Handler) RevokeOtherSessions(c echo.Context) error {
	userID, err := authutil.GetUserIDFromContext(c)
	if err != nil {
		zap.L().Error("Failed to parse user ID from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	if userID == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	refreshCookie, err := c.Cookie(models.CookieNameRefreshToken)
	if err != nil || refreshCookie == nil || refreshCookie.Value == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Refresh token not found")
	}

	tx, err := h.Repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer h.Repo.DeferRollback(c.Request().Context(), tx)

	refreshToken, err := h.Repo.GetRefreshTokenByToken(c.Request().Context(), tx, refreshCookie.Value)
	if err != nil {
		zap.L().Error("Failed to get refresh token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to revoke sessions")
	}

	if refreshToken == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Refresh token not found")
	}

	if refreshToken.UserID != *userID {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	_, err = h.Repo.UpdateRefreshTokensUsedAtExcept(c.Request().Context(), tx, *userID, refreshToken.ID, time.Now())
	if err != nil {
		zap.L().Error("Failed to revoke other refresh tokens", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to revoke sessions")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.SuccessMessage("Other sessions revoked successfully"))
}
