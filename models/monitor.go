package models

import "time"

// Monitor represents a monitor entity in the database
type Monitor struct {
	ID        int64     `json:"id" db:"id"`
	URL       string    `json:"url" db:"url"`
	Interval  int32     `json:"interval" db:"interval"`
	LastCheck time.Time `json:"last_check" db:"last_check"`
	NextCheck time.Time `json:"next_check" db:"next_check"`
}
