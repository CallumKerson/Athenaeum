#!/usr/bin/env bash
#MISE description="Check TOML file formatting. Primarily used by `hk` for pre-commit hooks."
#MISE sources=["**/*.toml"]
#MISE outputs={auto = true}
# shellcheck shell=bash

set -ex

taplo format --check --config ./.config/taplo.toml "$@"
