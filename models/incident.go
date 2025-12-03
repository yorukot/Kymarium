package models

import "time"

type IncidentStatus string

const (
	IncidentStatusDetected      IncidentStatus = "detected"
	IncidentStatusInvestigating IncidentStatus = "investigating"
	IncidentStatusIdentified    IncidentStatus = "identified"
	IncidentStatusMonitoring    IncidentStatus = "monitoring"
	IncidentStatusResolved      IncidentStatus = "resolved"
)

type IncidentEventType string

const (
	IncidentEventTypeDetected         IncidentEventType = "detected"
	IncidentEventTypeNotificationSent IncidentEventType = "notification_sent"
	IncidentEventTypeManuallyResolved IncidentEventType = "manually_resolved"
	IncidentEventTypeAutoResolved     IncidentEventType = "auto_resolved"
	IncidentEventTypeUnpublished      IncidentEventType = "unpublished"
	IncidentEventTypePublished        IncidentEventType = "published"
	IncidentEventTypeInvestigating    IncidentEventType = "investigating"
	IncidentEventTypeIdentified       IncidentEventType = "identified"
	IncidentEventTypeUpdate           IncidentEventType = "update"
	IncidentEventTypeMonitoring       IncidentEventType = "monitoring"
)

// Incident represents an incident record in the database
type Incident struct {
	ID         int64          `json:"id,string" db:"id"`
	MonitorID  int64          `json:"monitor_id,string" db:"monitor_id"`
	Status     IncidentStatus `json:"status" db:"status"`
	StartedAt  time.Time      `json:"started_at" db:"started_at"`
	ResolvedAt *time.Time     `json:"resolved_at,omitempty" db:"resloved_at"` // Note: schema has typo "resloved_at"
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
}

// IncidentEvent represents an incident event record in the database
type IncidentEvent struct {
	ID         int64             `json:"id,string" db:"id"`
	IncidentID int64             `json:"incident_id,string" db:"incident_id"`
	Message    string            `json:"message" db:"message"`
	EventType  IncidentEventType `json:"event_type" db:"event_type"`
	Public     bool              `json:"public" db:"public"`
	CreatedAt  time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at" db:"updated_at"`
}
