#!/usr/bin/env bash
set -Eeuo pipefail

REPO="${FTN_REPO:-https://github.com/beparykamrul-dev/FTN-AI.git}"
REF="${FTN_REF:-main}"
INSTALL_DIR="${FTN_INSTALL_DIR:-/opt/ftn-ai}"

log(){ printf '[FTN-INSTALL] %s\n' "$*"; }
fail(){ printf '[FTN-INSTALL][ERROR] %s\n' "$*" >&2; exit 1; }

[ "$(id -u)" -eq 0 ] || fail 'Run as root: curl -fsSL <installer> | sudo bash'
command -v apt-get >/dev/null 2>&1 || fail 'Supported host: Debian/Ubuntu with apt-get'
export DEBIAN_FRONTEND=noninteractive

log 'Installing host prerequisites'
apt-get update
apt-get install -y ca-certificates curl git openssl jq

if ! command -v docker >/dev/null 2>&1; then
  log 'Installing Docker Engine'
  curl -fsSL https://get.docker.com | sh
fi
systemctl enable --now docker
command -v docker >/dev/null 2>&1 || fail 'Docker installation failed'
docker compose version >/dev/null 2>&1 || fail 'Docker Compose v2 is required'

if [ -d "$INSTALL_DIR/.git" ]; then
  log "Updating existing checkout: $INSTALL_DIR"
  git -C "$INSTALL_DIR" fetch --prune origin
  git -C "$INSTALL_DIR" checkout -q "$REF"
  git -C "$INSTALL_DIR" reset --hard -q "origin/$REF"
else
  log "Cloning FTN-AI ($REF) -> $INSTALL_DIR"
  mkdir -p "$(dirname "$INSTALL_DIR")"
  git clone --depth 1 --branch "$REF" "$REPO" "$INSTALL_DIR"
fi

cd "$INSTALL_DIR"
chmod +x deploy/one-click/bootstrap.sh
exec deploy/one-click/bootstrap.sh
