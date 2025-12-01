package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/yorukot/knocker/models"
)

// NotificationPayload represents a notification dispatch request.
type NotificationPayload struct {
	MonitorID      int64       `json:"monitor_id,string"`
	NotificationID int64       `json:"notification_id,string"`
	Region         string      `json:"region"`
	Ping           models.Ping `json:"status"`
}

func NewNotificationDispatch(payload NotificationPayload) (*asynq.Task, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeNotificationDispatch, body), nil
}
