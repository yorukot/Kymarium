package models

import "time"

// Ping represents a ping record in the database
type Ping struct {
	Time      time.Time `json:"time" db:"time"`
	MonitorID int64     `json:"monitor_id" db:"monitor_id"`
	Latency   int16     `json:"latency" db:"latency"`
}
