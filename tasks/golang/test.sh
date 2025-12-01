#!/usr/bin/env bash
#MISE sources=["go.mod", "go.sum", "**/*.go", "**/testdata/**/*", "**/*.gohtml"]
#MISE outputs={auto = true}

set -ex

go test ./...
