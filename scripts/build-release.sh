#!/usr/bin/env bash

set -euo pipefail

APP_NAME="${APP_NAME:-pit}"
VERSION="${VERSION:-}"
DIST_DIR="${DIST_DIR:-dist}"
PACKAGE="${PACKAGE:-.}"

TARGETS=(
  "darwin/amd64"
  "darwin/arm64"
  "linux/amd64"
  "linux/arm64"
  "windows/amd64"
)

usage() {
  printf 'Usage: VERSION=0.1.0 %s\n' "$0"
  printf '       %s 0.1.0\n' "$0"
}

if [[ $# -gt 1 ]]; then
  usage >&2
  exit 2
fi

if [[ $# -eq 1 ]]; then
  VERSION="$1"
fi

if [[ -z "$VERSION" ]]; then
  VERSION="$(git describe --tags --exact-match 2>/dev/null || true)"
fi

if [[ -z "$VERSION" ]]; then
  printf 'error: version is required. Pass it as an argument or VERSION env var.\n\n' >&2
  usage >&2
  exit 2
fi

VERSION="${VERSION#v}"

for cmd in go tar zip shasum; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    printf 'error: required command not found: %s\n' "$cmd" >&2
    exit 1
  fi
done

printf 'Starting release build for %s v%s\n' "$APP_NAME" "$VERSION"

rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

export CGO_ENABLED=0

for target in "${TARGETS[@]}"; do
  IFS='/' read -r GOOS GOARCH <<< "$target"

  artifact_base="${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}"
  binary_name="$APP_NAME"
  if [[ "$GOOS" == "windows" ]]; then
    binary_name="${binary_name}.exe"
  fi

  build_dir="${DIST_DIR}/${artifact_base}"
  mkdir -p "$build_dir"

  printf 'Building %s/%s\n' "$GOOS" "$GOARCH"
  GOOS="$GOOS" GOARCH="$GOARCH" go build \
    -trimpath \
    -ldflags="-s -w" \
    -o "${build_dir}/${binary_name}" \
    "$PACKAGE"

  cp README.md LICENSE "$build_dir/"

  if [[ "$GOOS" == "windows" ]]; then
    archive_name="${artifact_base}.zip"
    (cd "$DIST_DIR" && zip -qr "$archive_name" "$artifact_base")
  else
    archive_name="${artifact_base}.tar.gz"
    tar -C "$DIST_DIR" -czf "${DIST_DIR}/${archive_name}" "$artifact_base"
  fi

  rm -rf "$build_dir"
  printf 'Created %s\n' "${DIST_DIR}/${archive_name}"
done

(cd "$DIST_DIR" && shasum -a 256 * > SHA256SUMS)

printf 'Release artifacts are ready in %s\n' "$DIST_DIR"
