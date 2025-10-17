//go:build windows
// +build windows

package main

import (
	"context"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

// DefaultProbeTimeout is used when no explicit timeout is provided.
const DefaultProbeTimeout = 1500 * time.Millisecond

// detectTheme detects the Windows appearance (dark/light) on Windows.
// It returns ThemeDark, ThemeLight, or ThemeUnknown.
func detectTheme(ctx context.Context) Theme {
	// honor context cancellation / timeout
	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), DefaultProbeTimeout)
	}
	select {
	case <-ctx.Done():
		return ThemeUnknown
	default:
	}

	if mode, ok := probeRegistry(ctx); ok {
		return FromString(mode)
	}
	return ThemeUnknown
}

// probeRegistry reads the relevant personalization registry values to determine theme.
// It prefers the AppsUseLightTheme value (affects apps), and falls back to SystemUsesLightTheme.
// Returns (mode, true) if a definitive mode was determined, otherwise ("", false).
func probeRegistry(ctx context.Context) (string, bool) {
	// If context is already done, bail out early.
	select {
	case <-ctx.Done():
		return "", false
	default:
	}

	const personalizePath = `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`
	k, err := registry.OpenKey(registry.CURRENT_USER, personalizePath, registry.QUERY_VALUE)
	if err != nil {
		// Unable to open personalize key — treat as unknown
		return "", false
	}
	defer k.Close()

	// Helper to read a DWORD value (0 or 1).
	readValue := func(name string) (string, bool) {
		// Check for context cancelation before doing the read.
		select {
		case <-ctx.Done():
			return "", false
		default:
		}

		val, _, err := k.GetIntegerValue(name)
		if err != nil {
			return "", false
		}
		switch val {
		case 0:
			return "dark", true
		case 1:
			return "light", true
		default:
			// Unexpected numeric value; treat as unknown but considered a probe result.
			return "unknown", true
		}
	}

	// Prefer AppsUseLightTheme (affects apps)
	if mode, ok := readValue("AppsUseLightTheme"); ok {
		return mode, true
	}

	// Fall back to system-level preference
	if mode, ok := readValue("SystemUsesLightTheme"); ok {
		return mode, true
	}

	return "", false
}

// normalizeFromString maps common outputs to "dark" or "light". If ambiguous,
// return "unknown". This mirrors the macOS implementation for consistency.
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
