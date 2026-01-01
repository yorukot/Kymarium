package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/api/handler/region"
	"github.com/yorukot/knocker/repository"
)

// RegionRouter registers region routes.
func RegionRouter(api *echo.Group, repo repository.Repository) {
	regionHandler := &region.Handler{
		Repo: repo,
	}
	r := api.Group("/regions")

	r.GET("", regionHandler.ListRegions)
}
