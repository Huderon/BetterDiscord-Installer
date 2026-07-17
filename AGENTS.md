# AGENT OPERATING NOTES

These notes are curated for downstream agents. Treat them as binding guidance: each section funnels a different expectation for builds, style, testing and compliance with upstream tooling.

## 1. Project layout recap

- `frontend/` (Bun + Svelte Kit) builds the installer UI that gets embedded into Go via `//go:embed all:frontend/build` in `main.go`.
- `types/` and `api/` contain Go helpers consumed by Wails bindings (see `frontend/src/lib/wailsjs/runtime`).
- `build/` (generated artifacts) holds product icons, release metadata and is not usually edited manually.
- `wails.json` orchestrates hooks that connect Bun scripts to Go builds; keep it synchronized with the commands you run locally.

## 2. Build, dev & release commands

Follow this checklist before you claim a build is validated.

1. `cd frontend && bun install`
   * Bootstraps the renderer toolchain. Required before `wails dev/build` and before running any Bun scripts.
   * The hook in `wails.json` (`frontend:install`) runs this automatically during `wails build`, but you may rerun it to refresh Bun locks.
2. `bun run --bun vite dev`
   * Runs the renderer dev server only. Use when you want instant UI feedback and can ignore the Go side.
   * Requires the backend to be started separately (e.g., `wails dev`) if you need runtime bindings.
3. `wails dev`
   * Launches the Go server, compiles the renderer on demand, and proxies events between them. This is the closest analog to running the installed app locally.
   * Useful for interactive testing (dialogs, runtime events, system prompts).
4. `wails build`
   * Produces the final native binaries. It triggers `frontend:build` (`bun run --bun vite build`), compresses assets, and packages installers.
   * Always run this on the target branch before tagging; artifacts are what release workflows expect.

### Frontend-only helpers

- `bun run --bun vite build` – rebuild the renderer bundle when changing UI/TS/CSS without touching Go.
- `bun run --bun svelte-kit sync` – keeps `frontend/.svelte-kit/generated` in sync; required before `svelte-check`, Vite builds, or when adding routes.
- `bun run --bun svelte-check --tsconfig ./tsconfig.json` – enforces Svelte/TypeScript static checks. Run with `--watch` only when doing extensive UI refactors.
- `bun run --bun eslint .` – lints every JS/TS/Svelte file using the shared config in `frontend/eslint.config.js`.

### Go-specific scripts

- `go test ./...` – quick sanity check for the backend, primarily hitting `types`. Use it as a pre-commit gate for Go logic.
- `gofmt -w .` – required formatting step for `.go` files. Run it in each directory you touch (e.g., `api`, `types`, root) to keep diffs tidy.

## 3. Running isolated tests

- Use the `-run` regex flag when a bundle is too big.
  * `go test ./types -run TestDiscordChannel_Name` – runs just one suite in `types/channel_test.go`. Swap `Name` for `String`, `Exe`, `ParseChannel` as needed.
  * Regex examples: `-run Exe` runs the `Exe` table and the `CurrentPlatform` check.
- Remember `runtime.GOOS` gating.
  * Many tests call `t.Skipf` when `runtime.GOOS` does not match the expectation. Run the subset aligned with your OS; the tests are guards not cross-platform stubs.
- Frontend tests are purely static (no Jest or Playwright).
  * Use `svelte-kit sync` + `svelte-check` to validate a changed component. Filter the log by file name if you only need one component.

## 4. Style rules – Go / backend

1. **Imports**
   * Standard library imports first (e.g., `context`, `embed`, `log`).
   * Blank line, then external dependencies (Wails, pkg/browser) followed by a blank line.
   * Local packages (`installer/api`, `installer/types`, `installer/utils`) go last.
2. **Formatting**
   * Always run `gofmt`. Tabs are canonical Go indentation.
   * Tables (like the ones in `channel_test.go`) are aligned for readability.
3. **Types & naming**
   * Exported structs and interfaces follow PascalCase (e.g., `App`, `Controller`).
   * Private helpers stay lowercase (e.g., `func newController()` would be unexported if defined).
   * Constants are grouped by responsibility and defined with `iota` (`Stable`, `Canary`, `PTB`).
4. **Error handling**
   * Prefer early returns instead of nested conditionals.
   * When something fails (e.g., `install.InstallBD()`), emit the failure event (`runtime.EventsEmit(action.ctx, "failure")`) and `return` immediately.
   * Swallow errors only when a retry loop is not viable; always log or emit events so the renderer can surface the failure.
5. **Context & runtime**
   * Store the Wails `context.Context` in the struct once (`SetContext`). Do not re-retrieve it; reuse the stored value for dialogs, events, etc.
   * Use the context for runtime calls (e.g., `runtime.MessageDialog`, `runtime.OpenDirectoryDialog`). Guard all runtime returns for errors and treat empty strings as cancellations.
6. **Logging/event protocol**
   * `Controller.Write` redirects fmt output together with emitting a `log` event. Use it to pipe backend state to the UI.
   * Emit semantic events (`log`, `success`, `failure`, `reset`). Do not invent new event names without updating both sides.
7. **Assets & versions**
   * UI assets live in `frontend/build`. Any change to Svelte requires re-running `bun run --bun vite build` so the Go embed rule picks it up.
   * Version string is set via `ldflags` in `main.go`. Compare against GitHub release data using `utils.CompareVersions` to decide if updates are necessary.
8. **Dialogs**
   * Wrap each dialog in a helper that returns string results (`BrowseForDiscord`, `ConfirmAction`, `ShowNotice`). If the call errors or returns an empty string, propagate an empty response and let the renderer treat it as cancel.

## 5. Style rules – Frontend (Svelte + Bun)

### Imports & modules

- Prefer `$lib/` aliases for shared assets, stores, components (see `+page.svelte` or `+layout.svelte`).
- Local constants (e.g., `const installPaths = ...`) live near top-level script scope for clarity.
- Keep third-party imports grouped (Svelte, Wails runtime, helpers) and sorted within each category.

### Script style

- Always use `<script lang="ts">`.
- Default to `const` unless a binding reassigns; Svelte stores should be accessed through `$state` and `$derived` wrappers.
- Introspect derived state with expressions like `$derived(window.navigator.platform.startsWith("Mac"))` (see `+layout.svelte`).
- Inline helper functions (`log`, `succeed`, `reset`) should be declared near the store definitions.

### Components & layout

- Use the shared `<Page>`, `<TextDisplay>`, `<ProgressBar>`, `Titlebar`, and `Footer` components to keep layout consistent.
- Where the UI depends on state (e.g., `currentAction`), render the SVG logic inline inside named slot helpers (`{#snippet icon()}`) instead of spreading props.
- `listenFor` / `unlistenFor` lifecycles should be paired in `onMount`/`onDestroy` blocks to avoid leaking event listeners.

### Styling

- Keep component styles scoped with `<style>` blocks. Avoid global selectors.
- Reference CSS custom properties (`--bg2`, `--bg2-alt`) for colors; these live in `frontend/src/lib/styles/theme.css`.
- Use textures/gradients and layered pseudo-elements (see `.installer-body::after` in `+layout.svelte`) to maintain the polished aesthetic.
- Avoid brittle absolute positioning; prefer flex layouts and `contain: strict` for window surfaces.

### State & events

- Store logs and progress inside `$state` wrappers, and mutate them via `logs.push(...)` or direct assignments when necessary.
- Expose `active` flags for async work, so the page can wire `canGoNext={!active}`.
- Treat `listenFor` payloads as trusted strings but still guard against empty events (e.g., call `message.trim()` before pushing into logs).

### TypeScript organization

- Use `type` aliases for simple shapes (`type Snippet = ...`). Keep them near the components that consume them.
- Always keep Bun config in `frontend/package.json` and `bun.lock` updated in parallel when dependencies change.
- For generated bindings, prefer the Wails TS typing output placed in `frontend/src/lib/wailsjs`. Import from `@api` (generated alias) to keep runtime definitions stable.

## 6. Release & CI context

- The GitHub Workflow `release-pipeline.yml` automates tagging, building, artifact uploads, and draft releases.
- Key steps include merging `development` into `release`, bumping versions via `node` snippet, running `yarn install && yarn dist`, and uploading per-platform installers.
- Keep release notes in sync with real version tags. The pipeline expects tags like `vX.Y.Z` and strips the leading `v` when writing `package.json`.

## 7. Documentation & collaborator expectations

- Reference `README.md` for high-level orientation and `CONTRIBUTING.md` for style guide summaries.
- Emulate the verbose variable names and namespace-heavy imports suggested in `CONTRIBUTING.md` (e.g., avoid destructuring core Node modules if the repo avoids it).
- When documenting new behavior, add sections to `README.md` or a dedicated `docs/` file; keep AGENTS focused on machine-readable guidance.

## 8. Cursor / Copilot compliance

- There are no `.cursor/rules/` or `.cursorrules` directories in this repo; no cursor-specific gating exists.
- There is likewise no `.github/copilot-instructions.md`; default Copilot settings apply.

If future agents add cursor or Copilot rules, append them to this section so every agent reads the new constraints first.

## 9. Security, secrets, and environment

- Do not commit credentials, tokens, `.env` files, or batch scripts that expose signing keys. The release workflow injects secrets via GitHub Actions (`CSC_LINK`, `CSC_KEY_PASSWORD`, `CI_PAT`); local testing should rely on safely stored equivalents.
- Configuration files such as `wails.json`, `frontend/package.json`, and `go.mod` are tracked; avoid adding untracked copies or credentials in the repository tree.
- When discussing or documenting security-sensitive flows (installer auto-update, runtime dialogs that trigger downloads), describe the discovery/confirmation process rather than printing secrets in logs.

## 10. Help & escalation

- If you find conflicting conventions between AGENTS and `CONTRIBUTING.md`, escalate by opening an issue so the maintainers can harmonize the documents.
- Use the Wails community (https://wails.io/support) or BetterDiscord Discord (linked from README) when runtime integration questions exceed local knowledge.
- Document any new tooling (e.g., a new CLI helper) within this file so future agents immediately know how to operate the repo.
