# Core Memories

Athenaeum: audiobook server serving podcast RSS feeds.
Go module `github.com/CallumKerson/Athenaeum`.

## Source map

- `cmd/athenaeum/` — CLI entrypoint, viper/cobra config, service composition/DI (`cmd.go`)
- `pkg/` — domain layer: `audiobooks` (Audiobook, Genre entities + description formatting), `caching/response`, `client` (HTTP client wrapper), `m4b` (metadata struct)
- `internal/media/service/` — scans filesystem for `.m4b` + `.toml` pairs
- `internal/audiobooks/service/` — orchestrates scanning, storage, notifications; filter pattern for queries (`filters.go`)
- `internal/podcasts/service/` — builds RSS 2.0/iTunes feeds from audiobooks
- `internal/adapters/` — external integrations: `bolt` (BoltDB store), `alfgmp4` (M4B metadata reader), `logrus` (logging adapter)
- `internal/notifiers/overcast/` — pushes updates to Overcast podcast app
- `internal/memcache/` — in-memory LRU+TTL cache for HTTP responses
- `internal/transport/http/` — chi router, handlers, middleware (incl. caching middleware)
- `internal/testing/` — shared test helpers (`dataloader`, `testbooks`)

Full architecture narrative and data flow already documented in project root `CLAUDE.md` — read that first for the conceptual model; this memory is for locating things.

For commands/build tooling: `mem:suggested_commands`.
For task-completion checklist: `mem:task_completion`.
For stack/versions: `mem:tech_stack`.
For code style: `mem:conventions`.
