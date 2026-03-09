package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UtilsTestSuite defines the test suite for utility functions
type UtilsTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest sets up the test environment
func (suite *UtilsTestSuite) SetupTest() {
	tempDir, err := os.MkdirTemp("", "genlinx_test")
	suite.Require().NoError(err)
	suite.tempDir = tempDir
}

// TearDownTest cleans up the test environment
func (suite *UtilsTestSuite) TearDownTest() {
	_ = os.RemoveAll(suite.tempDir)
}

// TestNormalizePath tests path normalization
func (suite *UtilsTestSuite) TestNormalizePath() {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "forward slashes on Windows",
			input:    "C:/Program Files/AMX",
			expected: filepath.FromSlash("C:/Program Files/AMX"),
		},
		{
			name:     "relative path",
			input:    "./relative/path",
			expected: filepath.FromSlash("relative/path"),
		},
		{
			name:     "path with dots",
			input:    "./foo/../bar",
			expected: filepath.FromSlash("bar"),
		},
		{
			// Regression: filepath.Clean("") returns "." which would then be
			// treated as a non-empty path and override config defaults.
			name:     "empty string stays empty",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := NormalizePath(tt.input)
			assert.Equal(suite.T(), tt.expected, result)
		})
	}
}

// TestNormalizePaths tests multiple path normalization
func (suite *UtilsTestSuite) TestNormalizePaths() {
	input := []string{
		"C:/Program Files/AMX",
		"./relative/path",
		"./foo/../bar",
	}

	expected := []string{
		filepath.FromSlash("C:/Program Files/AMX"),
		filepath.FromSlash("relative/path"),
		filepath.FromSlash("bar"),
	}

	result := NormalizePaths(input)
	assert.Equal(suite.T(), expected, result)
}

// TestFileExists tests file existence checking
func (suite *UtilsTestSuite) TestFileExists() {
	// Test existing file
	existingFile := filepath.Join(suite.tempDir, "test.txt")
	err := os.WriteFile(existingFile, []byte("test"), 0o644)
	suite.Require().NoError(err)

	assert.True(suite.T(), FileExists(existingFile))

	// Test non-existing file
	nonExistingFile := filepath.Join(suite.tempDir, "nonexistent.txt")
	assert.False(suite.T(), FileExists(nonExistingFile))
}

// TestFindFilesByExtension tests finding files by extension
func (suite *UtilsTestSuite) TestFindFilesByExtension() {
	// Create test files
	err := os.WriteFile(filepath.Join(suite.tempDir, "test1.axs"), []byte("test"), 0o644)
	suite.Require().NoError(err)
	err = os.WriteFile(filepath.Join(suite.tempDir, "test2.axi"), []byte("test"), 0o644)
	suite.Require().NoError(err)
	err = os.WriteFile(filepath.Join(suite.tempDir, "test3.txt"), []byte("test"), 0o644)
	suite.Require().NoError(err)

	// Create subdirectory with a file that must NOT be returned (non-recursive).
	subDir := filepath.Join(suite.tempDir, "subdir")
	err = os.MkdirAll(subDir, 0o755)
	suite.Require().NoError(err)
	err = os.WriteFile(filepath.Join(subDir, "test4.axs"), []byte("test"), 0o644)
	suite.Require().NoError(err)

	// Test finding .axs files — only the top-level file should be returned.
	files, err := FindFilesByExtension(suite.tempDir, ".axs")
	suite.Require().NoError(err)
	assert.Len(suite.T(), files, 1)

	expectedFiles := []string{
		filepath.Join(suite.tempDir, "test1.axs"),
	}

	assert.ElementsMatch(suite.T(), expectedFiles, files)

	// Test finding .axi files
	files, err = FindFilesByExtension(suite.tempDir, ".axi")
	suite.Require().NoError(err)
	assert.Len(suite.T(), files, 1)
	assert.Equal(suite.T(), filepath.Join(suite.tempDir, "test2.axi"), files[0])
}

// TestUtilsTestSuite runs the test suite
func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}
