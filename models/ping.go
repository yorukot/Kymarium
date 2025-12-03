package models

import (
	"encoding/json"
	"time"
)

type PingStatus string

const (
	PingStatusSuccessful PingStatus = "successful"
	PingStatusFailed     PingStatus = "failed"
	PingStatusTimeout    PingStatus = "timeout"
)

// Ping represents a ping record in the database
type Ping struct {
	Time      time.Time       `json:"time" db:"time"`
	MonitorID int64           `json:"monitor_id,string" db:"monitor_id"`
	Latency   int             `json:"latency" db:"latency"`
	Status    PingStatus      `json:"status" db:"status"`
	Region    string          `json:"region" db:"region"`
	Data      json.RawMessage `json:"data,omitempty" db:"data"`
}
