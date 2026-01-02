package user

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	authutil "github.com/yorukot/kymarium/utils/auth"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

// RevokeSession godoc
// @Summary Revoke session
// @Description Revokes a refresh token session for the authenticated user
// @Tags users
// @Produce json
// @Param sessionID path string true "Session ID"
// @Success 200 {object} response.SuccessResponse "Session revoked successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid session ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Session not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/me/sessions/{sessionID}/revoke [post]
func (h *Handler) RevokeSession(c echo.Context) error {
	sessionID, err := strconv.ParseInt(c.Param("sessionID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid session ID")
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

	updated, err := h.Repo.UpdateRefreshTokenUsedAtByID(c.Request().Context(), tx, *userID, sessionID, time.Now())
	if err != nil {
		zap.L().Error("Failed to revoke refresh token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to revoke session")
	}

	if !updated {
		return echo.NewHTTPError(http.StatusNotFound, "Session not found")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.SuccessMessage("Session revoked successfully"))
}
