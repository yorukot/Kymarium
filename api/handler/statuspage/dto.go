package statuspage

import (
	"github.com/yorukot/kymarium/models"
)

type statusPageGroupInput struct {
	ID        *int64                       `json:"id,string,omitempty" form:"id"`
	Name      string                       `json:"name" form:"name" validate:"required,min=1,max=255"`
	Type      models.StatusPageElementType `json:"type" form:"type" validate:"required,oneof=historical_timeline current_status_indicator"`
	SortOrder int                          `json:"sort_order" form:"sort_order" validate:"min=1"`
}

type statusPageMonitorInput struct {
	ID        *int64                       `json:"id,string,omitempty" form:"id"`
	MonitorID int64                        `json:"monitor_id,string" form:"monitor_id" validate:"required"`
	GroupID   *int64                       `json:"group_id,string,omitempty" form:"group_id"`
	Name      string                       `json:"name" form:"name" validate:"required,min=1,max=255"`
	Type      models.StatusPageElementType `json:"type" form:"type" validate:"required,oneof=historical_timeline current_status_indicator"`
	SortOrder int                          `json:"sort_order" form:"sort_order" validate:"min=1"`
}

type statusPageElementInput struct {
	ID        *int64                       `json:"id,string,omitempty" form:"id"`
	Name      string                       `json:"name" form:"name" validate:"required,min=1,max=255"`
	Type      models.StatusPageElementType `json:"type" form:"type" validate:"required,oneof=historical_timeline current_status_indicator"`
	SortOrder int                          `json:"sort_order" form:"sort_order" validate:"min=1"`
	Monitor   bool                         `json:"monitor" form:"monitor"`
	MonitorID *int64                       `json:"monitor_id,string,omitempty" form:"monitor_id"`
	Monitors  []statusPageMonitorInput     `json:"monitors" form:"monitors" validate:"dive"`
}

type statusPageUpsertRequest struct {
	Title    string                   `json:"title" form:"name" validate:"required,min=1,max=255"`
	Slug     string                   `json:"slug" form:"slug" validate:"required,min=3,max=255"`
	Icon     []byte                   `json:"icon,omitempty" form:"icon"`
	Elements []statusPageElementInput `json:"elements" form:"elements" validate:"dive"`
	Groups   []statusPageGroupInput   `json:"groups" form:"groups" validate:"dive"`
	Monitors []statusPageMonitorInput `json:"monitors" form:"monitors" validate:"dive"`
}

type statusPageElementResponse struct {
	ID           string                       `json:"id"`
	StatusPageID string                       `json:"status_page_id"`
	Name         string                       `json:"name"`
	Type         models.StatusPageElementType `json:"type"`
	SortOrder    int                          `json:"sort_order"`
	Monitor      bool                         `json:"monitor"`
	MonitorID    *string                      `json:"monitor_id,omitempty"`
	Monitors     []models.StatusPageMonitor   `json:"monitors"`
}

type statusPageResponse struct {
	StatusPage models.StatusPage           `json:"status_page"`
	Elements   []statusPageElementResponse `json:"elements"`
}
