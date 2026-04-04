// Copyright (c) 2026 Rudy Hartwig.
// Licensed under the MIT License.

// Package cmd contains all cli commands for watchdoc
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "watchdoc",
	Short: "Automates creation of file headers",
	Long: `Watchdoc autoamtes creation of file headers across codebases based on a
configuration file. It can scan existing files for missing headers and
fill them in. It can watch for new files and create headers immediately.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
