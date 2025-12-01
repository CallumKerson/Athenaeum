#!/usr/bin/env bash
#MISE description="Fix formatting issues in Markdown, YAML, JSON, and HTML files."
#MISE sources=["**/*.md", "**/*.yaml", "**/*.json", "**/*.html"]
#MISE outputs={auto = true}

set -ex

prettier --write --ignore-path=.gitignore --ignore-path=.config/.prettierignore "${@:-.}"
