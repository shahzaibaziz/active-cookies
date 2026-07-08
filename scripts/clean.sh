#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

BIN_DIR="${BIN_DIR:-bin}"
BINARY="${BINARY:-${BIN_DIR}/main}"
IMAGE_NAME="${IMAGE_NAME:-activecookie}"
IMAGE="${IMAGE:-${IMAGE_NAME}}"
DEV_IMAGE="${DEV_IMAGE:-activecookie-dev}"

echo "==> removing local binaries"
rm -rf "${BIN_DIR}"
rm -f activecookie

echo "==> removing tmp directory"
rm -rf tmp

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
