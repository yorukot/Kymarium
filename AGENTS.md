# Repository Guidelines

Use this guide to work consistently in the Knocker Go codebase and ship changes safely.

## Project Structure & Module Organization
- Entry point: `cmd/main.go` starts the API, worker, and schedular (or targeted via `api`/`worker`/`schedular` args).
- HTTP layer: `api/` handlers, middleware, routing; generated Swagger lives in `api/docs/`.
- Data layer: `db/` connections, `repository/` query helpers, `models/` domain structs, `migrations/` SQL.
- Background work: `worker/` for jobs, `schedular/` for cron-like tasks.
- Utilities: `utils/` for config, logging (zap), IDs, helpers. Artifacts land in `tmp/`. Local services defined in `compose.yaml`.

## Build, Test, and Development Commands
- `docker compose up -d` — start TimescaleDB and Dragonfly locally.
- `make build` — build `tmp/knocker` binary. `make run` runs API, worker, and schedular together.
- Dev reload: `make dev` for all, or `make api` / `make worker` / `make schedular` for focused hot reload (air).
- Tests: `make test` or `go test ./...`. Lint/format: `make lint` (runs `go fmt`, `go vet`, `golint`).
- Docs: `make generate-docs` regenerates Swagger from `cmd/main.go`. Seed sample data: `make seed`. Clean artifacts: `make clean`.

## Coding Style & Naming Conventions
- Go defaults: tabs; always commit gofmt output.
- Keep package-scoped names short and descriptive; place request/response DTOs near their handlers.
- Prefix repository methods by resource (e.g., `UserCreate`, `TeamList`); add new env vars in `utils/config` with documented defaults.
- Use `utils/logger` for structured logs and pass contexts for request-scoped work.

## Testing Guidelines
- Prefer table-driven tests alongside code in the same package; name `TestFunction_Scenario`.
- Cover new public functions plus validation and data-access edges; prefer fakes over real services.
- Run `go test ./...` before pushing; keep fixtures near the tested package.

## Commit & Pull Request Guidelines
- Follow Conventional Commits seen in `git log` (`feat: ...`, `fix: ...`, `chore: ...`).
- PRs should state scope, rationale, impact, and linked issues; note config or migration changes and backward-compat expectations.
- Include curl examples or screenshots for API-visible changes and list which tests you ran.

## Security & Configuration Tips
- `.env` auto-loads; never commit secrets—use example placeholders. Set `JWT_SECRET_KEY`, OAuth, and DB/Redis values per `utils/config`.
- Avoid default passwords locally; when migrations change, document rollback behavior in the PR.
