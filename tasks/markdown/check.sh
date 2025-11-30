#!/usr/bin/env bash
#MISE description="Check Markdown files for style issues. Primarily used by `hk` for pre-commit hooks."
#MISE sources=["**/*.md"]
#MISE outputs={auto = true}
# shellcheck shell=bash

set -ex

if [ $# -eq 0 ]; then
	markdownlint-cli2 --config ./.config/.markdownlint-cli2.yaml "**/*.md"
else
	markdownlint-cli2 --config ./.config/.markdownlint-cli2.yaml "$@"
fi
