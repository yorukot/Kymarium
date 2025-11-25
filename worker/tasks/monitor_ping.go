package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/yorukot/knocker/models"
)

func NewMonitorPing(monitor models.Monitor) (*asynq.Task, error) {
 	payload, err := json.Marshal(monitor)
    if err != nil {
        return nil, err
    }

	return asynq.NewTask(TypeMonitorPing, payload), nil
}
