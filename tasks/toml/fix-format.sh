#!/usr/bin/env bash
#MISE description="Format TOML files. Primarily used by `hk` for pre-commit hooks."
#MISE sources=["**/*.toml"]
#MISE outputs={auto = true}
# shellcheck shell=bash

set -ex

taplo format --config ./.config/taplo.toml "$@"
