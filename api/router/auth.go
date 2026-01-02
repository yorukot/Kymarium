package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/api/handler/auth"
	"github.com/yorukot/kymarium/api/middleware"
	"github.com/yorukot/kymarium/repository"
	"github.com/yorukot/kymarium/utils/config"
)

// AuthRouter registers auth-related routes.
func AuthRouter(api *echo.Group, repo repository.Repository) {
	oauthConfig, err := config.GetOAuthConfig()
	if err != nil {
		panic("Failed to initialize OAuth config: " + err.Error())
	}

	authHandler := &auth.Handler{
		Repo:        repo,
		OAuthConfig: oauthConfig,
	}
	r := api.Group("/auth")

	r.GET("/oauth/:provider", authHandler.OAuthEntry, middleware.AuthOptionalMiddleware)
	r.GET("/oauth/:provider/callback", authHandler.OAuthCallback)

	r.GET("/status", authHandler.Status, middleware.AuthRequiredMiddleware)
	r.POST("/register", authHandler.Register)
	r.GET("/verify", authHandler.VerifyEmail)
	r.POST("/verify/resend", authHandler.ResendVerification)
	r.POST("/login", authHandler.Login)
	r.POST("/logout", authHandler.Logout, middleware.AuthOptionalMiddleware)
	r.POST("/refresh", authHandler.RefreshToken)
}
