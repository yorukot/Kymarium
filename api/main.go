package api

import (
	"net/http"
	"strings"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/yorukot/knocker/api/middleware"
	"github.com/yorukot/knocker/api/router"
	swaggerDocs "github.com/yorukot/knocker/docs"
	"github.com/yorukot/knocker/repository"
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
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     frontendOrigins(env.FrontendDomain),
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Setup routes
	repo := repository.New(db)
	routes(e, repo)
	e.Logger.Infof("Starting server on port %s in %s mode", env.AppPort, env.AppEnv)
	e.Logger.Fatal(e.Start(":" + env.AppPort))
}

// routes sets up the API routes
func routes(e *echo.Echo, repo repository.Repository) {
	// Development-only routes
	if config.Env().AppEnv == config.AppEnvDev {
		// Swagger documentation route
		e.GET("/swagger/*", echoSwagger.WrapHandler)
		// Scalar API reference route for the generated Swagger spec
		e.GET("/reference", scalarDocsHandler())
	}

	// User routes
	api := e.Group("/api")
	router.AuthRouter(api, repo)
	router.UserRouter(api, repo)
	router.TeamRouter(api, repo)
	router.TeamMemberRouter(api, repo)
	router.TeamInviteRouter(api, repo)
	router.InviteTokenRouter(api, repo)
	router.RegionRouter(api, repo)
	router.NotificationRouter(api, repo)
	router.MonitorRouter(api, repo)
	router.IncidentRouter(api, repo)
	router.StatusPageRouter(api, repo)
	router.PublicStatusPageRouter(api, repo)
}

func scalarDocsHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		html, err := scalar.ApiReferenceHTML(&scalar.Options{
			// Use generated swagger spec as inline content to avoid filesystem lookups
			SpecContent: swaggerDocs.SwaggerInfo.ReadDoc(),
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Ridash API Reference",
			},
			DarkMode: true,
		})
		if err != nil {
			zap.L().Error("failed to generate Scalar docs", zap.Error(err))
			return c.String(http.StatusInternalServerError, "could not render API reference")
		}

		return c.HTML(http.StatusOK, html)
	}
}

// frontendOrigins builds allowed origins for CORS from the configured frontend domain.
func frontendOrigins(domain string) []string {
	trimmed := strings.TrimSpace(domain)
	if trimmed == "" {
		return nil
	}

	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return []string{trimmed}
	}

	return []string{"https://" + trimmed, "http://" + trimmed}
}
