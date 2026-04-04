package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/localalbino/watchdoc/internal/lib"
)

func TestInitCreatesConfig(t *testing.T) {
	dir := t.TempDir()

	rootCmd.SetArgs([]string{"init", "--root", dir})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	path := filepath.Join(dir, "watchdoc.json")
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("watchdoc.json was not created: %v", err)
	}

	var config lib.Config
	if err = json.Unmarshal(contents, &config); err != nil {
		t.Fatalf("watchdoc.json is not valid JSON: %v", err)
	}

	if len(config.Extensions) == 0 {
		t.Error("expected extensions to be populated")
	}
	if len(config.ExcludeDirs) == 0 {
		t.Error("expected exclude_dirs to be populated")
	}
}

func TestInitDoesNotOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "watchdoc.json")

	original := []byte(`{"author":"original"}`)
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatalf("failed to create existing config: %v", err)
	}

	rootCmd.SetArgs([]string{"init", "--root", dir})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	if string(contents) != string(original) {
		t.Errorf("init overwrote existing config\ngot:  %q\nwant: %q", string(contents), string(original))
	}
}
