// Copyright (c) 2026 Rudy Hartwig.
// Licensed under the MIT License.

package cmd

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/localalbino/watchdoc/internal/lib"
)

var testConfig = lib.Config{
	Author:    "Test Author",
	Copyright: "Copyright (c) 2026 Test.",
	CreatedAt: false,
	FileName:  false,
	Extensions: map[string]lib.Extension{
		"go":  {CommentSyntax: "//"},
		"py":  {CommentSyntax: "#"},
		"js":  {CommentSyntax: "//"},
		"sql": {CommentSyntax: "--"},
		"sh":  {CommentSyntax: "#"},
	},
}

func TestWriteHeaderFileTypes(t *testing.T) {
	tests := []struct {
		filename string
		want     string
	}{
		{
			"main.go",
			"// Copyright (c) 2026 Test.\n// \n// Author: Test Author\n\n",
		},
		{
			"script.py",
			"# Copyright (c) 2026 Test.\n# \n# Author: Test Author\n\n",
		},
		{
			"app.js",
			"// Copyright (c) 2026 Test.\n// \n// Author: Test Author\n\n",
		},
		{
			"query.sql",
			"-- Copyright (c) 2026 Test.\n-- \n-- Author: Test Author\n\n",
		},
		{
			"run.sh",
			"# Copyright (c) 2026 Test.\n# \n# Author: Test Author\n\n",
		},
	}

	dir := t.TempDir()

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			path := filepath.Join(dir, tt.filename)
			if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}

			lib.WriteHeader(&testConfig, path)

			contents, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read file after WriteHeader: %v", err)
			}

			if !strings.HasPrefix(string(contents), tt.want) {
				t.Errorf("unexpected header\ngot:  %q\nwant: %q", string(contents), tt.want)
			}
		})
	}
}

func TestWriteHeaderSkipsShebang(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "run.sh")
	original := "#!/usr/bin/env bash\necho hello\n"

	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	lib.WriteHeader(&testConfig, path)

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(contents) != original {
		t.Errorf("shebang file should be unchanged\ngot:  %q\nwant: %q", string(contents), original)
	}
}

func TestWriteHeaderUnknownExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "notes.txt")

	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	lib.WriteHeader(&testConfig, path)

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(contents) != "hello" {
		t.Errorf("file should be unchanged for unknown extension, got: %q", string(contents))
	}
}

func TestAddSubdirsExcludesDirs(t *testing.T) {
	// Build a temp tree:
	//   root/
	//     included/
	//       nested/
	//     node_modules/
	//     .git/
	root := t.TempDir()
	included := filepath.Join(root, "included")
	nested := filepath.Join(root, "included", "nested")
	nodeModules := filepath.Join(root, "node_modules")
	dotGit := filepath.Join(root, ".git")

	for _, dir := range []string{included, nested, nodeModules, dotGit} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("failed to create dir %v: %v", dir, err)
		}
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer w.Close()

	excludeDirs := []string{"node_modules", ".git"}
	addSubdirs(w, root, excludeDirs)

	watched := w.WatchList()

	// root and included/ and included/nested/ should be watched
	for _, want := range []string{root, included, nested} {
		if !slices.Contains(watched, want) {
			t.Errorf("expected %q to be watched, watched list: %v", want, watched)
		}
	}

	// excluded dirs should not be watched
	for _, unwanted := range []string{nodeModules, dotGit} {
		if slices.Contains(watched, unwanted) {
			t.Errorf("expected %q to be excluded, but it was watched", unwanted)
		}
	}
}

func TestWriteHeaderIdempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "main.go")

	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	lib.WriteHeader(&testConfig, path)

	first, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	lib.WriteHeader(&testConfig, path)

	second, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(first) != string(second) {
		t.Errorf("WriteHeader is not idempotent\nfirst:  %q\nsecond: %q", string(first), string(second))
	}
}
