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

type IncidentSeverity string

const (
	IncidentSeverityEmergency IncidentSeverity = "emergency"
	IncidentSeverityCritical  IncidentSeverity = "critical"
	IncidentSeverityMajor     IncidentSeverity = "major"
	IncidentSeverityMinor     IncidentSeverity = "minor"
	IncidentSeverityInfo      IncidentSeverity = "info"
)

type EventType string

const (
	IncidentEventTypeDetected         EventType = "detected"
	IncidentEventTypeNotificationSent EventType = "notification_sent"
	IncidentEventTypeManuallyResolved EventType = "manually_resolved"
	IncidentEventTypeAutoResolved     EventType = "auto_resolved"
	IncidentEventTypeUnpublished      EventType = "unpublished"
	IncidentEventTypePublished        EventType = "published"
	IncidentEventTypeInvestigating    EventType = "investigating"
	IncidentEventTypeIdentified       EventType = "identified"
	IncidentEventTypeUpdate           EventType = "update"
	IncidentEventTypeMonitoring       EventType = "monitoring"
)

// Incident represents an incident record in the database
type Incident struct {
	ID          int64            `json:"id,string" db:"id"`
	Title       *string          `json:"title,omitempty" db:"title"`
	Status      IncidentStatus   `json:"status" db:"status"`
	Severity    IncidentSeverity `json:"severity" db:"severity"`
	IsPublic    bool             `json:"is_public" db:"is_public"`
	AutoResolve bool             `json:"auto_resolve" db:"auto_resolve"`
	StartedAt   time.Time        `json:"started_at" db:"started_at"`
	ResolvedAt  *time.Time       `json:"resolved_at,omitempty" db:"resolved_at"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
}

// IncidentWithMonitorID decorates an incident with the related monitor id.
type IncidentWithMonitorID struct {
	Incident
	MonitorID int64 `json:"monitor_id,string" db:"monitor_id"`
}

// IncidentMonitor links incidents to monitors.
type IncidentMonitor struct {
	ID         int64 `json:"id,string" db:"id"`
	IncidentID int64 `json:"incident_id,string" db:"incident_id"`
	MonitorID  int64 `json:"monitor_id,string" db:"monitor_id"`
}

// EventTimeline stores public or internal updates for an incident.
type EventTimeline struct {
	ID         int64     `json:"id,string" db:"id"`
	IncidentID int64     `json:"incident_id,string" db:"event_id"`
	CreatedBy  *int64    `json:"created_by,string,omitempty" db:"created_by"`
	Message    string    `json:"message" db:"message"`
	EventType  EventType `json:"event_type" db:"event_type"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
