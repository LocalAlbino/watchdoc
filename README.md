# watchdoc

Watchdoc is a Go CLI tool that automatically adds file headers — author, copyright, creation date, and more — to new files as you create them. It runs in the background and reacts to your filesystem, so it works the same way regardless of which editor, IDE, or tool you use to create files.

It is designed for teams and projects that enforce a header convention. Run `watchdoc watch` once at the start of your session and never think about headers again. Run `watchdoc scan` in CI to catch any files that slipped through.

## Installation

```bash
go install github.com/localalbino/watchdoc@latest
```

Or build from source:

```bash
make build
# outputs ./bin/watchdoc.exe
```

## Usage

### 1. Initialize a config

Run this once in your project root:

```bash
watchdoc init
```

This creates a `watchdoc.json` with sensible defaults for common languages and a standard set of excluded directories. Open it and fill in your `author` and `copyright` fields.

To initialize in a different directory:

```bash
watchdoc init --root ./my-project
```

### 2. Start watching

```bash
watchdoc watch
```

Run this from your project root. Watchdoc will recursively watch all non-excluded directories for new files and immediately prepend a header when one is detected.

### 3. Scan for missing headers

```bash
watchdoc scan
```

Walks the project recursively and reports any files (whose extensions are listed in `extensions`) that are missing a header. Exits with a non-zero status if any are found, making it suitable as a CI check.

```bash
watchdoc scan --fix
```

Same as above, but writes a best-effort header to each missing file. Because the scanner does not know who created the file or when, it uses your local `author` from `watchdoc.json`, sets `Created: unknown`, and appends a note marking the header as auto-generated. These files are flagged as `[temp]` on subsequent scans to remind contributors to review and update them.

## Configuration

`watchdoc.json` lives in your project root and controls everything:

```json
{
  "author": "Your Name",
  "copyright": "Copyright (c) 2026 Your Name.\nLicensed under the MIT License.",
  "copyright_only": false,
  "created_at": true,
  "file_name": false,
  "exclude_dirs": [".git", "bin", "node_modules"],
  "extensions": {
    "go": { "comment_syntax": "//" },
    "py": { "comment_syntax": "#" },
    "sql": { "comment_syntax": "--" }
  }
}
```

| Field            | Type             | Description                                                 |
| ---------------- | ---------------- | ----------------------------------------------------------- |
| `author`         | string           | Written as `Author: <value>`                                |
| `copyright`      | string           | Supports `\n` for multi-line blocks                         |
| `copyright_only` | bool             | When true, only the copyright block is written              |
| `created_at`     | bool             | Appends `Created: <YYYY-MM-DD>` using today's date          |
| `file_name`      | bool             | Appends `File: <path>` using the path relative to the root  |
| `exclude_dirs`   | array of strings | Directories to ignore entirely (matched by name, recursive) |
| `extensions`     | map              | File extensions to handle, each with a `comment_syntax`     |

Files whose extension is not listed in `extensions` are silently skipped.

### Example header output

For a Go file with the config above:

```
// Copyright (c) 2026 Your Name.
// Licensed under the MIT License.
//
// Created: 2026-04-04
// Author: Your Name
```

## Things to note

**Editor-agnostic by design.** Watchdoc listens at the OS filesystem level via [fsnotify](https://github.com/fsnotify/fsnotify). It does not integrate with any editor, which means it works consistently across Vim, Neovim, VS Code, JetBrains, and anything else — including scripts and generators that create files programmatically.

**Handles atomic saves.** Editors like Vim and Neovim save files by writing to a temporary file and renaming it over the original. This rename triggers a `Create` event, which could cause headers to be re-added on every save. Watchdoc detects this by checking whether the file already starts with the expected header and skips it if so.

**Shebangs are preserved.** Files starting with `#!` are never modified, since the shebang must remain on the first line to be recognized by the OS.

**Headers are only added on file creation.** Watchdoc listens for `Create` events only, not `Write` events. Editing and saving an existing file will never trigger a header write.

**New directories are picked up automatically.** If you create a new directory while `watchdoc watch` is running, it is added to the watch list immediately (unless it matches an `exclude_dirs` entry).

**Headers are cross-contributor safe.** The existence check only looks at the copyright and fileName fields — the author is written last and ignored during the check. This means a file headered by one contributor will not be re-headered when another contributor runs `watch` with a different `author` in their local config.

**`scan --fix` headers are marked for review.** Files fixed by `watchdoc scan --fix` receive a header with `Created: unknown` and an embedded auto-gen note. These show up as `[temp]` on subsequent scans so they are not silently treated as fully resolved.
