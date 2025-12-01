package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/api/handler/notification"
	"github.com/yorukot/knocker/api/middleware"
	"github.com/yorukot/knocker/repository"
)

// Auth router going to route register signin etc
func NotificationRouter(api *echo.Group, repo repository.Repository) {
	notificationHandler := &notification.NotificationHandler{
		Repo: repo,
	}
	r := api.Group("/teams/:teamID/notifications", middleware.AuthRequiredMiddleware)

	r.POST("", notificationHandler.New)
	r.GET("", notificationHandler.ListNotifications)
	r.GET("/:id", notificationHandler.GetNotification)
	r.PATCH("/:id", notificationHandler.UpdateNotification)
	r.DELETE("/:id", notificationHandler.DeleteNotification)
	r.POST("/:id/test", notificationHandler.TestNotification)
}
