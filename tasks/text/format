#!/usr/bin/env bash
#MISE sources=["**/*.md", "**/*.yaml", "**/*.json"]
#MISE outputs={auto = true}

set -ex

prettier '**/*.{md,yaml,json}' --write --ignore-path=.gitignore --ignore-path=.config/.prettierignore
