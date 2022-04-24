#!/usr/bin/bash

set -euo pipefail

output=$(gofmt -s -w -l "$@")

if [ "$output" != "" ]; then
    git add "$@"
fi