#!/usr/bin/env bash
#MISE description="Check issues in Go files using golangci-lint."

set -ex

golangci-lint run --config ./.config/.golangci.yaml --fix=false --allow-parallel-runners "$@"
