package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/localalbino/watchdoc/internal/lib"
)

var scanConfig = lib.Config{
	Author:    "Test Author",
	Copyright: "Copyright (c) 2026 Test.",
	CreatedAt: true,
	FileName:  false,
	ExcludeDirs: []string{
		"vendor",
		"node_modules",
	},
	Extensions: map[string]lib.Extension{
		"go": {CommentSyntax: "//"},
		"py": {CommentSyntax: "#"},
		"js": {CommentSyntax: "//"},
	},
}

// buildTree creates a directory structure under root from a map of
// relative path → file contents.
func buildTree(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("failed to create dir for %v: %v", rel, err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("failed to write %v: %v", rel, err)
		}
	}
}

func TestScanDetectsMissingHeaders(t *testing.T) {
	dir := t.TempDir()
	buildTree(t, dir, map[string]string{
		"main.go":    "package main\n",
		"handler.go": "package main\n",
		"readme.txt": "not tracked\n",
	})

	missing, temp := scanDir(&scanConfig, dir, false)

	if missing != 2 {
		t.Errorf("expected 2 missing, got %d", missing)
	}
	if temp != 0 {
		t.Errorf("expected 0 temp, got %d", temp)
	}
}

func TestScanIgnoresHeaderedFiles(t *testing.T) {
	dir := t.TempDir()

	// Write a file then run WriteHeader so it has a valid header
	path := filepath.Join(dir, "main.go")
	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	lib.WriteHeader(&scanConfig, path)

	missing, temp := scanDir(&scanConfig, dir, false)

	if missing != 0 {
		t.Errorf("expected 0 missing, got %d", missing)
	}
	if temp != 0 {
		t.Errorf("expected 0 temp, got %d", temp)
	}
}

func TestScanSkipsExcludedDirs(t *testing.T) {
	dir := t.TempDir()
	buildTree(t, dir, map[string]string{
		"main.go":              "package main\n",     // missing — should be counted
		"vendor/dep.go":        "package dep\n",      // excluded — should not be counted
		"node_modules/lib.js":  "export default {}\n", // excluded — should not be counted
	})

	missing, _ := scanDir(&scanConfig, dir, false)

	if missing != 1 {
		t.Errorf("expected 1 missing (excluded dirs skipped), got %d", missing)
	}
}

func TestScanFixWritesTempHeader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "main.go")
	if err := os.WriteFile(path, []byte("package main\n"), 0o644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	scanDir(&scanConfig, dir, true)

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if !strings.Contains(string(contents), lib.ScanAutoGenMarker) {
		t.Errorf("expected auto-gen marker in fixed file, got:\n%s", string(contents))
	}
	if !strings.Contains(string(contents), "Created: unknown") {
		t.Errorf("expected 'Created: unknown' in fixed file, got:\n%s", string(contents))
	}
}

func TestScanDetectsTempHeaders(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "main.go")
	if err := os.WriteFile(path, []byte("package main\n"), 0o644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	// First pass with fix — produces a temp header
	scanDir(&scanConfig, dir, true)

	// Second pass — should detect it as temp, not missing
	missing, temp := scanDir(&scanConfig, dir, false)

	if missing != 0 {
		t.Errorf("expected 0 missing, got %d", missing)
	}
	if temp != 1 {
		t.Errorf("expected 1 temp, got %d", temp)
	}
}
