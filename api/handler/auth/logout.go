package auth

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/utils/config"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

// Logout godoc
// @Summary Logout
// @Description Logs the user out by invalidating refresh token and clearing cookies
// @Tags auth
// @Produce json
// @Success 200 {object} response.SuccessResponse "Logged out"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/logout [post]
func (h *Handler) Logout(c echo.Context) error {
	expireAuthCookies(c)

	refreshCookie, err := c.Cookie(models.CookieNameRefreshToken)
	if err != nil || refreshCookie == nil || refreshCookie.Value == "" {
		return c.JSON(http.StatusOK, response.SuccessMessage("Logged out"))
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
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get refresh token")
	}

	if refreshToken != nil && refreshToken.UsedAt == nil {
		now := time.Now()
		refreshToken.UsedAt = &now
		if err := h.Repo.UpdateRefreshTokenUsedAt(c.Request().Context(), tx, *refreshToken); err != nil {
			zap.L().Error("Failed to update refresh token used_at", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update refresh token")
		}
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		zap.L().Error("Failed to commit transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.SuccessMessage("Logged out"))
}

func expireAuthCookies(c echo.Context) {
	accessCookie := http.Cookie{
		Name:     models.CookieNameAccessToken,
		Value:    "",
		Path:     "/api",
		Domain:   cookieDomain(),
		HttpOnly: true,
		Secure:   config.Env().AppEnv == config.AppEnvProd,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	}

	refreshCookie := http.Cookie{
		Name:     models.CookieNameRefreshToken,
		Value:    "",
		Path:     "/api/auth/refresh",
		Domain:   cookieDomain(),
		HttpOnly: true,
		Secure:   config.Env().AppEnv == config.AppEnvProd,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	}

	oauthCookie := http.Cookie{
		Name:     models.CookieNameOAuthSession,
		Value:    "",
		Path:     "/api/auth/oauth",
		Domain:   cookieDomain(),
		HttpOnly: true,
		Secure:   config.Env().AppEnv == config.AppEnvProd,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(&accessCookie)
	c.SetCookie(&refreshCookie)
	c.SetCookie(&oauthCookie)
}
