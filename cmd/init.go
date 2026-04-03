package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/localalbino/watchdoc/internal"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a watchdoc configuration for this project",
	Long: `Initializes a new watchdoc configuration for the current directory.
This config will be used for both the 'scan' and 'watch' commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := internal.Config{
			Author:    "_",
			Copyright: "_",
			CreatedAt: false,
			FileName:  false,
			ExcludeDirs: []string{
				// General
				".git", "bin", "dist", "out", "build", "vendor", "tmp", "log", "logs", "coverage",
				// JavaScript / Node
				"node_modules", ".cache", ".turbo",
				// JS framework build/cache dirs
				".next", ".nuxt", ".svelte-kit", ".vite",
				// Python
				"__pycache__", ".venv", "venv", ".pytest_cache", ".mypy_cache",
				// Java / Kotlin (Maven + Gradle)
				"target", ".gradle",
				// .NET / C#
				"obj", ".vs",
				// IDEs
				".idea", ".vscode", ".eclipse",
			},
			Extensions: map[string]internal.Extension{
				// JavaScript / TypeScript
				"js":  {CommentSyntax: "//"},
				"ts":  {CommentSyntax: "//"},
				"jsx": {CommentSyntax: "//"},
				"tsx": {CommentSyntax: "//"},
				// Frontend frameworks
				"vue":    {CommentSyntax: "//"},
				"svelte": {CommentSyntax: "//"},
				// Styles
				"scss": {CommentSyntax: "//"},
				"sass": {CommentSyntax: "//"},
				// Backend languages
				"go":   {CommentSyntax: "//"},
				"py":   {CommentSyntax: "#"},
				"rb":   {CommentSyntax: "#"},
				"php":  {CommentSyntax: "//"},
				"java": {CommentSyntax: "//"},
				"cs":   {CommentSyntax: "//"},
				"rs":   {CommentSyntax: "//"},
				"kt":   {CommentSyntax: "//"},
				// Shell / scripting
				"sh":  {CommentSyntax: "#"},
				"sql": {CommentSyntax: "--"},
			},
		}

		root, err := cmd.Flags().GetString("root")
		if err != nil {
			log.Fatalln("Unable to create configuration file")
		}

		if _, err := os.Stat(root); os.IsNotExist(err) {
			log.Fatalln("Directory does not exist: " + root)
		}

		path := filepath.Join(root, "watchdoc.json")

		if _, err := os.Stat(path); err == nil {
			fmt.Println("Configuration file already exists at " + path)
			return
		}

		file, err := os.Create(path)
		if err != nil {
			log.Fatalln("Unable to create configuration file at " + path)
		}

		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatalln("Unable to create configuration file at " + path)
		}

		_, err = file.Write(configBytes)
		if err != nil {
			log.Fatalln("Failed to write to configuration file at " + path)
		}

		fmt.Println("Successfully created configuration at " + path)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().String("root", ".", "Specifies the root directory for the config")
}
