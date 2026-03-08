package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// FindFilesByExtension finds all files with the specified extension in dir.
// Only the immediate directory is searched — subdirectories are not descended
// into, matching the behaviour of the TypeScript implementation.
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

// NormalizeConfigPaths normalizes all paths in a configuration
func NormalizeConfigPaths(paths []string) []string {
	return NormalizePaths(paths)
}

// GetModuleName gets the module name from package.json or returns default
func GetModuleName() string {
	// Try to read from package.json
	if FileExists("package.json") {
		// For now, return default. In a real implementation, you'd parse the JSON
		return "genlinx"
	}
	return "genlinx"
}

// GetAppVersion gets the application version
func GetAppVersion() string {
	// Try to read from package.json or use default
	if FileExists("package.json") {
		// For now, return default. In a real implementation, you'd parse the JSON
		return "2.7.0"
	}
	return "2.7.0"
}

// EnsureDir ensures a directory exists, creating it if necessary
func EnsureDir(dir string) error {
	if !DirExists(dir) {
		return os.MkdirAll(dir, 0o755)
	}
	return nil
}

// IsWindows checks if running on Windows
func IsWindows() bool {
	return os.Getenv("OS") == "Windows_NT" || strings.Contains(strings.ToLower(os.Getenv("OS")), "windows")
}

// PromptUser prompts the user for input (simplified implementation)
func PromptUser(message string) (string, error) {
	fmt.Print(message)
	var input string
	_, err := fmt.Scanln(&input)
	return strings.TrimSpace(input), err
}

// SelectFiles allows user to select from a list of files
func SelectFiles(files []string, message string) ([]string, error) {
	if len(files) == 0 {
		return []string{}, nil
	}

	if len(files) == 1 {
		return files, nil
	}

	fmt.Println(message)
	for i, file := range files {
		fmt.Printf("%d. %s\n", i+1, file)
	}

	fmt.Print("Enter file numbers (comma-separated) or 'all': ")
	input, err := PromptUser("")
	if err != nil {
		return nil, err
	}

	if strings.ToLower(input) == "all" {
		return files, nil
	}

	var selected []string
	parts := strings.Split(input, ",")
	for _, part := range parts {
		var index int
		if _, err := fmt.Sscanf(strings.TrimSpace(part), "%d", &index); err == nil {
			if index > 0 && index <= len(files) {
				selected = append(selected, files[index-1])
			}
		}
	}

	return selected, nil
}
