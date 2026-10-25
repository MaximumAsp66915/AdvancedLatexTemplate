#!/bin/bash
# SPDX-License-Identifier: MPL-2.0
# build-release.sh -- build the mktex executables for Windows, macOS and Linux.
#
#   tools/build-release.sh [VERSION]       e.g. tools/build-release.sh v1.0.0
#
# Output (git-ignored): dist/mktex-windows-amd64.exe, mktex-windows-arm64.exe,
#   mktex-macos-arm64, mktex-macos-amd64, mktex-linux-amd64, mktex-linux-arm64, SHA256SUMS
# Each file is one self-contained program: templates, engine and fonts are embedded.
# Needs only Go (https://go.dev/dl/). The build uses exactly the files git would commit
# (tracked + not ignored), so local-only files are never embedded.
set -eu
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
VERSION="${1:-dev}"
export PATH="$PATH:/usr/local/go/bin"
command -v go >/dev/null 2>&1 || { echo "Error: Go is not installed (https://go.dev/dl/)" >&2; exit 1; }

WORK="$(mktemp -d "${TMPDIR:-/tmp}/atl-release.XXXXXX")"
( cd "$ROOT" && git ls-files -z --cached --others --exclude-standard | tar --null -T - --ignore-failed-read -cf - 2>/dev/null ) | tar -xf - -C "$WORK"
mkdir -p "$ROOT/dist"

TARGETS="windows/amd64 windows/arm64 darwin/arm64 darwin/amd64 linux/amd64 linux/arm64"
for t in $TARGETS; do
  goos="${t%/*}"; goarch="${t#*/}"
  os="$goos"; [[ "$goos" == darwin ]] && os=macos
  ext=""; [[ "$goos" == windows ]] && ext=".exe"
  out="$ROOT/dist/mktex-$os-$goarch$ext"
  ( cd "$WORK" && CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
      go build -trimpath -ldflags "-s -w -X main.version=$VERSION" -o "$out" ./cmd/mktex )
  printf '  built %-28s %s\n' "$(basename "$out")" "$(du -h "$out" | cut -f1)"
done

cd "$ROOT/dist"
if command -v sha256sum >/dev/null 2>&1; then sha256sum mktex-* > SHA256SUMS; else shasum -a 256 mktex-* > SHA256SUMS; fi
echo "  wrote SHA256SUMS"
echo "Publish with:  gh release create $VERSION dist/mktex-* dist/SHA256SUMS --title \"$VERSION\" --notes-file <notes>"
