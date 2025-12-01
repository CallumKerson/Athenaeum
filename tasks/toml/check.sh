#!/usr/bin/env bash
#MISE description="Check TOML files. Primarily used by `hk` for pre-commit hooks."
#MISE sources=["**/*.toml"]
#MISE outputs={auto = true}
# shellcheck shell=bash

set -ex

taplo lint --config ./.config/taplo.toml "$@"
