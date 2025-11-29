# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

This project uses [mise](https://mise.jdx.dev/) for task running and tool
version management. Run `mise tasks` to see all available tasks.

### Building and Testing

- **Compile only**: `mise run golang:compile` - Compiles binary for current OS/architecture
- **Linux compile**: `mise run golang:linux-compile` - Cross-compiles for Linux deployment
- **Test**: `mise run test` or `go test ./...` - Run all tests
- **Format**: `mise run format` - Runs all formatters (Go and text files)
- **Lint**: `mise run lint` - Runs all linters (Go and text files)
- **Pre-commit**: `mise run pre-commit` - Runs format, lint, and test (recommended before committing)
- **CI**: `mise run ci` - Runs lint and test (used in CI pipeline)

### Docker

- **Docker build and run**: `mise run docker:run` - Runs pre-commit checks, builds Linux binary, and starts with docker-compose

### Task Structure

Tasks are organized in the `tasks/` directory by category:

- `golang/` - Go-specific tasks (compile, test, lint, format, mod-tidy)
- `text/` - Text file tasks (markdown/YAML formatting and linting)
- `docker/` - Container tasks

## Architecture Overview

Athenaeum is an audiobook server that provides podcast feeds, built with a
layered architecture following clean architecture principles.

### Core Architecture

- **Domain Layer** (`pkg/`): Core entities (Audiobook, Genre, etc.) and domain interfaces
- **Service Layer** (`internal/*/service/`): Business logic for media scanning, audiobook management, and podcast generation
- **Adapter Layer** (`internal/adapters/`): External integrations (BoltDB storage, M4B metadata extraction, logging)
- **Transport Layer** (`internal/transport/http/`): HTTP handlers, middleware, and REST endpoints
- **CLI Layer** (`cmd/athenaeum/`): Command-line interface and service composition

### Key Services

- **Media Service**: Scans filesystem for `.m4b` audiobook files and their `.toml` metadata companions
- **Audiobooks Service**: Orchestrates scanning, storage (BoltDB), and third-party notifications (Overcast)
- **Podcasts Service**: Generates RSS 2.0 feeds with iTunes compatibility from audiobook collections

### Data Flow

1. Media Service scans filesystem for `.m4b` + `.toml` file pairs
2. Audiobooks Service extracts metadata and stores in BoltDB via Bolt adapter
3. Podcasts Service generates filtered RSS feeds (by author, genre, narrator, etc.)
4. HTTP Transport serves feeds at `/podcast/feed.rss` and media files at `/media/*`

### Storage

- **Primary**: BoltDB embedded database (JSON serialization of audiobook entities)
- **Caching**: In-memory LRU cache with TTL for HTTP responses
- **Media**: Direct filesystem serving of `.m4b` files

### Configuration

- YAML/TOML config file (default: `~/.athenaeum/config.yaml`)
- Required: `Host` (external URL) and `Media.Root` (audiobook directory path)
- Audiobook layout: `$MEDIA_ROOT/Author/Audiobook/Audiobook.m4b` + `Audiobook.toml`

### Testing

- Unit tests use testify/assert and testify/require
- Integration tests with test audiobooks in `testdata/`
- HTTP handler tests use `gopkg.in/h2non/baloo.v3` for API testing

## Import Organization

- Standard library imports first
- Third-party imports second
- Local imports last with company prefix `github.com/CallumKerson`
- Project imports use full path: `github.com/CallumKerson/Athenaeum/internal/...`

## Code Patterns

- **Options Pattern**: Services accept `opts ...Option` for configuration
- **Interface Segregation**: Clear interfaces for storage (`AudiobookStore`), metadata reading (`M4BMetadataReader`), logging
- **Filter Pattern**: Functional composition for audiobook queries (`AuthorFilter`, `GenreFilter`, etc.)
- **Dependency Injection**: Services depend on interfaces, configured in `cmd/athenaeum/cmd.go`
