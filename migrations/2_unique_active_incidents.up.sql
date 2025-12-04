-- Ensure only one non-resolved incident exists per monitor.
CREATE UNIQUE INDEX IF NOT EXISTS unique_active_incident_per_monitor
ON incidents (monitor_id)
WHERE status <> 'resolved';
