package utils

import "strings"

// FormatVersion normalizes a version string with a single leading 'v'.
func FormatVersion(version string) string {
	trimmed := strings.TrimSpace(version)
	trimmed = strings.TrimPrefix(trimmed, "v")
	if trimmed == "" {
		return "v0.0.0"
	}
	return "v" + trimmed
}