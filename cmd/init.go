package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/localalbino/watchdoc/internal"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a watchdoc configuration for this project",
	Long: `Initializes a new watchdoc configuration for the current directory.
This config will be used for both the 'scan' and 'watch' commands.`,
	Run: run,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().String("root", ".", "Specifies the root directory for the config")
}

func run(cmd *cobra.Command, args []string) {
	config := internal.Config{
		Author:    "_",
		Copyright: "_",
		CreatedAt: false,
	}

	path, err := cmd.Flags().GetString("root")
	if err != nil {
		fmt.Println("Unable to create configuration file")
		return
	}

	if path[len(path)-1] == '/' {
		path += "watchdoc.json"
	} else {
		path += "/watchdoc.json"
	}

	if _, err := os.Stat(path); err == nil {
		fmt.Println("Configuration file already exists at " + path)
		return
	}

	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Unable to create configuration file at " + path)
		return
	}

	configBytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Unable to create configuration file at " + path)
		return
	}

	_, err = file.Write(configBytes)
	if err != nil {
		fmt.Println("Failed to write to configuration file at " + path)
		return
	}

	fmt.Println("Successfully created configuration at " + path)
}
