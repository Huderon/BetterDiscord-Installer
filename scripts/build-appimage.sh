#!/bin/sh
# Packages the built Linux binary as an AppImage. Invoked by the GoReleaser
# post-build hook on the linux build with the binary path as the argument.
# The AppImage is uploaded via release.extra_files; it is not self-contained —
# like the zip binary it needs webkit2gtk-4.1 from the host system.
set -eu

cd "$(dirname "$0")/.."
BIN="${1:?binary path required}"

OUT_DIR=dist/appimage
APPDIR="$OUT_DIR/BetterDiscord-Installer.AppDir"

rm -rf "$APPDIR"
mkdir -p "$APPDIR/usr/bin"
cp "$BIN" "$APPDIR/usr/bin/BetterDiscord-Installer"
cp build/appicon.png "$APPDIR/betterdiscord-installer.png"
ln -sf usr/bin/BetterDiscord-Installer "$APPDIR/AppRun"

# The AppImage format requires an embedded .desktop file for menu integration
# while running; nothing is installed on the user's system.
cat > "$APPDIR/betterdiscord-installer.desktop" <<EOF
[Desktop Entry]
Type=Application
Name=BetterDiscord Installer
Exec=BetterDiscord-Installer
Icon=betterdiscord-installer
Categories=Utility;
Terminal=false
EOF

TOOL="$OUT_DIR/appimagetool"
if [ ! -x "$TOOL" ]; then
    curl -fsSL -o "$TOOL" "https://github.com/AppImage/appimagetool/releases/download/continuous/appimagetool-x86_64.AppImage"
    chmod +x "$TOOL"
fi

# --appimage-extract-and-run avoids needing FUSE on CI runners.
ARCH=x86_64 "$TOOL" --appimage-extract-and-run "$APPDIR" \
    "$OUT_DIR/BetterDiscord-Installer_linux_amd64.AppImage"
