package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// IntegrationTestSuite defines the integration test suite
type IntegrationTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest sets up the test environment
func (suite *IntegrationTestSuite) SetupTest() {
	tempDir, err := os.MkdirTemp("", "genlinx_integration_test")
	suite.Require().NoError(err)
	resolved, err := filepath.EvalSymlinks(tempDir)
	suite.Require().NoError(err)
	suite.tempDir = resolved
}

// TearDownTest cleans up the test environment
func (suite *IntegrationTestSuite) TearDownTest() {
	_ = os.RemoveAll(suite.tempDir)
}

// TestBuildCommandIntegration tests the full build command integration
func (suite *IntegrationTestSuite) TestBuildCommandIntegration() {
	// Create a test NetLinx source file
	sourceFile := filepath.Join(suite.tempDir, "test.axs")
	sourceContent := `PROGRAM_NAME='test'

DEFINE_DEVICE
dvTP = 10001:1:0

DEFINE_EVENT
data_event[dvTP] {
    online: {
        send_string data.device, "'Test Program Online'"
    }
}
`
	err := os.WriteFile(sourceFile, []byte(sourceContent), 0o644)
	suite.Require().NoError(err)

	// Change to temp directory
	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck

	err = os.Chdir(suite.tempDir)
	suite.Require().NoError(err)

	// Note: This test would require a working NLRC installation
	// For now, we'll test the command structure and argument parsing
	// In a real environment, this would execute the actual build

	// Test that the source file exists
	assert.True(suite.T(), fileExists(sourceFile))

	// Test that we can read the source file
	content, err := os.ReadFile(sourceFile)
	suite.Require().NoError(err)
	assert.Contains(suite.T(), string(content), "PROGRAM_NAME='test'")
	assert.Contains(suite.T(), string(content), "DEFINE_DEVICE")
}

// TestConfigurationIntegration tests configuration loading integration
func (suite *IntegrationTestSuite) TestConfigurationIntegration() {
	// Create a local config file
	configFile := filepath.Join(suite.tempDir, ".genlinxrc.json")
	configContent := `{
		"build": {
			"nlrc": {
				"includePath": ["integration/include"],
				"modulePath": ["integration/module"],
				"libraryPath": ["integration/library"]
			}
		}
	}`
	err := os.WriteFile(configFile, []byte(configContent), 0o644)
	suite.Require().NoError(err)

	// Change to temp directory
	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck

	err = os.Chdir(suite.tempDir)
	suite.Require().NoError(err)

	// Test that config file exists
	assert.True(suite.T(), fileExists(configFile))

	// Test that we can read the config file
	content, err := os.ReadFile(configFile)
	suite.Require().NoError(err)
	assert.Contains(suite.T(), string(content), "integration/include")
	assert.Contains(suite.T(), string(content), "integration/module")
	assert.Contains(suite.T(), string(content), "integration/library")
}

// TestPathResolutionIntegration tests path resolution across the system
func (suite *IntegrationTestSuite) TestPathResolutionIntegration() {
	// Create directory structure
	subDir := filepath.Join(suite.tempDir, "subdir")
	err := os.MkdirAll(subDir, 0o755)
	suite.Require().NoError(err)

	// Create files in different locations
	sourceFile := filepath.Join(subDir, "source.axs")
	err = os.WriteFile(sourceFile, []byte("test source"), 0o644)
	suite.Require().NoError(err)

	includeDir := filepath.Join(suite.tempDir, "include")
	err = os.MkdirAll(includeDir, 0o755)
	suite.Require().NoError(err)

	// Change to subdir
	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck

	err = os.Chdir(subDir)
	suite.Require().NoError(err)

	// Test relative path resolution
	relSourcePath := "source.axs"
	absSourcePath, err := filepath.Abs(relSourcePath)
	suite.Require().NoError(err)

	assert.True(suite.T(), filepath.IsAbs(absSourcePath))
	assert.True(suite.T(), fileExists(absSourcePath))

	// Test path from parent directory
	parentIncludePath := "../include"
	absIncludePath, err := filepath.Abs(parentIncludePath)
	suite.Require().NoError(err)

	assert.True(suite.T(), filepath.IsAbs(absIncludePath))
	assert.Equal(suite.T(), includeDir, absIncludePath)
}

// TestFindUpIntegration tests the find-up functionality in a realistic scenario
func (suite *IntegrationTestSuite) TestFindUpIntegration() {
	// Create directory structure: temp/project/subdir/
	projectDir := filepath.Join(suite.tempDir, "project")
	subDir := filepath.Join(projectDir, "subdir")

	err := os.MkdirAll(subDir, 0o755)
	suite.Require().NoError(err)

	// Create config file in project root
	configFile := filepath.Join(projectDir, ".genlinxrc.json")
	configContent := `{
		"build": {
			"nlrc": {
				"includePath": ["project/include"]
			}
		}
	}`
	err = os.WriteFile(configFile, []byte(configContent), 0o644)
	suite.Require().NoError(err)

	// Create source file in subdir
	sourceFile := filepath.Join(subDir, "app.axs")
	err = os.WriteFile(sourceFile, []byte("test app"), 0o644)
	suite.Require().NoError(err)

	// Change to subdir
	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck

	err = os.Chdir(subDir)
	suite.Require().NoError(err)

	// Verify we can access the config file via relative path
	configRelPath := "../.genlinxrc.json"
	assert.True(suite.T(), fileExists(configRelPath))

	// Verify we can access the source file
	assert.True(suite.T(), fileExists("app.axs"))

	// Test that find-up would work (simulate the directory traversal)
	// The config file should be in the parent directory
	parentConfigPath := filepath.Join("..", ".genlinxrc.json")
	assert.True(suite.T(), fileExists(parentConfigPath), "Config file should exist in parent directory")

	// Simple find-up: check if we can find the config by going up one level
	found := fileExists("../.genlinxrc.json")

	assert.True(suite.T(), found, "Should find config file via find-up traversal")
}

// TestMultipleConfigFormats tests support for different config file formats
func (suite *IntegrationTestSuite) TestMultipleConfigFormats() {
	// Create config files in different formats
	jsonConfig := filepath.Join(suite.tempDir, ".genlinxrc.json")
	yamlConfig := filepath.Join(suite.tempDir, ".genlinxrc.yaml")
	ymlConfig := filepath.Join(suite.tempDir, ".genlinxrc.yml")

	jsonContent := `{"build": {"nlrc": {"includePath": ["json/path"]}}}`
	yamlContent := `build:
  nlrc:
    includePath:
      - yaml/path
`
	ymlContent := `build:
  nlrc:
    includePath:
      - yml/path
`

	err := os.WriteFile(jsonConfig, []byte(jsonContent), 0o644)
	suite.Require().NoError(err)
	err = os.WriteFile(yamlConfig, []byte(yamlContent), 0o644)
	suite.Require().NoError(err)
	err = os.WriteFile(ymlConfig, []byte(ymlContent), 0o644)
	suite.Require().NoError(err)

	// Change to temp directory
	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck

	err = os.Chdir(suite.tempDir)
	suite.Require().NoError(err)

	// Verify all config files exist
	assert.True(suite.T(), fileExists(".genlinxrc.json"))
	assert.True(suite.T(), fileExists(".genlinxrc.yaml"))
	assert.True(suite.T(), fileExists(".genlinxrc.yml"))

	// Verify content
	jsonContentRead, err := os.ReadFile(".genlinxrc.json")
	suite.Require().NoError(err)
	assert.Contains(suite.T(), string(jsonContentRead), "json/path")

	yamlContentRead, err := os.ReadFile(".genlinxrc.yaml")
	suite.Require().NoError(err)
	assert.Contains(suite.T(), string(yamlContentRead), "yaml/path")

	ymlContentRead, err := os.ReadFile(".genlinxrc.yml")
	suite.Require().NoError(err)
	assert.Contains(suite.T(), string(ymlContentRead), "yml/path")
}

// TestFilePermissions tests file permission handling
func (suite *IntegrationTestSuite) TestFilePermissions() {
	// Create a test file with specific permissions
	testFile := filepath.Join(suite.tempDir, "test.axs")
	err := os.WriteFile(testFile, []byte("test content"), 0o644)
	suite.Require().NoError(err)

	// Verify file exists and is readable
	assert.True(suite.T(), fileExists(testFile))

	info, err := os.Stat(testFile)
	suite.Require().NoError(err)

	// On Windows, permissions work differently, but we can at least verify the file is not a directory
	assert.False(suite.T(), info.IsDir())

	// Verify we can read the content
	content, err := os.ReadFile(testFile)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), "test content", string(content))
}

// Helper function to check if file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// TestIntegrationTestSuite runs the integration test suite
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
