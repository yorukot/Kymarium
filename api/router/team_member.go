package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/api/handler/team"
	"github.com/yorukot/knocker/api/middleware"
	"github.com/yorukot/knocker/repository"
)

// TeamMemberRouter handles team member-related routes.
func TeamMemberRouter(api *echo.Group, repo repository.Repository) {
	teamHandler := &team.Handler{Repo: repo}

	r := api.Group("/teams/:teamID/members", middleware.AuthRequiredMiddleware)
	r.GET("", teamHandler.ListMembers)
	r.DELETE("/:userID", teamHandler.RemoveMember)
}
