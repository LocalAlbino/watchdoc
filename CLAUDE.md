# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Watchdoc is a Go CLI tool that automates adding file headers (author, copyright, etc.) across codebases. It scans for files missing headers and inserts them based on a `watchdoc.json` config file. A `watch` mode is planned to add headers to new files automatically.

## Build Commands

```bash
make          # fmt + build (default)
make build    # compiles to ./bin/watchdoc.exe
make fmt      # runs go fmt
go test ./... # run all tests
```

## Architecture

The project uses [Cobra](https://github.com/spf13/cobra) for CLI structure.

```
main.go               → calls cmd.Execute()
cmd/root.go           → defines root "watchdoc" command
cmd/init.go           → "init" subcommand: creates watchdoc.json with defaults
internal/config.go    → Config and Extension structs (JSON-serializable)
```

### Config structure (`internal/config.go`)

```go
type Config struct {
    Author      string
    Copyright   string
    CreatedAt   bool
    FileName    bool
    ExcludeDirs []string
    Extensions  map[string]Extension  // keyed by file extension (e.g. "go", "py")
}

type Extension struct {
    CommentSyntax string  // e.g. "//" or "#"
}
```

The `init` command writes a `watchdoc.json` into `--root` (default: `.`) with sensible defaults for common languages and excluded dirs (`.git`, `bin`, `node_modules`, etc.).
