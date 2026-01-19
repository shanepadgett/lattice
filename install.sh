#!/usr/bin/env bash

set -euo pipefail

REPO_OWNER="shanepadgett"
REPO_NAME="lattice"
PROJECT_NAME="lcss"
BINARY_NAME="lcss"
DEFAULT_OUT_DIR="$HOME/.local/bin"
CHECKSUMS_FILE="checksums.txt"

VERSION=""
OUT_DIR="$DEFAULT_OUT_DIR"

usage() {
  cat <<'EOF'
Usage:
  install.sh [--version vX.Y.Z] [--out-dir <path>] [--help]

Options:
  --version   Release tag to install (defaults to latest)
  --out-dir   Install directory (defaults to ~/.local/bin)
  --help      Show this help
EOF
}

log() {
  printf "[lcss] %s\n" "$1"
}

fail() {
  printf "[lcss] error: %s\n" "$1" >&2
  exit 1
}

ensure_command() {
  local name="$1"
  if ! command -v "$name" >/dev/null 2>&1; then
    fail "Missing required command: $name"
  fi
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --version)
      if [[ $# -lt 2 ]]; then
        fail "--version requires a value"
      fi
      VERSION="$2"
      shift 2
      ;;
    --out-dir)
      if [[ $# -lt 2 ]]; then
        fail "--out-dir requires a value"
      fi
      OUT_DIR="$2"
      shift 2
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      fail "Unknown argument: $1"
      ;;
  esac
done

ensure_command curl
ensure_command tar
ensure_command sha256sum

resolve_version() {
  if [[ -n "$VERSION" ]]; then
    printf "%s" "$VERSION"
    return
  fi

  local latest_url="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest"
  local latest
  latest=$(curl -fsSL "$latest_url" | awk -F '"' '/"tag_name"/ { print $4; exit }')
  if [[ -z "$latest" ]]; then
    fail "Unable to resolve latest release tag"
  fi
  printf "%s" "$latest"
}

normalize_arch() {
  local arch="$1"
  case "$arch" in
    x86_64|amd64)
      printf "x86_64"
      ;;
    arm64|aarch64)
      printf "arm64"
      ;;
    *)
      fail "Unsupported architecture: $arch"
      ;;
  esac
}

normalize_os() {
  local os="$1"
  case "$os" in
    Linux)
      printf "Linux"
      ;;
    Darwin)
      printf "Darwin"
      ;;
    *)
      fail "Unsupported OS: $os"
      ;;
  esac
}

main() {
  local version
  local os
  local arch
  local asset
  local base_url
  local tmpdir
  local archive_path
  local checksums_path
  local expected_checksum

  version=$(resolve_version)
  os=$(normalize_os "$(uname -s)")
  arch=$(normalize_arch "$(uname -m)")
  asset="${PROJECT_NAME}_${os}_${arch}.tar.gz"
  base_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${version}"

  log "Installing ${BINARY_NAME} ${version} to ${OUT_DIR}"

  tmpdir=$(mktemp -d)
  trap 'rm -rf "$tmpdir"' EXIT

  archive_path="$tmpdir/$asset"
  checksums_path="$tmpdir/$CHECKSUMS_FILE"

  curl -fsSL "$base_url/$asset" -o "$archive_path"
  curl -fsSL "$base_url/$CHECKSUMS_FILE" -o "$checksums_path"

  expected_checksum=$(awk -v file="$asset" '$2 == file { print $1 }' "$checksums_path")
  if [[ -z "$expected_checksum" ]]; then
    fail "Checksum not found for $asset"
  fi

  echo "${expected_checksum}  ${archive_path}" | sha256sum -c - >/dev/null 2>&1 || fail "Checksum verification failed"

  mkdir -p "$OUT_DIR"
  tar -xzf "$archive_path" -C "$tmpdir"
  if [[ ! -f "$tmpdir/$BINARY_NAME" ]]; then
    fail "Binary not found in archive"
  fi

  mv "$tmpdir/$BINARY_NAME" "$OUT_DIR/$BINARY_NAME"
  chmod +x "$OUT_DIR/$BINARY_NAME"

  log "Installed ${OUT_DIR}/${BINARY_NAME}"

  if [[ ":$PATH:" != *":${OUT_DIR}:"* ]]; then
    log "Add ${OUT_DIR} to PATH to use '${BINARY_NAME}'"
  fi
}

main
