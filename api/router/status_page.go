package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/api/handler/statuspage"
	"github.com/yorukot/kymarium/api/middleware"
	"github.com/yorukot/kymarium/repository"
)

// StatusPageRouter handles status page routes.
func StatusPageRouter(api *echo.Group, repo repository.Repository) {
	handler := &statuspage.Handler{Repo: repo}

	r := api.Group("/teams/:teamID/status-pages", middleware.AuthRequiredMiddleware(repo))
	r.POST("", handler.CreateStatusPage)
	r.GET("", handler.ListStatusPages)
	r.GET("/:id", handler.GetStatusPage)
	r.PUT("/:id", handler.UpdateStatusPage)
	r.DELETE("/:id", handler.DeleteStatusPage)
}

// PublicStatusPageRouter handles public status page routes.
func PublicStatusPageRouter(api *echo.Group, repo repository.Repository) {
	handler := &statuspage.Handler{Repo: repo}
	api.GET("/status-pages/:slug", handler.GetPublicStatusPage)
}
