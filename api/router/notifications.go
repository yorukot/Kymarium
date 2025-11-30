package router

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/api/handler/notification"
	"github.com/yorukot/knocker/api/middleware"
)

// Auth router going to route register signin etc
func NotificationRouter(api *echo.Group, db *pgxpool.Pool) {
	notificationHandler := &notification.NotificationHandler{
		DB: db,
	}
	r := api.Group("/teams/:teamID/notifications", middleware.AuthRequiredMiddleware)

	r.POST("/", notificationHandler.New)
	r.GET("/", notificationHandler.ListNotifications)
	r.GET("/:id", notificationHandler.GetNotification)
	r.PATCH("/:id", notificationHandler.UpdateNotification)
	r.DELETE("/:id", notificationHandler.DeleteNotification)
}
