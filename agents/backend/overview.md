# Kymarium System Overview

Use this as the starting point for understanding how Kymarium runs, how monitors are scheduled, and where to look when extending behavior.

## Runtime components
- API (`api/`): Echo HTTP server, JWT auth, Swagger docs in `api/docs/`, routes defined under `api/router/` and handlers under `api/handler/`.
- Scheduler (`schedular/`): polls for monitors whose `next_check` has elapsed and enqueues work to Asynq queues.
- Worker (`worker/`): consumes Asynq queues; runs monitor checks (`monitor:ping:{region}`) and dispatches notifications (`notification:dispatch`). Only consumes `monitor:ping:{APP_REGION}` for the region this worker serves.
- Data layer (`repository/`, `models/`, `db/`, `migrations/`): Postgres via pgx; Redis/Dragonfly for queues; see `compose.yaml` for local services.
- Utilities (`utils/`): config/env loading, logging (`utils/logger`), ID generation (`utils/id`), helpers.

## Monitor check lifecycle
1. Scheduler tick (`schedular/scheduler.go`) runs every 2s. It fetches monitors due for checking with `repository.ListMonitorsDueForCheck`, ordered by `next_check`.
2. For each monitor and each configured region (`APP_REGIONS`), scheduler enqueues an Asynq task built by `worker/tasks.NewMonitorPing` using type `monitor:ping:{region}`. It concurrently updates `last_checked` and `next_check` with jitter via `BatchUpdateMonitorsLastChecked`.
3. Worker `HandleStartServiceTask` (`worker/handler/monitor_ping.go`) runs the monitor through `core/monitor.Run`, capturing status, latency, and any detail message. Ping results default to `failed` with `latency=0` when execution errors.
4. Pings are buffered and persisted in batches to the `pings` table by `PingRecorder` (`worker/handler/ping_recorder.go` -> `repository.BatchInsertPings`). Flush interval is 1s with a ~1000 ping batch size.
5. Incident evaluation runs immediately after each ping inside `processIncident` (`worker/handler/monitor_ping.go`). See `agents/backend/monitoring.md` for thresholds, event rules, and when notifications fire.
6. When an incident needs user visibility, notification tasks are enqueued (`notification:dispatch`). `HandleNotificationDispatch` loads the monitor/notification pair and sends via `core/notification.Send`.

## Key models and schema
- Monitors: `models/monitor.go` with interval, failure/recovery thresholds, raw JSON config decoded by `HTTPConfig`/`PingConfig` from `models/monitorm/*`.
- Pings: `models/ping.go` with status enum (`successful`, `failed`, `timeout`), latency ms, region, and timestamp.
- Incidents and events: `models/incident.go`; only one active incident per monitor (see `migrations/2_unique_active_incidents.up.sql`). Events capture timeline changes and are written both automatically and via the API.

## Local development
- Start dependencies: `docker compose up -d` (TimescaleDB, Dragonfly/Redis).
- Run everything: `make run` (or `make api` / `make worker` / `make schedular` for focused hot reload with air).
- Build/test: `make build`, `make test`; lint via `make lint`. Swagger docs regenerate with `make generate-docs`.
