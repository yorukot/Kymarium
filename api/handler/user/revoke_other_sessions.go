package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/models"
	authutil "github.com/yorukot/kymarium/utils/auth"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

// RevokeOtherSessions godoc
// @Summary Revoke other sessions
// @Description Revokes all active sessions except the current one
// @Tags users
// @Produce json
// @Success 200 {object} response.SuccessResponse "Other sessions revoked successfully"
// @Failure 400 {object} response.ErrorResponse "Session not found"
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

	sessionCookie, err := c.Cookie(models.CookieNameSession)
	if err != nil || sessionCookie == nil || sessionCookie.Value == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Session not found")
	}

	tx, err := h.Repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer h.Repo.DeferRollback(c.Request().Context(), tx)

	session, err := h.Repo.GetSessionByToken(c.Request().Context(), tx, sessionCookie.Value)
	if err != nil {
		zap.L().Error("Failed to get session", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to revoke sessions")
	}

	if session == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Session not found")
	}

	if session.UserID != *userID {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	_, err = h.Repo.DeleteSessionsExceptToken(c.Request().Context(), tx, *userID, session.Token)
	if err != nil {
		zap.L().Error("Failed to revoke other sessions", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to revoke sessions")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.SuccessMessage("Other sessions revoked successfully"))
}
