package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/api/handler/user"
	"github.com/yorukot/kymarium/api/middleware"
	"github.com/yorukot/kymarium/repository"
)

// UserRouter handles user-related routes
func UserRouter(api *echo.Group, repo repository.Repository) {
	userHandler := &user.Handler{
		Repo: repo,
	}

	r := api.Group("/users", middleware.AuthRequiredMiddleware)
	r.GET("/me", userHandler.GetMe)
	r.PATCH("/me", userHandler.UpdateMe)
	r.GET("/me/account", userHandler.GetAccount)
	r.PATCH("/me/password", userHandler.UpdatePassword)
	r.GET("/me/invites", userHandler.ListInvites)
	r.GET("/me/sessions", userHandler.ListSessions)
	r.POST("/me/sessions/:sessionID/revoke", userHandler.RevokeSession)
	r.POST("/me/sessions/revoke-others", userHandler.RevokeOtherSessions)
}
