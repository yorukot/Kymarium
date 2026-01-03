package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/api/handler/notification"
	"github.com/yorukot/kymarium/api/middleware"
	"github.com/yorukot/kymarium/repository"
)

// NotificationRouter registers notification routes.
func NotificationRouter(api *echo.Group, repo repository.Repository) {
	notificationHandler := &notification.Handler{
		Repo: repo,
	}
	r := api.Group("/teams/:teamID/notifications", middleware.AuthRequiredMiddleware(repo))

	r.POST("", notificationHandler.New)
	r.GET("", notificationHandler.ListNotifications)
	r.GET("/:id", notificationHandler.GetNotification)
	r.PATCH("/:id", notificationHandler.UpdateNotification)
	r.DELETE("/:id", notificationHandler.DeleteNotification)
	r.POST("/:id/test", notificationHandler.TestNotification)
}
