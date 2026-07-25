# Tech Stack

- Go 1.24 (module), toolchain go1.25.11 pinned in go.mod; mise pins go 1.25.6 as the tool version.
- Task runner: mise (`.config/mise.toml`, tasks in `tasks/` dir, included via `task_config.includes`).
- Key deps: chi (router), cobra+viper (CLI/config), bbolt (storage), pelletier/go-toml/v2, sirupsen/logrus + CallumKerson/loggerrific (logging), CallumKerson/podcasts (RSS/iTunes feed gen), alfg/mp4 (M4B parsing), gomarkdown/markdown, shopspring/decimal, testify (tests), h2non/baloo.v3 (HTTP handler tests), golang.org/x/time (rate limiting).
- Linting: golangci-lint v2 (`.config/.golangci.yaml`), config uses `linters.default: none` with an explicit enable list — dupl threshold 100, funlen 100 lines/50 statements, gocyclo min-complexity 15, varnamelen enabled (short var names flagged).
- Text tooling: prettier (md/yaml/json), markdownlint-cli2, yamllint, taplo (TOML).
- Docker: Dockerfile + docker-compose.yaml, driven by `mise run docker:run`.
- Release: goreleaser (`.goreleaser.yaml`).
