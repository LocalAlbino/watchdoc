package lib

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type headerWriter func(config *Config, path string, commentSyntax string, header *strings.Builder)

func WriteHeader(config *Config, path string) {
	splitPath := strings.Split(path, ".")
	pathExt := splitPath[len(splitPath)-1]
	extension, ok := config.Extensions[pathExt]
	if !ok {
		return
	}

	commentSyntax := extension.CommentSyntax
	if commentSyntax != "" {
		commentSyntax += " "
	}

	// Get pre-existing file contents to ensure that they
	// can be appended back after the header
	existingContents, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to create header for file '%v'\n", path)
		log.Println("Unable to read existing file contents")
		return
	}

	header := strings.Builder{}
	writeCalls := []headerWriter{
		writeCopyright,
		writeFileName,
		writeAuthor,
		writeCreatedAt,
	}

	if config.CopyrightOnly {
		writeCopyright(config, path, commentSyntax, &header)
	} else {
		for i, writeFn := range writeCalls {
			writeFn(config, path, commentSyntax, &header)
			if i == 0 {
				// We do one extra line between the copyright and everything else
				header.WriteString(commentSyntax + "\n")
			}
		}
	}

	if strings.HasPrefix(string(existingContents), "#!") {
		return
	}

	if strings.HasPrefix(string(existingContents), header.String()) {
		return
	}

	// Append back on any pre-existing file contents and write with header
	header.WriteString("\n" + string(existingContents))
	if err = os.WriteFile(path, []byte(header.String()), 0o644); err != nil {
		log.Printf("Unable to write to file '%v'. Skipping it\n", path)
	}

	fmt.Println("Successfully created header")
}

func writeCopyright(config *Config, path string, commentSyntax string, header *strings.Builder) {
	if config.Copyright == "" {
		return
	}

	for line := range strings.SplitSeq(config.Copyright, "\n") {
		header.WriteString(commentSyntax + line + "\n")
	}
}

func writeAuthor(config *Config, path string, commentSyntax string, header *strings.Builder) {
	if config.Author == "" {
		return
	}

	header.WriteString(commentSyntax + "Author: " + config.Author + "\n")
}

func writeCreatedAt(config *Config, path string, commentSyntax string, header *strings.Builder) {
	if !config.CreatedAt {
		return
	}

	today := time.Now().Format(time.RFC3339)
	today = strings.Split(today, "T")[0]
	header.WriteString(commentSyntax + "Created: " + today + "\n")
}

func writeFileName(config *Config, path string, commentSyntax string, header *strings.Builder) {
	if !config.FileName {
		return
	}

	header.WriteString(commentSyntax + "File: " + path + "\n")
}
