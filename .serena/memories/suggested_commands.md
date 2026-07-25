# Suggested Commands

All via mise (`mise tasks` lists everything).
Task names are `namespace:task`, invoked as `mise run <name>`.

- `mise run golang:compile` — build binary for current OS/arch (also embeds git-derived version/commit/date via ldflags)
- `mise run golang:linux-compile` — cross-compile for Linux (used before docker build)
- `mise run test` (alias, depends on `golang:test`) or `go test ./...` — run all tests
- `mise run format` — runs all `*:format` tasks (golang:format runs `golangci-lint run --fix`, depends on mod-tidy; text:format runs prettier)
- `mise run lint` — runs all `*:lint` tasks (golangci-lint, mod-tidy-check, prettier --check, markdownlint-cli2, yamllint); waits for format
- `mise run ci` — lint + test (what CI runs)
- `mise run pre-commit` — format + lint + test (run this before committing, per user's global CLAUDE.md instructions)
- `mise run docker:run` — pre-commit + linux-compile + `docker compose up --build`
- `go mod tidy` — required when go.mod/go.sum drift; `mise run golang:mod-tidy-check` enforces this in CI via `go mod tidy -diff`

Darwin-specific notes: none beyond mise itself; standard unix tools behave normally on this repo.
