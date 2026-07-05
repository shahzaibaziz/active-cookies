#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

echo "==> running tests"
go test ./... -v
