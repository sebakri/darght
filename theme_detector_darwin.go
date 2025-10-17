//go:build darwin
// +build darwin

package main

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

// DefaultProbeTimeout is used when no explicit timeout is provided.
const DefaultProbeTimeout = 1500 * time.Millisecond

// detectTheme detects the macOS appearance (dark/light) on Darwin.
// It returns ThemeDark, ThemeLight, or ThemeUnknown.
func detectTheme(ctx context.Context) Theme {
	if mode, ok := probeDefaults(ctx); ok {
		return FromString(mode)
	}
	if mode, ok := probeAppleScript(ctx); ok {
		return FromString(mode)
	}
	// unable to determine confidently
	return ThemeUnknown
}

// probeDefaults tries to read AppleInterfaceStyle using `defaults`.
// Returns (mode, true) if a definitive mode was determined, otherwise ("", false).
func probeDefaults(ctx context.Context) (string, bool) {
	cmd := exec.CommandContext(ctx, "defaults", "read", "-g", "AppleInterfaceStyle")
	out, err := cmd.Output()
	if err != nil {
		trim := strings.TrimSpace(string(out))
		if trim == "" {
			return "light", true
		}
		return normalizeFromString(trim), true
	}
	trim := strings.TrimSpace(string(out))
	if trim == "" {
		return "light", true
	}
	return normalizeFromString(trim), true
}

// probeAppleScript uses osascript to query the appearance preferences.
func probeAppleScript(ctx context.Context) (string, bool) {
	script := `tell application "System Events" to tell appearance preferences to get dark mode`
	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}
	trim := strings.TrimSpace(string(out))
	if trim == "" {
		return "", false
	}
	if strings.EqualFold(trim, "true") || strings.EqualFold(trim, "yes") {
		return "dark", true
	}
	if strings.EqualFold(trim, "false") || strings.EqualFold(trim, "no") {
		return "light", true
	}
	return normalizeFromString(trim), true
}

// normalizeFromString maps common outputs to "dark" or "light". If ambiguous,
// return "unknown".
func normalizeFromString(s string) string {
	sl := strings.ToLower(strings.TrimSpace(s))
	if sl == "" {
		return "unknown"
	}
	if strings.Contains(sl, "dark") {
		return "dark"
	}
	if strings.Contains(sl, "light") {
		return "light"
	}
	return "unknown"
}
