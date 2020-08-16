// Package fs provides functions for filesystem operations.
package fs

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	// MarkdownOnly is a filter that only lets pass Markdown files.
	MarkdownOnly = func(file string) bool {
		return filepath.Ext(file) == ".md"
	}

	// NoUnderscores is a predefined filter that doesn't let pass
	// files starting with an underscore.
	NoUnderscores = func(file string) bool {
		filename := filepath.Base(file)
		return !strings.HasPrefix(filename, "_")
	}

	// ErrStreaming is returned from StreamFiles.
	ErrStreaming error = nil
)

// StreamFiles sends files in a given path that match the given
// filters through the files channel.
func StreamFiles(path string, files chan<- string, filters ...func(file string) bool) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		close(files)
		return nil
	}

	ErrStreaming = filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		for _, filter := range filters {
			if !filter(file) {
				return nil
			}
		}

		files <- file

		return nil
	})

	close(files)
	return ErrStreaming
}
