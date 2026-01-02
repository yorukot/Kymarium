package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/utils/config"
	"github.com/yorukot/kymarium/utils/encrypt"
	"go.uber.org/zap"
)

// authMiddlewareLogic is the logic for the auth middleware
func authMiddlewareLogic(token string) (*encrypt.AccessTokenClaims, error) {
	JWTSecret := encrypt.JWTSecret{
		Secret: config.Env().JWTSecretKey,
	}

	valid, claims, err := JWTSecret.ValidateAccessTokenAndGetClaims(token)
	if err != nil {
		zap.L().Error("Failed to validate access token", zap.Error(err))
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	if !valid {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	return &claims, nil
}

// AuthRequiredMiddleware is the middleware for the auth required
func AuthRequiredMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessCookie, err := c.Cookie(models.CookieNameAccessToken)
		if err != nil || accessCookie.Value == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		claims, err := authMiddlewareLogic(accessCookie.Value)
		if err != nil {
			return err
		}

		c.Set(string(UserIDKey), claims.Subject)
		return next(c)
	}
}

// AuthOptionalMiddleware is the middleware for the auth optional
func AuthOptionalMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessCookie, err := c.Cookie(models.CookieNameAccessToken)
		if err != nil || accessCookie.Value == "" {
			return next(c)
		}

		claims, err := authMiddlewareLogic(accessCookie.Value)
		if err != nil {
			// For optional auth, continue even if token is invalid
			return next(c)
		}

		c.Set(string(UserIDKey), claims.Subject)
		return next(c)
	}
}
