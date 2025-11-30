#!/usr/bin/env bash
#MISE description="Check YAML files. Primarily used by `hk` for pre-commit hooks."
#MISE sources=["**/*.yaml"]
#MISE outputs={auto = true}

set -ex

if [ $# -eq 0 ]; then
	yamllint -c .config/.yamllint.yaml .
else
	yamllint -c .config/.yamllint.yaml "$@"
fi
