#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

echo "==> formatting Go sources"
go fmt ./...
