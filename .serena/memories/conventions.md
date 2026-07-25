# Conventions

Base conventions (import order, options pattern, interface segregation, filter pattern, DI in cmd.go) are documented in root `CLAUDE.md` — read that first.

Additional non-obvious constraints enforced by golangci-lint (`.config/.golangci.yaml`), not just style preference:

- `gochecknoinits` — no `init()` functions.
- `err113` — don't use dynamic `errors.New`/`fmt.Errorf` inline for sentinel-style comparisons; define wrapped/sentinel errors.
- `varnamelen` enabled — very short variable names get flagged outside tiny scopes; prefer descriptive names except in short loops.
- `dupl` threshold 100 — near-duplicate code blocks >100 tokens will fail lint; extract shared helpers instead of copy-pasting across adapters/services.
- `funlen` — 100 lines / 50 statements max per function.
- `gocyclo` — min-complexity 15 triggers failure; keep branching low.
- `exhaustive` — switch statements over defined enums (e.g. Genre-like types) must handle all cases or explicitly default.
- Per-package `Option` functions live in a sibling `opts.go` file (e.g. `internal/audiobooks/service/opts.go`, `internal/media/service/opts.go`), separate from the main service file — follow this file-split convention when adding new options.
