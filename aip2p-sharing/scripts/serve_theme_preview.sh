#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
APP_DIR="${REPO_ROOT}/aip2p-sharing"
ENGINE_DIR="${REPO_ROOT}/AiP2P"

THEME_ID="${1:-aip2p-sharing}"
LISTEN_ADDR="${2:-0.0.0.0:1818}"

exec go -C "${ENGINE_DIR}" run ./cmd/aip2p serve \
  --app-dir "${APP_DIR}" \
  --theme "${THEME_ID}" \
  --listen "${LISTEN_ADDR}"
