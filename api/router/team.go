package router

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/api/handler/team"
	"github.com/yorukot/knocker/api/middleware"
)

// TeamRouter handles team-related routes
func TeamRouter(api *echo.Group, db *pgxpool.Pool) {
	teamHandler := &team.TeamHandler{
		DB: db,
	}

	r := api.Group("/teams", middleware.AuthRequiredMiddleware)
	r.GET("", teamHandler.ListTeams)
	r.POST("", teamHandler.CreateTeam)
	r.GET("/:id", teamHandler.GetTeam)
	r.PUT("/:id", teamHandler.UpdateTeam)
	r.DELETE("/:id", teamHandler.DeleteTeam)
}
