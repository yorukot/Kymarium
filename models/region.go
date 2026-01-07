package models

// Region represents an available monitoring region.
type Region struct {
	ID          int64  `json:"id,string" db:"id"`
	Name        string `json:"name" db:"name"`
}

// MonitorRegion links a monitor to a specific region with metadata.
type MonitorRegion struct {
	ID        int64 `json:"id,string" db:"id"`
	MonitorID int64 `json:"monitor_id,string" db:"monitor_id"`
	RegionID  int64 `json:"region_id,string" db:"region_id"`
}
