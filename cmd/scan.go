package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/localalbino/watchdoc/internal/lib"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans the project for files missing headers",
	Long: `Scans the project recursively for files that are missing headers,
reporting each one and exiting with a non-zero status if any are found.
Intended for use in CI pipelines. Run from the project root.`,
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

		fix, err := cmd.Flags().GetBool("fix")
		if err != nil {
			log.Fatalln("Unable to read flags")
		}

		missing, temp := scanDir(&config, ".", fix)

		fmt.Println()

		issues := missing + temp
		switch {
		case missing > 0 && temp > 0:
			fmt.Printf("%d file(s) missing headers, %d file(s) with temp headers\n", missing, temp)
		case missing > 0:
			fmt.Printf("%d file(s) missing headers\n", missing)
		case temp > 0:
			fmt.Printf("%d file(s) with temp headers\n", temp)
		default:
			fmt.Println("All files have headers.")
		}

		if issues > 0 {
			os.Exit(1)
		}
	},
}

func scanDir(config *lib.Config, root string, fix bool) (int, int) {
	var missing, temp int
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path != root && slices.Contains(config.ExcludeDirs, d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		ext := extensionOf(path)
		if _, ok := config.Extensions[ext]; !ok {
			return nil
		}

		contents, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("  unable to read '%v', skipping\n", path)
			return nil
		}

		if lib.HasHeader(config, path, contents) {
			if strings.Contains(string(contents), lib.ScanAutoGenMarker) {
				temp++
				fmt.Printf("  [temp]    %v\n", path)
			}
			return nil
		}

		missing++
		fmt.Printf("  [missing] %v\n", path)
		if fix {
			lib.WriteScanHeader(config, path)
		}

		return nil
	})
	return missing, temp
}

func extensionOf(path string) string {
	parts := strings.Split(path, ".")
	return parts[len(parts)-1]
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().Bool("fix", false, "Write temp headers to files missing headers")
}
