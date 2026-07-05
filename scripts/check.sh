#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

echo "==> running go vet"
go vet ./...

echo "==> running golangci-lint"
if ! command -v golangci-lint >/dev/null 2>&1; then
  echo "golangci-lint is not installed; use Dockerfile.dev or install it locally" >&2
  exit 1
fi
golangci-lint run ./...
