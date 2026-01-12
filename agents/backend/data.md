# Data and Repository Layer

Use this to work on persistence, transactions, IDs, and schema expectations.

## Postgres schema highlights
- Monitors: interval-driven jobs with `failure_threshold`, `recovery_threshold`, `last_checked`, `next_check`, JSON `config`, and `type` (`http` or `ping`).
- Pings: append-only history (`time`, `monitor_id`, `region`, `latency`, `status`). Schema enforces `ping_status` enum.
- Incidents: one active per monitor enforced by `unique_active_incident_per_monitor` index (`migrations/2_unique_active_incidents.up.sql`). Related `incident_events` capture timeline (`event_type` enum) with optional `created_by` user and `public` flag.
- Notifications: per-team channels with type (`discord`, `telegram`, `email` placeholder) and JSON `config`; junction table `monitor_notifications` associates monitors to notification IDs.
- Auth/teams: users, accounts, refresh tokens, teams, and team members back access control; see `migrations/1_initialize_schema.up.sql` for fields.

## Repository patterns (`repository/`)
- Interface lives in `repository/repository.go`; implementations in per-entity files. Mock version for tests in `repository/mock_repository.go` mirrors the interface.
- Transactions: start with `StartTransaction`, defer `DeferRollback`, and finish with `CommitTransaction`. Most handlers wrap reads and writes in a single transaction for consistency.
- Scan helpers: queries use `pgxscan` to map rows to structs; prefer `pgxscan.Select` for slices and `pgxscan.Get` for single rows.
- Naming: methods are resource-prefixed (e.g., `CreateMonitor`, `ListIncidentsByMonitorID`, `BatchInsertPings`). Keep new methods consistent with this convention.

## Monitor scheduling fields
- `last_checked` and `next_check` are updated in batches by the scheduler (`BatchUpdateMonitorsLastChecked`). Next check = `now + interval + jitter` where jitter is up to 30% of interval capped at 20s.
- `failure_threshold` and `recovery_threshold` drive incident detection/resolution (see `agents/backend/monitoring.md`). Store them as smallints but treat as ints in code.

## Ping persistence
- `PingRecorder` buffers and writes pings via `repository.BatchInsertPings` using COPY for throughput, with an upsert fallback to tolerate duplicate keys.
- `ListRecentPingsByMonitorIDAndRegion` fetches the newest pings per monitor/region to evaluate incidents. Keep indexes aligned if you change query patterns.

## Incidents and events
- Creation: `createIncidentIfAbsent` defends against races using a unique constraint; on conflict it reloads the open incident.
- Resolution: `MarkIncidentResolved` stamps `resolved_at` and `updated_at`; events log the change. Manual status changes go through `UpdateIncidentStatus` and write an event with the mapped type.
- Event ordering: `ListIncidentEventsByIncidentID` returns events ordered by `created_at`. Keep event types in sync with enums in `models/incident.go`.

## IDs and time handling
- IDs come from `utils/id` (sonyflake). Never rely on serial sequences.
- Store times in UTC; handlers call `time.Now().UTC()`. Schema uses timestamptz.

## Migrations and docs
- SQL migrations live under `migrations/` and are applied in order. Update both up/down files when adding changes.
- Swagger docs are generated from handler annotations; run `make generate-docs` after modifying API shapes that touch data payloads.
