package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/api/handler/monitor"
	"github.com/yorukot/knocker/api/middleware"
	"github.com/yorukot/knocker/repository"
)

// MonitorRouter handles monitor-related routes
func MonitorRouter(api *echo.Group, repo repository.Repository) {
	monitorHandler := &monitor.Handler{
		Repo: repo,
	}

	r := api.Group("/teams/:teamID/monitors", middleware.AuthRequiredMiddleware)
	r.POST("", monitorHandler.CreateMonitor)
	r.GET("", monitorHandler.ListMonitors)
	r.GET("/:id", monitorHandler.GetMonitor)
	r.PUT("/:id", monitorHandler.UpdateMonitor)
	r.DELETE("/:id", monitorHandler.DeleteMonitor)
	r.GET("/:id/analytics", monitorHandler.GetAnalytics)
}
