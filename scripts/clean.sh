#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

BINARY="${BINARY:-activecookie}"
IMAGE="${IMAGE:-activecookie}"
DEV_IMAGE="${DEV_IMAGE:-activecookie-dev}"

echo "==> removing local binary"
rm -f "${BINARY}"

remove_image() {
  local name="$1"
  if docker image inspect "${name}" >/dev/null 2>&1; then
    echo "==> removing docker image ${name}"
    docker rmi "${name}"
  fi
}

if command -v docker >/dev/null 2>&1; then
  remove_image "${IMAGE}"
  remove_image "${DEV_IMAGE}"
  remove_image "${IMAGE}:latest"
  remove_image "${DEV_IMAGE}:latest"
else
  echo "docker not found; skipped docker image cleanup"
fi
