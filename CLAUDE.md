# fizzy-cli

A Go CLI for the [Fizzy](https://fizzy.do) API. Fizzy is a kanban-style project management SaaS.

The binary is `fizzy` and is distributed via npm (`fizzy-cli` package); `scripts/postinstall.js` downloads the matching Go binary from the GitHub release on install. Local installs use `make install` to drop the binary into `$GOBIN`.

## Architecture

- `main.go` is a thin entrypoint that calls `cmd.Execute()`.
- `cmd/` holds every CLI command, one command per file using [`spf13/cobra`](https://github.com/spf13/cobra). Naming is `<resource>_<action>.go` (e.g. `board_create.go`, `card_assign.go`); the resource's parent command lives in `<resource>.go` (e.g. `board.go`, `card.go`). Tests sit alongside as `<file>_test.go`.
- `internal/app` exposes the `App` struct (`Client *fizzy.Client` + `Config *config.Config`). `App.New` is invoked once in `rootCmd.PersistentPreRun` and stashed in the cobra command context; commands recover it via `app.FromContext(cmd.Context())`.
- `internal/config` reads/writes `~/.config/fizzy-cli/config.json` (`SelectedAccount`, `SelectedBoard`, `CurrentUserID`). File mode is `0600`.
- `internal/ui` holds presentation helpers — one file per resource (`card_show.go`, `board_list.go`, …) plus `format.go` for shared helpers like `DisplayID` and `FormatTime`. Styling uses `lipgloss`.
- `internal/colors` is for color helpers.
- `internal/testutil` exposes `NewTestClient(baseURL, accountSlug, boardID, accessToken)` for wiring an `httptest.Server` into a `fizzy.Client`.

The actual API layer lives in [`github.com/rogeriopvl/fizzy-go`](https://github.com/rogeriopvl/fizzy-go) — **do not** reimplement endpoints here. Call the typed methods on `*fizzy.Client` and consume the structs it returns.

## API specs

The official Fizzy API specs live in `docs/api/`:

- `docs/api/README.md` — index/overview of the API
- `docs/api/sections/*.md` — one file per resource (boards, cards, columns, comments, identity, notifications, reactions, steps, tags, users, webhooks, pins, account, activities, authentication, exports, rich_text)

This CLI does not talk to the Fizzy API directly — `fizzy-go` does. So you generally don't need to read the specs to implement a command; rely on the typed methods and structs `fizzy-go` already exposes. The specs are kept here as reference for when you *do* need them: discovering which endpoints exist, checking request/response shapes when adding a new feature, or maintaining parity as the API evolves.

The specs are synced from the upstream `basecamp/fizzy` repo via `make sync-api-spec`. Run that if the local copy may be stale.

## Conventions

- New subcommands follow the existing pattern:
  1. Define `var fooBarCmd = &cobra.Command{...}` in `cmd/foo_bar.go`.
  2. The `Run` func delegates to a `handle<Action><Resource>` function (e.g. `handleCreateBoard`, `handleShowCard`) so it stays unit-testable. The handler takes `cmd *cobra.Command` plus any positional args destructured by name (e.g. `handleAssignCard(cmd *cobra.Command, cardNumber, userID string)`), not `args []string`.
  3. Inside the handler, recover the app via `app.FromContext(cmd.Context())` and bail with `"API client not available"` if `a == nil || a.Client == nil`.
  4. Read flags via `cmd.Flags().GetX(...)`. Use `MarkFlagRequired` for required flags. Use `Flags().Changed(...)` when distinguishing "unset" from "zero value" matters (e.g. partial updates).
  5. Wrap API errors as `fmt.Errorf("doing thing: %w", err)`.
  6. Write user-facing output to `cmd.OutOrStdout()` and errors to `cmd.OutOrStderr()` — never `fmt.Println`/`os.Stdout` directly, so tests can capture output.
  7. Register the command in `init()` by attaching it to its parent (e.g. `boardCmd.AddCommand(boardCreateCmd)`); the root-level resource command attaches itself with `rootCmd.AddCommand(...)`.
- Display logic belongs in `internal/ui`, not in `cmd/`. Commands fetch data and call a `ui.Display…` function.
- Every new command must ship with tests in the matching `_test.go`. Tests spin up an `httptest.NewServer`, build a client with `testutil.NewTestClient`, attach it via `(&app.App{Client: client}).ToContext(...)`, parse flags with `cmd.ParseFlags(...)`, and call the handler directly. Cover at minimum: success, API error, and "no client" paths.
- Don't write obvious comments. Only add a comment when the code itself needs clarification.
- Performance matters: prefer passing large structs by pointer, avoid unnecessary allocations.
- When in doubt, mirror the closest existing command/test pair.

## Build and run

All dev commands are in the `Makefile`:

- `make install` — build with the version from `package.json` baked in via ldflags and drop the binary into `$GOBIN` (or `$(go env GOPATH)/bin`). Use this for manual testing/debugging — then run `fizzy` normally.
- `make build` / `make build-dev` — build into `bin/fizzy` (release vs. debug ldflags).
- `make build-all` — cross-compile darwin/linux/windows for amd64/arm64 into `bin/`.
- `make test` — run tests via `gotestsum` (install with `make dev-tools`).
- `make run` — `go run .`.
- `make sync-api-spec` — refresh `docs/api/` from upstream.

The CLI version comes from `package.json` (parsed by the Makefile) and is injected into `cmd.Version` at build time via `-ldflags`. Bump `package.json` `version` when releasing.

Authentication: the binary reads `FIZZY_ACCESS_TOKEN` from the environment. Without it, `App.New` still returns an `App` with `Config` populated but `Client == nil`, so commands must guard for that.

## Releases

Releases are SemVer, driven from [Conventional Commits](https://www.conventionalcommits.org/) since the previous `vX.Y.Z` tag (`feat:`, `fix:`, `chore:`, `feat(scope): ...`, etc.). Each feature gets its own commit — don't bundle unrelated changes.

The release flow:

1. Working tree clean, `make test` green.
2. Pick the new version from the commits since the last tag:
   - `BREAKING CHANGE` footer or `!` after the type → major
   - `feat:` → minor
   - `fix:` / `chore:` / `docs:` / other → patch
3. Bump `version` in `package.json` (the Makefile reads it for ldflags).
4. Update `CHANGELOG.md` with a section for the new version.
5. Commit as `chore(release): vX.Y.Z`.
6. Tag: `git tag -a vX.Y.Z -m "vX.Y.Z"`.
7. Push commit and tag together: `git push --follow-tags`.
8. Create the GitHub release (`gh release create vX.Y.Z --notes-file <file>`) — no need to attach binaries by hand. The `release.yml` workflow listens on `release: created`, runs `make build-all`, and uploads `bin/fizzy-*` to the release. `scripts/postinstall.js` then fetches them on `npm install`.
9. Publish to npm so the new version is installable.

Do not add `Co-Authored-By` trailers (or any other AI-attribution trailer) to commit messages.
