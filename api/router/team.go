package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/api/handler/team"
	"github.com/yorukot/knocker/api/middleware"
	"github.com/yorukot/knocker/repository"
)

// TeamRouter handles team-related routes
func TeamRouter(api *echo.Group, repo repository.Repository) {
	teamHandler := &team.Handler{
		Repo: repo,
	}

	r := api.Group("/teams", middleware.AuthRequiredMiddleware)
	r.GET("", teamHandler.ListTeams)
	r.POST("", teamHandler.CreateTeam)
	r.GET("/:id", teamHandler.GetTeam)
	r.PUT("/:id", teamHandler.UpdateTeam)
	r.DELETE("/:id", teamHandler.DeleteTeam)
	r.POST("/:id/leave", teamHandler.LeaveTeam)
}
