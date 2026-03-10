package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ConfigTestSuite defines the test suite for configuration functions
type ConfigTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest sets up the test environment
func (suite *ConfigTestSuite) SetupTest() {
	tempDir, err := os.MkdirTemp("", "genlinx_config_test")
	suite.Require().NoError(err)
	suite.tempDir = tempDir
}

// TearDownTest cleans up the test environment
func (suite *ConfigTestSuite) TearDownTest() {
	_ = os.RemoveAll(suite.tempDir)
}

// TestLoadDefaultConfig tests loading default configuration
func (suite *ConfigTestSuite) TestLoadDefaultConfig() {
	config, err := LoadDefaultConfig()
	suite.Require().NoError(err)
	suite.Require().NotNil(config)

	// Test default values
	assert.Equal(suite.T(), "build.cfg", config.CFG.OutputFile)
	assert.Equal(suite.T(), "build.log", config.CFG.OutputLogFile)
	assert.Contains(suite.T(), config.Build.NLRC.IncludePath, filepath.FromSlash("C:/Program Files (x86)/Common Files/AMXShare/AXIs"))
	assert.Contains(suite.T(), config.Build.NLRC.ModulePath, filepath.FromSlash("C:/Program Files (x86)/Common Files/AMXShare/Duet/bundle"))
	assert.Contains(suite.T(), config.Build.NLRC.LibraryPath, filepath.FromSlash("C:/Program Files (x86)/Common Files/AMXShare/SYCs"))
}

// TestNormalizeConfigPaths tests path normalization in configuration
func (suite *ConfigTestSuite) TestNormalizeConfigPaths() {
	config := &Config{
		CFG: CFGConfig{
			IncludePath: []string{"C:/test/../normalized"},
			ModulePath:  []string{"./relative/path"},
		},
		Build: BuildConfig{
			NLRC: NLRCConfig{
				Path:        "C:/program files/../Program Files/NLRC.exe",
				IncludePath: []string{"./include"},
				ModulePath:  []string{"./module"},
				LibraryPath: []string{"./lib"},
			},
		},
		Archive: ArchiveConfig{
			ExtraFileSearchLocations: []string{"./search"},
		},
	}

	NormalizeConfigPaths(config)

	// Test normalized paths
	assert.Equal(suite.T(), filepath.FromSlash("C:/normalized"), config.CFG.IncludePath[0])
	assert.Equal(suite.T(), filepath.FromSlash("relative/path"), config.CFG.ModulePath[0])
	assert.Equal(suite.T(), filepath.FromSlash("C:/Program Files/NLRC.exe"), config.Build.NLRC.Path)
	assert.Equal(suite.T(), filepath.FromSlash("include"), config.Build.NLRC.IncludePath[0])
	assert.Equal(suite.T(), filepath.FromSlash("module"), config.Build.NLRC.ModulePath[0])
	assert.Equal(suite.T(), filepath.FromSlash("lib"), config.Build.NLRC.LibraryPath[0])
	assert.Equal(suite.T(), filepath.FromSlash("search"), config.Archive.ExtraFileSearchLocations[0])
}

// TestConfigTestSuite runs the test suite
func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
