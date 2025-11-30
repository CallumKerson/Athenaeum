#!/usr/bin/env bash
#MISE description="Fix formatting issues in Markdown, YAML, JSON, and HTML files."
#MISE sources=["**/*.md", "**/*.yaml", "**/*.json", "**/*.html"]
#MISE outputs={auto = true}

set -e

if [ $# -eq 0 ]; then
	set -x
	prettier '**/*.{md,yaml,json,html}' --write --ignore-path=.gitignore --ignore-path=.config/.prettierignore
else
	set -x
	prettier --write --ignore-path=.gitignore --ignore-path=.config/.prettierignore "$@"
fi
