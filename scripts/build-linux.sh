#!/usr/bin/env sh
set -eu

TARGET_ARCH="${1:-amd64}"
OUTPUT="dist/grsai-api-linux-${TARGET_ARCH}"
GO_BUILD_CACHE="${GOCACHE:-$(pwd)/.tmp/go-build-cache}"

mkdir -p dist
mkdir -p "${GO_BUILD_CACHE}"

GOCACHE="${GO_BUILD_CACHE}" CGO_ENABLED=0 GOOS=linux GOARCH="${TARGET_ARCH}" \
  go build -trimpath -ldflags="-s -w" -o "${OUTPUT}" .

echo "Built ${OUTPUT}"
