# darght — Minimal Dark/Light Theme CLI (macOS + Windows)

A tiny, focused command-line utility that prints the current system appearance as a single, machine-friendly token:
`dark`, `light`, or `unknown`. It's intentionally minimal and dependency-free — ideal for dotfiles,
shell scripts, and automation that needs to adapt to a user's theme.

## Highlights

- Small, single-binary tool written in Go.
- Cross-platform support: macOS (darwin) and Windows.
- Outputs exactly one token suitable for scripts: `dark`, `light`, or `unknown`.
- Build and release automation via GoReleaser and an optional Taskfile for convenience.
- Public domain license (Unlicense).

## Quick examples

Print the current theme:

```darght/README.md#L1-4
darght current
# -> dark
```

Show build/version info:

```darght/README.md#L5-8
darght version
# -> darght v0.1.0 (commit abc123, darwin/arm64)
```

## What `current` prints

- `dark` — when the system/theme is in dark mode.
- `light` — when the system/theme is in light mode.
- `unknown` — when the tool cannot determine the appearance reliably.

## Supported platforms & how detection works (brief)

- macOS:
  - Primary: `defaults read -g AppleInterfaceStyle`.
  - Fallback: `osascript` (AppleScript) probe when the primary method is inconclusive.
- Windows:
  - Primary: Read the user personalization registry value:
    `HKCU\Software\Microsoft\Windows\CurrentVersion\Themes\Personalize\AppsUseLightTheme`
    (value `0` typically means dark, `1` means light).
  - Typical probe is invoked via a short PowerShell command from the Go binary when running on Windows.
- In all cases the output is normalized to `dark`, `light`, or `unknown` for reliable scripting.

## Installation

### Prebuilt binaries

- Check the GitHub Releases page for prebuilt artifacts for macOS and Windows.
- Download the appropriate archive for your platform, extract, and move `darght` (or `darght.exe`)
  into a directory on your `PATH`.

Example (macOS):

```darght/README.md#L9-13
tar -xzf darght_X.Y.Z_darwin_amd64.tar.gz
sudo mv darght /usr/local/bin/
```

Example (Windows PowerShell):

```darght/README.md#L14-18
Expand-Archive darght_X.Y.Z_windows_amd64.zip -DestinationPath .\darght
Move-Item .\darght\darght.exe 'C:\Program Files\darght\darght.exe'
```

### Build from source

- On macOS (native build):

```darght/README.md#L19-23
go build -o darght .
# move to PATH
sudo mv darght /usr/local/bin/
```

- On Windows (native build in PowerShell or cmd):

```darght/README.md#L24-28
go build -o darght.exe .
# Move darght.exe to a directory on your PATH, e.g.:
Move-Item .\darght.exe C:\Users\<you>\bin\
```

- Cross-compile (from any host with Go toolchain) to build for Windows or another platform:

```darght/README.md#L29-33
# build for Windows (amd64)
GOOS=windows GOARCH=amd64 go build -o darght.exe .
# build for macOS (arm64)
GOOS=darwin GOARCH=arm64 go build -o darght .
```

### Using GoReleaser (local snapshot)

```darght/README.md#L34-37
goreleaser release --snapshot --clean
# or via a Taskfile helper:
# task dist
```

## Usage

- Print the current appearance (preferred for scripts):

```darght/README.md#L38-41
darght current
# prints exactly one word: dark, light, or unknown
```

- Display build metadata (version, commit, platform):

```darght/README.md#L42-44
darght version
```

## Examples

### Bash script example (macOS, Linux shells, WSL)

```darght/README.md#L45-51
if [ "$(darght current)" = "dark" ]; then
  echo "Enable dark-theme settings"
else
  echo "Enable light-theme settings"
fi
```

### PowerShell example (Windows)

```darght/README.md#L52-58
if ((darght current) -eq 'dark') {
  Write-Host "Enable dark-theme settings"
} else {
  Write-Host "Enable light-theme settings"
}
```

## Developer notes

- The code is written in Go. Platform-specific probes are gated by runtime detection / build tags:
  - macOS code may be guarded with `//go:build darwin`.
  - Windows code may be guarded with `//go:build windows`.
- The CLI provides two primary subcommands:
  - `current` — prints the theme token.
  - `version` — prints build metadata (version, commit, GOOS/GOARCH).
- Keep the tool small and deterministic — prefer stability and predictable output for scripts.

## Troubleshooting

- If `darght` prints `unknown`:
  - Ensure the binary is running on a supported platform (macOS or Windows).
  - On macOS: verify `defaults read -g AppleInterfaceStyle` behaves in your shell.
  - On Windows: ensure the registry value exists and PowerShell is available (it's included with modern Windows).
- If building on a non-native host, ensure your environment supports cross-compilation (CGO-free builds are recommended).

## Release workflow

- Tag and push a release to trigger CI (example):

```darght/README.md#L59-62
git tag vX.Y.Z
git push origin vX.Y.Z
```

- The project's CI (GoReleaser + Actions/etc.) can produce multi-arch artifacts for macOS and Windows and publish them to GitHub Releases.

## Contributing

- Contributions are welcome. Keep changes small and focused.
- Prefer bug fixes, packaging/CI improvements, and tests where applicable.
- Workflow:
  - Fork → branch → PR with a clear description and, if applicable, test coverage.

## Security & privacy

- `darght` does not send data anywhere. It only reads local system settings using standard OS facilities.

## License

- Released into the public domain under the Unlicense. See `LICENSE` for details.

## Acknowledgements

- Built to be small, dependable, and script-friendly. If you use `darght` in your dotfiles or projects, consider opening an issue to share feedback or request packaging (e.g., Homebrew tap or a Windows package).
