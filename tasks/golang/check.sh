#!/usr/bin/env bash
#MISE description="Check issues in Go files using golangci-lint."
#MISE sources=["**/*.go"]
#MISE outputs={auto = true}

set -e

# Count total Go files in the repo
total_go_files=$(find . -name "*.go" -type f | wc -l | tr -d ' ')

# If no arguments or all/more files passed, run full check
# Otherwise use golangci-lint's hook mode (check new changes only)
if [[ $# -eq 0 ]] || [[ $# -ge $total_go_files ]]; then
	set -x
	golangci-lint run --config ./.config/.golangci.yaml --fix=false --allow-parallel-runners
else
	set -x
	golangci-lint run --config ./.config/.golangci.yaml --fix=false --allow-parallel-runners --new-from-rev HEAD
fi
