package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/api/handler/team"
	"github.com/yorukot/kymarium/api/middleware"
	"github.com/yorukot/kymarium/repository"
)

// TeamMemberRouter handles team member-related routes.
func TeamMemberRouter(api *echo.Group, repo repository.Repository) {
	teamHandler := &team.Handler{Repo: repo}

	r := api.Group("/teams/:teamID/members", middleware.AuthRequiredMiddleware)
	r.GET("", teamHandler.ListMembers)
	r.DELETE("/:userID", teamHandler.RemoveMember)
}
