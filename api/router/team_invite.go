package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/api/handler/team_invite"
	"github.com/yorukot/knocker/api/middleware"
	"github.com/yorukot/knocker/repository"
)

// TeamInviteRouter handles team invite-related routes.
func TeamInviteRouter(api *echo.Group, repo repository.Repository) {
	handler := &team_invite.TeamInviteHandler{Repo: repo}

	r := api.Group("/teams/:teamID/invites", middleware.AuthRequiredMiddleware)
	r.POST("", handler.CreateInvite)
	r.GET("", handler.ListInvites)
	r.DELETE("/:inviteID", handler.DeleteInvite)
	r.PATCH("/:inviteID", handler.UpdateInvite)
}
