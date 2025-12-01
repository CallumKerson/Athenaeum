#!/usr/bin/env bash
#MISE description="Format Markdown files. Primarily used by `hk` for pre-commit hooks."
#MISE sources=["**/*.md"]
#MISE outputs={auto = true}

set -e

set -ex

markdownlint-cli2 --fix --config ./.config/.markdownlint-cli2.yaml "${@:-.}"
