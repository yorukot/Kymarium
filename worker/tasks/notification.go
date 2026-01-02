package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/yorukot/kymarium/models"
)

// NotificationPayload represents a notification dispatch request.
type NotificationPayload struct {
	TeamID         int64       `json:"team_id,string"`
	MonitorID      int64       `json:"monitor_id,string"`
	NotificationID int64       `json:"notification_id,string"`
	RegionID       int64       `json:"region_id,string"`
	Ping           models.Ping `json:"ping"`
	Detail         string      `json:"detail,omitempty"`
}

// NewNotificationDispatch builds an Asynq task to send a notification.
func NewNotificationDispatch(payload NotificationPayload) (*asynq.Task, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeNotificationDispatch, body), nil
}
