#!/usr/bin/env bash
# Builds static binaries for Linux and macOS (Intel + Apple Silicon / ARM).
set -euo pipefail
mkdir -p dist
for target in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64; do
  os=${target%/*}; arch=${target#*/}
  echo "building dist/ccmatrix-$os-$arch"
  GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -o "dist/ccmatrix-$os-$arch" .
done
