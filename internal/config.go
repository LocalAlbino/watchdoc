// Package internal contains watchdoc lib structs and functions
package internal

type Config struct {
	// File header specifics
	Author    string `json:"author"`
	Copyright string `json:"copyright"`
	CreatedAt bool   `json:"created_at"`
	FileName  bool   `json:"file_name"`

	// CLI behavior specifics
	ExcludeDirs []string             `json:"exclude_dirs"`
	Extensions  map[string]Extension `json:"extensions"`
}

type Extension struct {
	CommentSyntax string `json:"comment_syntax"`
}
