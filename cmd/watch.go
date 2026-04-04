// Copyright (c) 2026 Rudy Hartwig.
// Licensed under the MIT License.

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/fsnotify/fsnotify"
	"github.com/localalbino/watchdoc/internal/lib"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watches for new files and adds headers based on the config",
	Long: `Watches for new files and adds headers based on your project's
configuration file. This command should be run from the project root.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("watchdoc.json"); err != nil {
			log.Fatalln("Configuration file not found in the current directory.\nRun 'watchdoc init' to create one.")
		}

		file, err := os.ReadFile("watchdoc.json")
		if err != nil {
			log.Fatalln("Unable to read configuration file.")
		}

		var config lib.Config
		if err = json.Unmarshal(file, &config); err != nil {
			log.Fatalf("Unable to read configuration file: %v", err)
		}

		// Normalize exclude path names once before recursively adding directories.
		// This should prevent issues arising from mismatches in how the user types
		// names in the config and what each os.DirEntry gives for its name.
		excludeDirs := make([]string, 0, len(config.ExcludeDirs))
		for _, dir := range config.ExcludeDirs {
			stats, err := os.Stat(dir)
			if err != nil {
				continue
			} else if !stats.IsDir() {
				log.Fatalf("File '%v' is listed as an excluded directory but is not a directory\n", stats.Name())
			}

			excludeDirs = append(excludeDirs, stats.Name())
		}

		w, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatalln("Unable to create file watcher")
		}
		defer w.Close()

		go listen(&config, w, excludeDirs)

		addSubdirs(w, ".", excludeDirs)
		fmt.Println("Watching for new files...")

		<-make(chan struct{}) // Blocks return so that we keep listening forever
	},
}

func addSubdirs(w *fsnotify.Watcher, dir string, excludeDirs []string) {
	w.Add(dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalln("Unable to create file watcher")
	}

	for _, entry := range entries {
		if entry.IsDir() && !slices.Contains(excludeDirs, entry.Name()) {
			addSubdirs(w, filepath.Join(dir, entry.Name()), excludeDirs)
		}
	}
}

func listen(config *lib.Config, w *fsnotify.Watcher, excludeDirs []string) {
	go func() {
		for err := range w.Errors {
			fmt.Printf("Watcher error: %v\n", err)
		}
	}()

	for event := range w.Events {
		if event.Has(fsnotify.Create) {
			file, err := os.Stat(event.Name)
			if err != nil {
				fmt.Printf("Unable to open file '%v'\n", event.Name)
				continue
			} else if file.IsDir() {
				if !slices.Contains(excludeDirs, filepath.Base(event.Name)) {
					w.Add(event.Name)
					fmt.Printf("Added new directory '%v' to be watched\n", event.Name)
				}
				continue
			}

			fmt.Printf("New file '%v' detected\n", event.Name)
			lib.WriteHeader(config, event.Name)
		}
	}
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
