package models

import "time"

// MonitorAnalyticsBucket represents a single aggregated bucket for a monitor.
// It is sourced from the Timescale continuous aggregate monitor_30min_summary.
type MonitorAnalyticsBucket struct {
	Bucket     time.Time `json:"bucket" db:"bucket"`
	RegionID   int64     `json:"region_id" db:"region_id"`
	TotalCount int64     `json:"total_count" db:"total_count"`
	GoodCount  int64     `json:"good_count" db:"good_count"`
	P50Ms      float64   `json:"p50_ms" db:"p50_ms"`
	P75Ms      float64   `json:"p75_ms" db:"p75_ms"`
	P90Ms      float64   `json:"p90_ms" db:"p90_ms"`
	P95Ms      float64   `json:"p95_ms" db:"p95_ms"`
	P99Ms      float64   `json:"p99_ms" db:"p99_ms"`
}
