package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/api/handler/incident"
	"github.com/yorukot/knocker/api/middleware"
	"github.com/yorukot/knocker/repository"
)

// IncidentRouter handles incident-related routes.
func IncidentRouter(api *echo.Group, repo repository.Repository) {
	incidentHandler := &incident.Handler{
		Repo: repo,
	}

	// Monitor-scoped read/update for backwards compatibility
	r := api.Group("/teams/:teamID/incidents", middleware.AuthRequiredMiddleware)
	r.POST("", incidentHandler.CreateIncident)
	r.GET("", incidentHandler.ListIncidents)
	r.GET("/:incidentID", incidentHandler.GetIncident)
	r.GET("/:incidentID/events", incidentHandler.ListIncidentEvents)
	r.POST("/:incidentID/events", incidentHandler.CreateIncidentEvent)
	r.POST("/:incidentID/status", incidentHandler.UpdateIncidentStatus)
	r.PATCH("/:incidentID", incidentHandler.UpdateIncident)
}
