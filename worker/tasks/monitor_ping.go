package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/yorukot/kymarium/models"
)

// MonitorPingPayload represents the payload for a monitor ping task
type MonitorPingPayload struct {
	Monitor  models.Monitor `json:"monitor"`
	RegionID int64          `json:"region"`
}

// NewMonitorPing builds an Asynq task for monitor pinging.
func NewMonitorPing(monitor models.Monitor, regionID int64) (*asynq.Task, error) {
	payload := MonitorPingPayload{
		Monitor:  monitor,
		RegionID: regionID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeMonitorPingPattern, payloadBytes), nil
}
