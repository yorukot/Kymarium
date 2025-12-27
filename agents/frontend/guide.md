# Frontend Agent Guide

- Scope: `website/` only; avoid backend Go files unless a task specifically calls for backend changes to support the UI.
- Coordinate with API contracts in `api/` and shared models under `models/` when consuming backend data; do not modify schemas unless requested.
- Keep environment secrets out of git; reuse existing config patterns and `.env` examples where available.
- Match existing design system and build tooling inside `website/`; prefer localized changes over sweeping refactors unless asked.
- Before shipping, run the frontend lint/build/test commands defined in `website/` tooling (e.g., package scripts) to catch regressions.
- Svelte/SvelteKit workflow: start by listing available docs via the Svelte MCP `list-sections` tool, fetch all relevant sections with `get-documentation`, run `svelte-autofixer` on any Svelte code before sharing it, and offer a Playground link with `playground-link` when the user wants to preview code that was not written to the repo.
- Current sidebar task context: default to pure UI/UX, but API wiring is allowed when explicitly requested; lean on shadcn tokens from `src/routes/layout.css` instead of hardcoded colors, radii, or spacing; prioritize UX clarity before visual polish.
