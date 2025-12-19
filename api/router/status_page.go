package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/api/handler/statuspage"
	"github.com/yorukot/knocker/api/middleware"
	"github.com/yorukot/knocker/repository"
)

// StatusPageRouter handles status page routes.
func StatusPageRouter(api *echo.Group, repo repository.Repository) {
	handler := &statuspage.Handler{Repo: repo}

	r := api.Group("/teams/:teamID/status-pages", middleware.AuthRequiredMiddleware)
	r.POST("", handler.CreateStatusPage)
	r.PUT(":id", handler.UpdateStatusPage)
}
