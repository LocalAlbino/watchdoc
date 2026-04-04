# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Watchdoc is a Go CLI tool that automates adding file headers (author, copyright, etc.) across codebases. It watches for new files and inserts headers immediately based on a `watchdoc.json` config file. Designed as a convention-over-config tool that works independently of any editor or IDE.

## Build Commands

```bash
make          # fmt + build (default)
make build    # compiles to ./bin/watchdoc.exe
make fmt      # runs go fmt
make test     # run all tests
```

## Architecture

The project uses [Cobra](https://github.com/spf13/cobra) for CLI structure.

```
main.go                  → calls cmd.Execute()
cmd/root.go              → defines root "watchdoc" command
cmd/init.go              → "init" subcommand: creates watchdoc.json with defaults
cmd/watch.go             → "watch" subcommand: watches for new files and writes headers
cmd/scan.go              → "scan" subcommand: walks the project, reports files missing headers, optional --fix
internal/lib/config.go   → Config and Extension structs (JSON-serializable)
internal/lib/header.go   → header building and writing logic
```

### Config structure (`internal/lib/config.go`)

```go
type Config struct {
    Author        string
    Copyright     string
    CopyrightOnly bool
    CreatedAt     bool
    FileName      bool
    ExcludeDirs   []string
    Extensions    map[string]Extension  // keyed by file extension (e.g. "go", "py")
}

type Extension struct {
    CommentSyntax string  // e.g. "//" or "#"
}
```

The `init` command writes a `watchdoc.json` into `--root` (default: `.`) with sensible defaults for common languages and excluded dirs (`.git`, `bin`, `node_modules`, etc.).

## Key behaviors to be aware of

- `WriteHeader` skips files starting with `#!` (shebangs must stay at line 1)
- `WriteHeader` is idempotent — the existence check uses only copyright and fileName (stable fields); author is written last and excluded from the check so headers from other contributors are not overwritten, and `Created: unknown` from scan --fix is also tolerated
- `watch` only triggers on `fsnotify.Create` events, not writes, so saving an existing file never re-adds a header
- New directories created while watching are automatically added to the watcher, unless they match an `exclude_dirs` entry
- `scan` exits non-zero if any files are missing headers or have temp headers, making it suitable for CI
- `scan --fix` writes a best-effort header with `Created: unknown` and embeds `ScanAutoGenMarker` so subsequent scans can identify it as a temp header needing human review
