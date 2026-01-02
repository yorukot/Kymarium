# Repository Guidelines

Use this guide to work consistently in the Kymarium Go codebase and ship changes safely. Role-specific guides: `agents/backend/guide.md` for Go services and `agents/frontend/guide.md` for `website/` work. Backend deep dives: `agents/backend/overview.md` (runtime flow), `agents/backend/monitoring.md` (pings/incidents), `agents/backend/data.md` (schema/repository), `agents/backend/api.md` (HTTP/auth), and `agents/backend/notifications.md` (queues/dispatch). Agents run with the Context7 MCP server available—use it to pull fresh, official docs for libraries and cloud services before coding or reviewing.
- Only read or modify the `website/` folder when the task explicitly requires a frontend change; otherwise leave it untouched to stay focused on backend work.

## Mandatory Workflow Rules
- Before taking any action, read the relevant guide: `agents/backend/*` for backend changes or `agents/frontend/*` for website work; do this step-by-step as you proceed.
- For every step of work, consult the newest official docs via Context7 (or another available MCP server) to confirm current APIs/behaviors before implementing or reviewing.
- If anything is unclear, pause and ask questions rather than assuming; keep asking until requirements are fully understood.
- You may ask as many clarifying questions as needed at any point in the workflow.

## Project Structure & Module Organization
- Entry point: `cmd/main.go` starts API, worker, and schedular (or targeted via `api`/`worker`/`schedular` args).
- HTTP: handlers/middleware/routing in `api/`; Swagger in `api/docs/`.
- Data: `db/` connections, `repository/` queries, `models/` structs, `migrations/` SQL (one active incident per monitor enforced in `migrations/2_*`). IDs come from `utils/id` (sonyflake).
- Background: `worker/` (Asynq consumers, ping handling, notifications) and `schedular/` (queues monitor tasks, updates next_check).
- Utilities: `utils/` for config (`utils/config`), logging (`utils/logger`), IDs, helpers. Artifacts in `tmp/`. Local services in `compose.yaml`.

## How the system operates
- Scheduler polls every 2s for monitors whose `next_check` has elapsed, enqueues `monitor:ping:{region}` tasks for each region (`APP_REGIONS`), and updates `last_checked`/`next_check` with jitter. Worker consumes only the current `APP_REGION` queue.
- Worker executes monitors via `core/monitor.Run` (HTTP or ping), buffers ping results to `pings` through `PingRecorder`, then evaluates incidents (`worker/handler/monitor_ping.go`). Notifications dispatch asynchronously when incidents are created or resolved.
- Refer to `agents/monitoring.md` before modifying ping execution, thresholds, incident logic, or notification triggers.

## API Design Conventions
- RESTful resources (e.g., `/users`, `/teams/{id}`) with Auth-required routes in `api/router/`.
- Handlers live under `api/handler/<resource>/`, one handler per file (e.g., `list.go`, `create.go`); define routes per resource file in `api/router/` (e.g., `monitor.go`, `incident.go`).
- Keep DTOs beside their handlers; generated docs live with the API (`make generate-docs`).

## Coding Style & Naming Conventions
- Go defaults: tabs; always commit gofmt output. Follow the Uber Go style guide.
- Keep package-scoped names short and descriptive; place request/response DTOs near handlers.
- Prefix repository methods by resource (e.g., `UserCreate`, `TeamList`); add env vars in `utils/config` with documented defaults.
- Use `utils/logger` for structured logs and pass contexts for request-scoped work.

## Testing Guidelines
- Prefer table-driven tests in the same package; name `TestFunction_Scenario`.
- Cover new public functions plus validation and data-access edges; prefer fakes over real services.
- Run `go test ./...` before pushing; keep fixtures near the tested package.

## Security & Configuration Tips
- `.env` auto-loads; never commit secrets - use example placeholders. Set `JWT_SECRET_KEY`, OAuth, and DB/Redis values per `utils/config`.
- Avoid default passwords locally; when migrations change, document rollback behavior in the PR.
