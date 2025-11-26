package router

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/api/handler/auth"
	"github.com/yorukot/knocker/api/middleware"
	"github.com/yorukot/knocker/utils/config"
)

// Auth router going to route register signin etc
func AuthRouter(api *echo.Group, db *pgxpool.Pool) {
	oauthConfig, err := config.GetOAuthConfig()
	if err != nil {
		panic("Failed to initialize OAuth config: " + err.Error())
	}

	authHandler := &auth.AuthHandler{
		DB:          db,
		OAuthConfig: oauthConfig,
	}
	r := api.Group("/auth")

	r.GET("/oauth/:provider", authHandler.OAuthEntry, middleware.AuthOptionalMiddleware)
	r.GET("/oauth/:provider/callback", authHandler.OAuthCallback)

	r.POST("/register", authHandler.Register)
}
