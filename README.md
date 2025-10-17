# macOS Theme CLI — Minimal

A tiny macOS-only CLI that prints the current system appearance: `dark`, `light`, or `unknown`.

Build

```sh
go build .
```

Run

```sh
# prints one of: dark, light, unknown
./darght current
```

Implementation notes

- macOS only (build tag: `//go:build darwin`).
- Detection probes:
  - Primary: `defaults read -g AppleInterfaceStyle` — returns `Dark` when Dark Mode is active; empty stdout typically means Light.
  - Fallback: `osascript -e 'tell application "System Events" to tell appearance preferences to get dark mode'` — returns `true`/`false`.
- External probes use a short timeout to avoid hangs.

Programmatic usage

The repository includes a small detector implementation in `theme_detector_darwin.go`. You can call the package-level `detectTheme(ctx)` helper from Go code to get a `Theme` value (`ThemeDark`, `ThemeLight`, or `ThemeUnknown`).

License

- Add a license of your choice (no license included by default).
