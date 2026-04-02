package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/localalbino/watchdoc/internal"
	"github.com/sgtdi/fswatcher"
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

		var config *internal.Config
		if err = json.Unmarshal(file, config); err != nil {
			log.Fatalf("Unable to read configuration file: %v", err)
		}

		w, err := fswatcher.New(
			fswatcher.WithCooldown(time.Millisecond * 300),
		)
		if err != nil {
			log.Fatalln("Unable to create file watcher")
		}

		ctx := context.Background()
		go w.Watch(ctx)
		fmt.Println("Watching for new files...")
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
