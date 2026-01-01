package router

import (
	"github.com/labstack/echo/v4"
	invitetoken "github.com/yorukot/knocker/api/handler/invite_token"
	"github.com/yorukot/knocker/repository"
)

// InviteTokenRouter handles invite token routes.
func InviteTokenRouter(api *echo.Group, repo repository.Repository) {
	handler := &invitetoken.Handler{Repo: repo}

	r := api.Group("/invite-tokens")
	r.GET("/:token", handler.GetInviteToken)
}
