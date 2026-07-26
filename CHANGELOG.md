# Changelog

## [2.0.1](https://github.com/CallumKerson/Athenaeum/compare/v2.0.0...v2.0.1) (2026-07-26)


### Bug Fixes

* **deps:** update github.com/gomarkdown/markdown digest to 8435af3 ([#236](https://github.com/CallumKerson/Athenaeum/issues/236)) ([b646268](https://github.com/CallumKerson/Athenaeum/commit/b64626871af688235b9c76625cca0ef87c38d82b))
* **deps:** update module github.com/pelletier/go-toml/v2 to v2.4.3 ([#228](https://github.com/CallumKerson/Athenaeum/issues/228)) ([29a4af6](https://github.com/CallumKerson/Athenaeum/commit/29a4af623e059eaa5a9fa8c2aee77b736cbf4a31))

## [2.0.0](https://github.com/CallumKerson/Athenaeum/compare/v1.15.0...v2.0.0) (2026-07-25)


### ⚠ BREAKING CHANGES

* the server is removed. `athenaeum run` and `athenaeum update`, the `/update`, `/health`, `/ready` and `/version` endpoints, and the Docker image no longer exist. Configuration moves from YAML to TOML at `~/.config/athenaeum/athenaeum.toml`, and the Homebrew cask no longer installs a launchd service — run `brew services stop athenaeum` once after upgrading.

### Features

* replace the server with a static feed generator ([#224](https://github.com/CallumKerson/Athenaeum/issues/224)) ([2e34743](https://github.com/CallumKerson/Athenaeum/commit/2e34743c4af59e3e790718f65e2f5793c7952b4a))


### Bug Fixes

* **deps:** update github.com/gomarkdown/markdown digest to 37c66b8 ([#182](https://github.com/CallumKerson/Athenaeum/issues/182)) ([cbb2dc1](https://github.com/CallumKerson/Athenaeum/commit/cbb2dc15fed4ff76e3f5aa2a31e200b4879355fc))
* **deps:** update module github.com/go-chi/chi/v5 to v5.2.4 ([#168](https://github.com/CallumKerson/Athenaeum/issues/168)) ([aa090b0](https://github.com/CallumKerson/Athenaeum/commit/aa090b06bb4cc581ed6c96c40166d7499acc1208))
* **deps:** update module github.com/go-chi/chi/v5 to v5.2.5 ([#176](https://github.com/CallumKerson/Athenaeum/issues/176)) ([9c0974d](https://github.com/CallumKerson/Athenaeum/commit/9c0974d2b1383b442077745415a1a8f9a93f6a25))
* **deps:** update module github.com/pelletier/go-toml/v2 to v2.3.0 ([#200](https://github.com/CallumKerson/Athenaeum/issues/200)) ([b563413](https://github.com/CallumKerson/Athenaeum/commit/b563413fdaddf5b1786ce969422ff6d2686f7dac))
* **deps:** update module github.com/pelletier/go-toml/v2 to v2.3.1 ([#203](https://github.com/CallumKerson/Athenaeum/issues/203)) ([4582071](https://github.com/CallumKerson/Athenaeum/commit/4582071ba614d14708598dd70b51fa6114ff6470))
* **deps:** update module github.com/pelletier/go-toml/v2 to v2.4.0 ([#217](https://github.com/CallumKerson/Athenaeum/issues/217)) ([22f5e76](https://github.com/CallumKerson/Athenaeum/commit/22f5e763da93c7d6afdfacf482752c17928cb67c))
* **deps:** update module github.com/sirupsen/logrus to v1.9.4 ([#165](https://github.com/CallumKerson/Athenaeum/issues/165)) ([67db165](https://github.com/CallumKerson/Athenaeum/commit/67db165bca2e125a000ce7cccaec7260d552856b))
* **deps:** update module github.com/spf13/cobra to v1.10.2 ([#159](https://github.com/CallumKerson/Athenaeum/issues/159)) ([50e23f1](https://github.com/CallumKerson/Athenaeum/commit/50e23f1c2dd8c96a68c648eac3feeb915d8a1875))
* update release to use homebrew casks ([#154](https://github.com/CallumKerson/Athenaeum/issues/154)) ([4a9aa78](https://github.com/CallumKerson/Athenaeum/commit/4a9aa780d87ab80a2585c0c487bba322f3923c02))

## [1.15.0](https://github.com/CallumKerson/Athenaeum/compare/v1.14.30...v1.15.0) (2025-11-29)


### Features

* add support for episode images in podcast feeds ([#145](https://github.com/CallumKerson/Athenaeum/issues/145)) ([e3b2d16](https://github.com/CallumKerson/Athenaeum/commit/e3b2d16f15117ea38b5b70137fe02ac7e974afdb))

## [1.14.30](https://github.com/CallumKerson/Athenaeum/compare/v1.14.29...v1.14.30) (2025-11-29)


### Bug Fixes

* fix release process ([#146](https://github.com/CallumKerson/Athenaeum/issues/146)) ([84af719](https://github.com/CallumKerson/Athenaeum/commit/84af7198bfa72746f20487507bd8e9657fd0d80d))

## [1.14.29](https://github.com/CallumKerson/Athenaeum/compare/v1.14.28...v1.14.29) (2025-11-29)


### Bug Fixes

* **deps:** update github.com/gomarkdown/markdown digest to 2e2c118 ([#107](https://github.com/CallumKerson/Athenaeum/issues/107)) ([9cf6f0c](https://github.com/CallumKerson/Athenaeum/commit/9cf6f0cc79fb13bc32e52c0d9df462fd7e5ac351))
* **deps:** update module github.com/callumkerson/podcasts to v1 ([#104](https://github.com/CallumKerson/Athenaeum/issues/104)) ([dfbfa0f](https://github.com/CallumKerson/Athenaeum/commit/dfbfa0f74ce661ece701e9ae0bcda54e8eab6142))
