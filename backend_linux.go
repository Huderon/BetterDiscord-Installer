//go:build linux

package main

import "os"

// KWin on KDE Plasma (Wayland) force-applies server-side decorations to GTK
// windows that opt out of decorations, giving our frameless window a second,
// native titlebar. Under X11/XWayland the undecorated hint is honored, so
// prefer that backend and fall back to Wayland where X isn't available.
// Users can still override by setting GDK_BACKEND themselves.
func init() {
	if os.Getenv("GDK_BACKEND") == "" {
		os.Setenv("GDK_BACKEND", "x11,wayland")
	}
}
