#!/bin/sh
# Builds the frontend for embedding into the Go binary via go:embed.
set -eu

cd "$(dirname "$0")/../frontend"
bun install --frozen-lockfile
bun run build
