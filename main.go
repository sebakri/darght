package main

import (
	"context"
	"fmt"
	"os"
)

var (
	// These variables are populated at build time via -ldflags.
	// Example:
	//   -ldflags "-X main.version=v0.1.0 -X main.commit=abcd123 -X main.date=2025-10-17"
	version = "dev"
	commit  = "none"
	date    = ""
)

type Theme string

const (
	ThemeDark    Theme = "dark"
	ThemeLight   Theme = "light"
	ThemeUnknown Theme = "unknown"
)

func (t Theme) String() string {
	if t == "" {
		return string(ThemeUnknown)
	}
	return string(t)
}
func FromString(theme string) Theme {
	if theme == "" {
		return ThemeUnknown
	}

	switch theme {
	case string(ThemeDark):
		return ThemeDark
	case string(ThemeLight):
		return ThemeLight
	default:
		return ThemeUnknown
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(os.Args) == 1 {
		usage()
		return
	}
	switch os.Args[1] {
	case "current":
		theme := detectTheme(ctx)
		fmt.Println(theme.String())
		return
	case "version":
		printVersion()
		return
	case "-h", "--help", "help":
		fallthrough
	default:
		// Unknown subcommand — print usage then fall back to printing current theme.
		usage()
		return
	}
}

func printVersion() {
	fmt.Printf("darght %s\n", version)
	fmt.Printf("commit: %s\n", commit)
	if date != "" {
		fmt.Printf("built:  %s\n", date)
	}
}

func usage() {
	fmt.Println("Usage:\n\ndarght current\tPrint current theme (dark/light/unknown)\ndarght version\tPrint build version information")
}
