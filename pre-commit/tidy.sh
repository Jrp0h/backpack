#!/usr/bin/bash

set -euo pipefail

go mod tidy
git add go.mod