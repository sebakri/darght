package main

import (
	"context"
	"fmt"
	"os"
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
	case "-h", "--help", "help":
		fallthrough
	default:
		// Unknown subcommand — print usage then fall back to printing current theme.
		usage()
		return
	}
}

func usage() {
	fmt.Println("Usage:\n\ndarght current\tPrint current theme (dark/light/unknown)")
}
