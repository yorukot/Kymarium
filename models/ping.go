package models

import (
	"time"
)

// PingStatus represents the result of a monitor ping.
type PingStatus string

// PingStatus values.
const (
	PingStatusSuccessful PingStatus = "successful"
	PingStatusFailed     PingStatus = "failed"
	PingStatusTimeout    PingStatus = "timeout"
)

// Ping represents a ping record in the database
type Ping struct {
	Time      time.Time  `json:"time" db:"time"`
	MonitorID int64      `json:"monitor_id,string" db:"monitor_id"`
	RegionID  int64      `json:"region_id,string" db:"region_id"`
	Latency   int        `json:"latency" db:"latency"`
	Status    PingStatus `json:"status" db:"status"`
}
