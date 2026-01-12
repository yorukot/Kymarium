# Monitoring, Pings, and Incidents

Use this to reason about the monitoring pipeline, ping execution, incident thresholds, and when notifications fire.

## Monitors and scheduling details
- Monitors are due when `next_check <= now` (see `repository.ListMonitorsDueForCheck`). Scheduler polls every 2s and enqueues `monitor:ping:{region}` tasks for each region in `APP_REGIONS`.
- `last_checked` and `next_check` update in batches (`BatchUpdateMonitorsLastChecked`) with jitter up to 30% of the interval (capped at 20s) from `schedular/utils.go` to avoid thundering herds.

## Ping monitors
- Config (`models/monitorm/ping.go`): `host` (required), `timeout_seconds` (default 5s when zero), `packet_size` (default 56 bytes).
- Execution (`core/monitor/ping.go`): sends a single ICMP packet using `prometheus-community/pro-bing`, tries privileged ping first then falls back to unprivileged on permission errors. Timeout uses config or defaults to 5s; interval fixed at 1s. Latency is clamped to a 32-bit ms integer.
- Status mapping: reply received -> `successful`; context timeout/cancel or net timeout -> `timeout`; other errors -> `failed`. Detail/message is the error string when not successful.
- HTTP monitors also run through `core/monitor/http.go`: success is based on accepted status codes (defaults to 2xx) and supports upside-down mode to invert success. Errors/timeouts set status and message accordingly.

## Persisting ping results
- Each ping result is buffered in `PingRecorder` (`worker/handler/ping_recorder.go`) and written in batches via `repository.BatchInsertPings` (COPY fast-path; upsert fallback on duplicates). Flush cadence: 1s ticker; target batch size 1000 with an 80% flush threshold; flush failures fall back to re-queueing the ping in memory.

## Incident lifecycle (automatic)
- Trigger point: after every ping in `processIncident` (`worker/handler/monitor_ping.go`), scoped to the monitor + region of the ping.
- Failure detection:
  - Uses `monitor.FailureThreshold` (>0) and a window of `ceil(threshold * 1.5)` most recent pings for the same region (current ping + history via `ListRecentPingsByMonitorIDAndRegion`).
  - If `failureCount >= threshold`, enough samples exist, and no open incident exists, create an incident with status `detected` and write two events: `detected` and `notification_sent` (both public). Unique index ensures only one open incident per monitor.
  - If an incident is already open and the message changes, append an `update` event (public) but do not send notifications.
- Recovery detection:
  - Requires an open incident and `monitor.RecoveryThreshold` (>0).
  - If the latest `recoveryThreshold` pings for the region (including the current one) are all `successful`, mark the incident resolved (`MarkIncidentResolved`) and add an `auto_resolved` event.
- Messages and details: `incidentMessage` prefixes the region when present and falls back to ping detail/status text. Latency is not part of the message; it lives on the ping and notification payload.
- Notifications: only sent when `handleIncidentFailure` creates a new incident or `handleIncidentRecovery` resolves one. Notification tasks (`notification:dispatch`) include the ping snapshot and detail string.

## Manual incident actions (API)
- Create: `POST /teams/:teamID/monitors/:monitorID/incidents` (`api/handler/incident/create_incident.go`) creates a new incident when none is open. Default status `detected`; supplying `resolved` sets `resolved_at`. The first event matches the status and respects the optional `public` flag.
- Update status: `POST /teams/:teamID/monitors/:monitorID/incidents/:incidentID/status` (`api/handler/incident/update_incident_status.go`) changes status and logs a timeline event. Setting status to `resolved` stamps `resolved_at` and uses event type `manually_resolved`; other statuses map to corresponding event types (`investigating`, `identified`, `monitoring`, etc.).
- Timeline/events: list and append via `/incidents/:incidentID/events` handlers; events store creator (when known), message, event type, and `public` flag.
- Access control: all incident APIs require authenticated team membership (`middleware.AuthRequiredMiddleware` and repository checks for membership/monitor ownership) before returning incident data.
