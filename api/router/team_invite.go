package router

import (
	"github.com/labstack/echo/v4"
	teaminvite "github.com/yorukot/kymarium/api/handler/team_invite"
	"github.com/yorukot/kymarium/api/middleware"
	"github.com/yorukot/kymarium/repository"
)

// TeamInviteRouter handles team invite-related routes.
func TeamInviteRouter(api *echo.Group, repo repository.Repository) {
	handler := &teaminvite.Handler{Repo: repo}

	r := api.Group("/teams/:teamID/invites", middleware.AuthRequiredMiddleware(repo))
	r.POST("", handler.CreateInvite)
	r.GET("", handler.ListInvites)
	r.DELETE("/:inviteID", handler.DeleteInvite)
	r.PATCH("/:inviteID", handler.UpdateInvite)
}
