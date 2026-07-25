# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

This project uses [mise](https://mise.jdx.dev/) for task running and tool version management.
Run `mise tasks` to see all available tasks.

### Building and Testing

- **Compile**: `mise run golang:compile` - Compiles the binary for the current OS/architecture
- **Test**: `mise run test` or `go test ./...` - Run all tests
- **Fix**: `mise run fix-all` - Runs all formatters and auto-fixable linters via [hk](https://hk.jdx.dev/)
- **Check**: `mise run check-all` - Runs all linters without fixing
- **CI**: `mise run ci` - Runs `check-pr` and `test` (used in CI pipeline)

Run `mise run fix-all` and `mise run test` before committing.
Linters are also wired to git hooks through `.config/hk.pkl`: fast formatters on pre-commit, golangci-lint as well on pre-push.

## Architecture Overview

Athenaeum is a static site generator.
`athenaeum build` scans a library of `.m4b` audiobooks and writes a tree of podcast feeds, which any web server can serve.
There is no server process, no database and no request path.

### Packages

- `pkg/audiobooks` - The `Audiobook` entity, the `Genre` enum and description formatting
- `internal/scan` - Walks the media root, reads `.toml` metadata and `.m4b` durations, with a duration cache
- `internal/feed` - Renders a list of audiobooks as an iTunes-compatible RSS 2.0 feed
- `internal/site` - Decides which feeds exist at which paths (`plan.go`) and writes the tree (`build.go`)
- `internal/fsutil` - Atomic and write-if-changed file helpers
- `cmd/athenaeum` - Cobra commands and config loading
- `static`, `templates` - Embedded assets and the index page template

### Data Flow

1. `scan.Library` walks the media root for `.m4b` + `.toml` file pairs
2. `site.Plan` groups the books into feeds — all, by author, narrator, genre and tag — and assigns each a set of paths
3. `site.Build` renders each feed with `feed.Renderer` and writes it to every path in its plan
4. Stale files from previous builds are swept, and Overcast is optionally notified

### Constraints that must not be broken

- **A feed item's `<guid>` is its enclosure URL.**
  Changing how either is built makes every subscriber re-download the whole library.
  `feed.mediaURL` is the only place URLs are constructed, and file paths are assigned to `url.URL.Path` so that `?`, `#` and `%` in filenames are escaped rather than parsed.
- **Rebuilds must be idempotent.**
  Rendering is deterministic — no timestamps in the output — and `fsutil.WriteIfChanged` skips unchanged files so their mtimes, and therefore the web server's ETags, stay stable.
- **One book can be reachable under several spellings.**
  Names that normalise identically (`V.E. Schwab` / `V. E. Schwab`) each need their own feed path, and every genre gets a feed even when empty.
- **The build refuses to write into a directory it does not own**, recognised by the `.athenaeum-site` marker file,
  because it sweeps files it did not write.

### Configuration

- TOML config at `~/.config/athenaeum/athenaeum.toml` (XDG paths resolved by hand — `os.UserConfigDir()` returns
  `~/Library/...` on macOS and ignores `XDG_CONFIG_HOME`)
- Duration cache at `~/.cache/athenaeum/m4b.json`
- Required: `Host` (external URL), `Media.Root` (library path) and `Site.Root` (output path)
- Audiobook layout: `$MEDIA_ROOT/Author/Audiobook/Audiobook.m4b` + `Audiobook.toml`

### Testing

- Unit tests use testify/assert and testify/require
- `internal/feed/testdata` holds feeds captured from the previous server implementation; they are matched byte for byte
  as a parity guarantee for existing subscribers
- `internal/site/testdata/golden_manifest.txt` hashes the build plan rather than the output tree, because a case insensitive filesystem would fold distinct feed paths together and make a golden tree unportable.
  Regenerate with `go test ./internal/site -update`
- Test audiobooks live in the root `testdata/` directory

## Import Organization

- Standard library imports first
- Third-party imports second
- Local imports last with company prefix `github.com/CallumKerson`
- Project imports use full path: `github.com/CallumKerson/Athenaeum/internal/...`

## Code Patterns

- Prefer direct functions; structs are for data, not for behaviour that a function would express as well
- Comments explain why a thing is the way it is, not what the line does
