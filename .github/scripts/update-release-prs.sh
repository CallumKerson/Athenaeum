#!/usr/bin/env bash

set -euxo pipefail

REPO_OWNER="${GH_REPO%%/*}"

gh pr list --label "autorelease: pending" --state open --json number --jq '.[].number' | while read -r pr_number; do
    echo "Processing PR #${pr_number}"
    delay=1
    for attempt in 1 2 3 4 5; do
        if gh pr update-branch "$pr_number" --rebase; then
            echo "Successfully rebased PR #${pr_number}"
            break
        fi
        if [ "$attempt" -lt 5 ]; then
            echo "Attempt $attempt failed, retrying in ${delay}s..."
            sleep "$delay"
            delay=$((delay * 2))
        else
            echo "Could not rebase PR #${pr_number} after 5 attempts"
        fi
    done
done
