# Repository Guidelines

## Project Structure & Module Organization
- `app/` holds the Next.js App Router routes, layouts, and pages (e.g., `app/page.tsx`).
- `components/`, `hooks/`, and `lib/` provide shared UI, React hooks, and utilities.
- `public/` contains static assets served at the site root.
- Config lives in `next.config.ts`, `tsconfig.json`, `eslint.config.mjs`, and `postcss.config.mjs`.

## Build, Test, and Development Commands
- `pnpm dev` (or `npm run dev`): start the local Next.js dev server.
- `pnpm build`: production build.
- `pnpm start`: run the production server after a build.
- `pnpm lint`: run ESLint on the codebase.

## Coding Style & Naming Conventions
- TypeScript + React (Next.js App Router). Prefer idiomatic React patterns and co-locate route files under `app/`.
- Indentation follows the project’s default formatter (TypeScript/ESLint rules); keep style changes minimal and consistent.
- Use descriptive component names and place reusable UI in `components/` rather than `app/`.

## Testing Guidelines
- No test runner is configured in `package.json` yet. If you add tests, document the command in this file and keep tests alongside features or in a dedicated `__tests__/` folder.

## Commit & Pull Request Guidelines
- Commit messages follow a conventional prefix pattern (e.g., `feat:`, `refactor:`, `build:`, `frontend:`). Keep them short and scoped.
- PRs should include a clear summary, list of key changes, and screenshots/GIFs for UI updates.

## Security & Configuration Tips
- Do not commit secrets. Use environment files for local config and keep sensitive values out of git.
- Prefer existing config patterns in this repository over introducing new ones.
