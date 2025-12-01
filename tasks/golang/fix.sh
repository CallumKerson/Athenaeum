#!/usr/bin/env bash
#MISE description="Fix issues in Go files using golangci-lint."

set -ex

golangci-lint run --config ./.config/.golangci.yaml --fix=true --allow-parallel-runners "$@"
