#!/usr/bin/bash

set -euo pipefail

./project gen-docs
git add docs/md/* docs/man/*
