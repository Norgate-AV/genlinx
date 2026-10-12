package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// FindFilesByExtension finds all files with the specified extension in dir.
// Only the immediate directory is searched — subdirectories are not descended
// into.
func FindFilesByExtension(dir, ext string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ext) {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}

	return files, nil
}

// NormalizePath converts a path to the OS-specific format.
// This allows users to use forward slashes in config files on Windows.
// An empty string is returned unchanged; filepath.Clean("") returns "."
// which would otherwise be treated as a non-empty path.
func NormalizePath(path string) string {
	if path == "" {
		return ""
	}

	return filepath.FromSlash(filepath.Clean(path))
}

// NormalizePaths normalizes a slice of paths
func NormalizePaths(paths []string) []string {
	normalized := make([]string, len(paths))

	for i, path := range paths {
		normalized[i] = NormalizePath(path)
	}

	return normalized
}
