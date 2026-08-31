#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

name="golf-$(go env GOOS)-$(go env GOARCH)"
bin="$root/dist/$name"

rm -rf "$root/dist"
mkdir -p "$root/dist"

# Assets are compiled into the binary via go:embed.
# -s -w strips symbol/debug tables; -trimpath drops local file paths.
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "$bin" .

echo "Share this single binary: $bin"
echo "Run it from anywhere: ./$name"
