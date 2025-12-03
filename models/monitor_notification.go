package models

// MonitorNotification represents the junction table between monitors and notifications
// Note: The table name in the schema is "monitor_notificaiton" (with a typo)
type MonitorNotification struct {
	ID             int64 `json:"id,string" db:"id"`
	MonitorID      int64 `json:"monitor_id,string" db:"monitor_id"`
	NotificationID int64 `json:"notification_id,string" db:"notification_id"`
}
