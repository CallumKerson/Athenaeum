#!/usr/bin/env bash
#MISE description="Compiles for Linux"
#MISE sources=["**/*.go", "go.mod", "go.sum"]
#MISE outputs={auto = true}

set -ex

GOOS=linux go build -ldflags "-s -w" -o . ./...
