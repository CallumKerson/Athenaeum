#!/usr/bin/env bash
#MISE description="Check issues in Markdown, YAML, JSON, and HTML files."
#MISE sources=["**/*.md", "**/*.yaml", "**/*.json", "**/*.html"]
#MISE outputs={auto = true}

set -e

if [ $# -eq 0 ]; then
	set -x
	prettier '**/*.{md,yaml,json,html}' --check --ignore-path=.gitignore --ignore-path=.config/.prettierignore
else
	set -x
	prettier --check --ignore-path=.gitignore --ignore-path=.config/.prettierignore "$@"
fi
