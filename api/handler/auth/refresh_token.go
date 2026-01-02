package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/utils/config"
	"github.com/yorukot/kymarium/utils/encrypt"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

// +----------------------------------------------+
// | Refresh Token                                |
// +----------------------------------------------+

// RefreshToken godoc
// @Summary Refresh token
// @Description Refreshes the access token using the refresh token cookie and issues new access/refresh token cookies
// @Tags auth
// @Accept json
// @Produce json
// @Success 201 {object} response.SuccessResponse "Access token generated successfully, new refresh token set in cookie"
// @Failure 401 {object} response.ErrorResponse "Refresh token not found, invalid, or already used"
// @Failure 500 {object} response.ErrorResponse "Internal server error (transaction, database, or token generation failure)"
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(c echo.Context) error {
	userRefreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Refresh token not found")
	}

	// Begin the transaction
	tx, err := h.Repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer h.Repo.DeferRollback(c.Request().Context(), tx)

	// Get the refresh token by token
	checkedRefreshToken, err := h.Repo.GetRefreshTokenByToken(c.Request().Context(), tx, userRefreshToken.Value)
	if err != nil {
		zap.L().Error("Failed to get refresh token by token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get refresh token by token")
	}

	// If the refresh token is not found, return an error
	if checkedRefreshToken == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Refresh token not found")
	}

	// TODO: Need to tell the user might just been hacked
	if checkedRefreshToken.UsedAt != nil {
		zap.L().Warn("Refresh token already used",
			zap.Int64("user_id", checkedRefreshToken.UserID),
			zap.String("ip", c.RealIP()),
			zap.String("user_agent", c.Request().UserAgent()),
		)
		return echo.NewHTTPError(http.StatusUnauthorized, "Refresh token already used")
	}

	// Update the refresh token used_at
	now := time.Now()
	checkedRefreshToken.UsedAt = &now
	if err = h.Repo.UpdateRefreshTokenUsedAt(c.Request().Context(), tx, *checkedRefreshToken); err != nil {
		zap.L().Error("Failed to update refresh token used_at", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update refresh token used_at")
	}

	// Generate new refresh token
	newRefreshToken, err := generateTokenAndSaveRefreshToken(c, h.Repo, tx, checkedRefreshToken.UserID)
	if err != nil {
		zap.L().Error("Failed to generate refresh token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate refresh token")
	}

	// Commit the transaction
	h.Repo.CommitTransaction(c.Request().Context(), tx)

	// Generate the refresh token cookie
	refreshTokenCookie := generateRefreshTokenCookie(newRefreshToken)
	c.SetCookie(&refreshTokenCookie)

	// Generate AccessTokenClaims
	accessTokenClaims := encrypt.JWTSecret{
		Secret: config.Env().JWTSecretKey,
	}

	// Generate the access token
	accessToken, err := accessTokenClaims.GenerateAccessToken(config.Env().AppName, strconv.FormatInt(checkedRefreshToken.UserID, 10), time.Now().Add(time.Duration(config.Env().AccessTokenExpiresAt)*time.Second))
	if err != nil {
		zap.L().Error("Failed to generate access token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate access token")
	}

	accessTokenCookie := generateAccessTokenCookie(accessToken)
	c.SetCookie(&accessTokenCookie)

	return c.JSON(http.StatusCreated, response.SuccessMessage("Access token refreshed successfully"))
}
