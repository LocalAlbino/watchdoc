// Package internal contains watchdoc lib structs and functions
package internal

type Config struct {
	Author    string `json:"author"`
	Copyright string `json:"copyright"`
	CreatedAt bool   `json:"created_at"`
}
