#!/usr/bin/env bash
#MISE description="Format Markdown files. Primarily used by `hk` for pre-commit hooks."
#MISE sources=["**/*.md"]
#MISE outputs={auto = true}

set -e

if [ $# -eq 0 ]; then
	set -x
	markdownlint-cli2 --config ./.config/.markdownlint-cli2.yaml --fix "**/*.md"
else
	set -x
	markdownlint-cli2 --config ./.config/.markdownlint-cli2.yaml --fix "$@"
fi
