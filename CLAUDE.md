# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building and Testing

- **Build**: `task build` - Builds the project, runs all linting, formatting, and compiles binary
- **Compile only**: `task compile` - Compiles binary for current OS/architecture
- **Linux compile**: `task linux-compile` - Cross-compiles for Linux deployment
- **Test**: `task test` or `go test ./...` - Run all tests
- **Format and Lint**: `task tidy` - Runs formatting, import sorting, linting, and pre-commit hooks

### Docker

- **Docker build and run**: `task docker` - Builds Linux binary and runs with docker-compose

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
