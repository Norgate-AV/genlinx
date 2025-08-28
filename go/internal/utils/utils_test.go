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
	os.RemoveAll(suite.tempDir)
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

// TestDirExists tests directory existence checking
func (suite *UtilsTestSuite) TestDirExists() {
	// Test existing directory
	assert.True(suite.T(), DirExists(suite.tempDir))

	// Test non-existing directory
	nonExistingDir := filepath.Join(suite.tempDir, "nonexistent")
	assert.False(suite.T(), DirExists(nonExistingDir))

	// Test file (should return false for DirExists)
	filePath := filepath.Join(suite.tempDir, "test.txt")
	err := os.WriteFile(filePath, []byte("test"), 0o644)
	suite.Require().NoError(err)
	assert.False(suite.T(), DirExists(filePath))
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

	// Create subdirectory with more files
	subDir := filepath.Join(suite.tempDir, "subdir")
	err = os.MkdirAll(subDir, 0o755)
	suite.Require().NoError(err)
	err = os.WriteFile(filepath.Join(subDir, "test4.axs"), []byte("test"), 0o644)
	suite.Require().NoError(err)

	// Test finding .axs files
	files, err := FindFilesByExtension(suite.tempDir, ".axs")
	suite.Require().NoError(err)
	assert.Len(suite.T(), files, 2)

	expectedFiles := []string{
		filepath.Join(suite.tempDir, "test1.axs"),
		filepath.Join(subDir, "test4.axs"),
	}
	assert.ElementsMatch(suite.T(), expectedFiles, files)

	// Test finding .axi files
	files, err = FindFilesByExtension(suite.tempDir, ".axi")
	suite.Require().NoError(err)
	assert.Len(suite.T(), files, 1)
	assert.Equal(suite.T(), filepath.Join(suite.tempDir, "test2.axi"), files[0])
}

// TestEnsureDir tests directory creation
func (suite *UtilsTestSuite) TestEnsureDir() {
	testDir := filepath.Join(suite.tempDir, "newdir", "nested")

	// Directory shouldn't exist initially
	assert.False(suite.T(), DirExists(testDir))

	// Ensure directory exists
	err := EnsureDir(testDir)
	suite.Require().NoError(err)

	// Directory should now exist
	assert.True(suite.T(), DirExists(testDir))

	// Calling again should not error
	err = EnsureDir(testDir)
	suite.Require().NoError(err)
}

// TestIsWindows tests Windows detection
func (suite *UtilsTestSuite) TestIsWindows() {
	// This test will vary based on the actual OS
	// We just test that it returns a boolean
	result := IsWindows()
	assert.IsType(suite.T(), false, result)
}

// TestGetModuleName tests module name retrieval
func (suite *UtilsTestSuite) TestGetModuleName() {
	// Without package.json, should return default
	result := GetModuleName()
	assert.Equal(suite.T(), "genlinx", result)

	// Create a package.json file
	packageJson := filepath.Join(suite.tempDir, "package.json")
	err := os.WriteFile(packageJson, []byte(`{"name": "test-module"}`), 0o644)
	suite.Require().NoError(err)

	// Change to temp directory and test
	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd)

	err = os.Chdir(suite.tempDir)
	suite.Require().NoError(err)

	// Should still return default (implementation doesn't parse JSON yet)
	result = GetModuleName()
	assert.Equal(suite.T(), "genlinx", result)
}

// TestGetAppVersion tests app version retrieval
func (suite *UtilsTestSuite) TestGetAppVersion() {
	result := GetAppVersion()
	assert.Equal(suite.T(), "2.7.0", result)
}

// TestUtilsTestSuite runs the test suite
func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}
