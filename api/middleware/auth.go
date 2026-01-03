package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/repository"
	"go.uber.org/zap"
)

func attachUserIDFromSession(c echo.Context, repo repository.Repository, requireAuth bool) error {
	sessionCookie, err := c.Cookie(models.CookieNameSession)
	if err != nil || sessionCookie.Value == "" {
		if requireAuth {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}
		return nil
	}

	tx, err := repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer repo.DeferRollback(c.Request().Context(), tx)

	session, err := repo.GetSessionByToken(c.Request().Context(), tx, sessionCookie.Value)
	if err != nil {
		zap.L().Error("Failed to get session", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get session")
	}

	if session == nil {
		if requireAuth {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}
		return nil
	}

	if time.Now().After(session.ExpiresAt) {
		if _, err := repo.DeleteSessionByToken(c.Request().Context(), tx, sessionCookie.Value); err != nil {
			zap.L().Error("Failed to delete expired session", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete expired session")
		}
		if err := repo.CommitTransaction(c.Request().Context(), tx); err != nil {
			zap.L().Error("Failed to commit transaction", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
		}
		if requireAuth {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}
		return nil
	}

	if err := repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		zap.L().Error("Failed to commit transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	c.Set(string(UserIDKey), strconv.FormatInt(session.UserID, 10))
	return nil
}

// AuthRequiredMiddleware is the middleware for the auth required
func AuthRequiredMiddleware(repo repository.Repository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := attachUserIDFromSession(c, repo, true); err != nil {
				return err
			}
			return next(c)
		}
	}
}

// AuthOptionalMiddleware is the middleware for the auth optional
func AuthOptionalMiddleware(repo repository.Repository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := attachUserIDFromSession(c, repo, false); err != nil {
				// For optional auth, continue even if session lookup fails
				return next(c)
			}
			return next(c)
		}
	}
}
