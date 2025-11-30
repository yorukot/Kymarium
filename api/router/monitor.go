package router

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/api/handler/monitor"
	"github.com/yorukot/knocker/api/middleware"
)

// MonitorRouter handles monitor-related routes
func MonitorRouter(api *echo.Group, db *pgxpool.Pool) {
	monitorHandler := &monitor.MonitorHandler{
		DB: db,
	}

	r := api.Group("/teams/:teamID/monitors", middleware.AuthRequiredMiddleware)
	r.POST("/", monitorHandler.CreateMonitor)
	r.GET("/", monitorHandler.ListMonitors)
	r.GET("/:id", monitorHandler.GetMonitor)
	r.PUT("/:id", monitorHandler.UpdateMonitor)
	r.DELETE("/:id", monitorHandler.DeleteMonitor)
}
