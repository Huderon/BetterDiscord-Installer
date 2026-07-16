#!/bin/sh
# Renders the macOS Info.plist with the release version so GoReleaser can
# assemble the .app bundle archive from it. Usage: render-plist.sh <version>
set -eu

cd "$(dirname "$0")/.."
mkdir -p dist/darwin
sed "s|__VERSION__|${1:?version required}|g" build/darwin/Info.release.plist > dist/darwin/Info.plist
