package api

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/yorukot/knocker/api/middleware"
	"github.com/yorukot/knocker/api/router"
	_ "github.com/yorukot/knocker/docs"
	"github.com/yorukot/knocker/utils/config"
	"go.uber.org/zap"
)

// Run starts the API server
func Run(db *pgxpool.Pool) {
	zap.L().Info("Starting Ridash API server...")

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.ZapLogger(zap.L()))
	e.Use(echoMiddleware.Recover())

	env := config.Env()

	// Setup routes
	routes(e, db)
	e.Logger.Infof("Starting server on port %s in %s mode", env.AppPort, env.AppEnv)
	e.Logger.Fatal(e.Start(":" + env.AppPort))
}

// routes sets up the API routes
func routes(e *echo.Echo, db *pgxpool.Pool) {
	if config.Env().AppEnv == config.AppEnvDev {
		// Swagger documentation route
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// User routes
	api := e.Group("/api")
	router.AuthRouter(api, db)
}
