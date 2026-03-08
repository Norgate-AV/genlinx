package options

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/Norgate-AV/genlinx-go/internal/config"
)

// OptionsTestSuite defines the test suite for options functions
type OptionsTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest sets up the test environment
func (suite *OptionsTestSuite) SetupTest() {
	tempDir, err := os.MkdirTemp("", "genlinx_options_test")
	suite.Require().NoError(err)
	suite.tempDir = tempDir
}

// TearDownTest cleans up the test environment
func (suite *OptionsTestSuite) TearDownTest() {
	os.RemoveAll(suite.tempDir)
}

// TestLoadBuildOptions tests loading and merging build options
func (suite *OptionsTestSuite) TestLoadBuildOptions() {
	// Create CLI options
	cliOpts := &CLIOptions{
		SourceFiles: []string{"test.axs"},
		IncludePath: []string{"cli/include"},
		ModulePath:  []string{"cli/module"},
		Verbose:     true,
	}

	// Test loading options
	opts, configInfo, err := LoadBuildOptions(cliOpts)
	suite.Require().NoError(err)
	suite.Require().NotNil(opts)
	suite.Require().NotNil(configInfo)

	// Verify config loading info
	suite.True(configInfo.DefaultLoaded)
	suite.NotNil(configInfo.GlobalResult.Config)
	suite.NotNil(configInfo.LocalResult.Config)

	// Verify CLI options are included (should be prepended to config paths)
	// Source files get resolved to absolute paths
	// Note: paths are resolved relative to current working directory (options package dir)

	// Just check that the paths contain the expected relative components
	sourcePathFound := false
	for _, path := range opts.SourceFiles {
		if strings.Contains(path, "test.axs") {
			sourcePathFound = true
			break
		}
	}
	assert.True(suite.T(), sourcePathFound, "Should contain source file with test.axs")

	includePathFound := false
	for _, path := range opts.IncludePath {
		if strings.Contains(path, "cli") && strings.Contains(path, "include") {
			includePathFound = true
			break
		}
	}
	assert.True(suite.T(), includePathFound, "Should contain include path with cli/include")

	modulePathFound := false
	for _, path := range opts.ModulePath {
		if strings.Contains(path, "cli") && strings.Contains(path, "module") {
			modulePathFound = true
			break
		}
	}
	assert.True(suite.T(), modulePathFound, "Should contain module path with cli/module")
	assert.True(suite.T(), opts.Verbose)
}

// TestPrependAndDeduplicate tests array merging with prepending and deduplication
func (suite *OptionsTestSuite) TestPrependAndDeduplicate() {
	existing := []string{"existing1", "existing2", "duplicate"}
	newItems := []string{"new1", "duplicate", "new2"}

	result := prependAndDeduplicate(existing, newItems)

	// New items should be first
	assert.Equal(suite.T(), "new1", result[0])
	assert.Equal(suite.T(), "duplicate", result[1]) // Only once
	assert.Equal(suite.T(), "new2", result[2])

	// Existing items should follow (without duplicates)
	assert.Contains(suite.T(), result, "existing1")
	assert.Contains(suite.T(), result, "existing2")

	// Duplicate should only appear once
	duplicateCount := 0
	for _, item := range result {
		if item == "duplicate" {
			duplicateCount++
		}
	}
	assert.Equal(suite.T(), 1, duplicateCount)
}

// TestResolvePaths tests path resolution to absolute paths
func (suite *OptionsTestSuite) TestResolvePaths() {
	opts := &BuildOptions{
		SourceFiles: []string{"relative/file.axs"},
		CFGFiles:    []string{"./config.cfg"},
		IncludePath: []string{"./include"},
		ModulePath:  []string{"./module"},
		LibraryPath: []string{"./lib"},
		OutputPath:  "./output",
		NLRCPath:    "./nlrc.exe",
		ShellPath:   "./shell.exe",
	}

	// Change to temp directory for relative path testing
	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd)

	err = os.Chdir(suite.tempDir)
	suite.Require().NoError(err)

	err = opts.resolvePaths()
	suite.Require().NoError(err)

	// Verify all paths are absolute
	for _, path := range opts.SourceFiles {
		assert.True(suite.T(), filepath.IsAbs(path), "Source file should be absolute: %s", path)
	}
	for _, path := range opts.CFGFiles {
		assert.True(suite.T(), filepath.IsAbs(path), "CFG file should be absolute: %s", path)
	}
	for _, path := range opts.IncludePath {
		assert.True(suite.T(), filepath.IsAbs(path), "Include path should be absolute: %s", path)
	}
	for _, path := range opts.ModulePath {
		assert.True(suite.T(), filepath.IsAbs(path), "Module path should be absolute: %s", path)
	}
	for _, path := range opts.LibraryPath {
		assert.True(suite.T(), filepath.IsAbs(path), "Library path should be absolute: %s", path)
	}
	assert.True(suite.T(), filepath.IsAbs(opts.OutputPath), "Output path should be absolute: %s", opts.OutputPath)
	assert.True(suite.T(), filepath.IsAbs(opts.NLRCPath), "NLRC path should be absolute: %s", opts.NLRCPath)
	assert.True(suite.T(), filepath.IsAbs(opts.ShellPath), "Shell path should be absolute: %s", opts.ShellPath)
}

// TestLoadGlobalConfig tests global configuration loading
func (suite *OptionsTestSuite) TestLoadGlobalConfig() {
	// Redirect the global config search to an empty temp dir so the test is
	// not affected by any real config file present on the developer's machine.
	suite.T().Setenv("GENLINX_CONFIG_DIR", suite.tempDir)

	// Test loading global config
	result, err := loadGlobalConfig()
	suite.Require().NoError(err)
	// Config should be empty since no global config exists in test environment
	assert.NotNil(suite.T(), result.Config)
	assert.False(suite.T(), result.Found)
}

// TestLoadLocalConfig tests local configuration loading with find-up
func (suite *OptionsTestSuite) TestLoadLocalConfig() {
	// Create a config file in the temp directory
	configFile := filepath.Join(suite.tempDir, ".genlinxrc.json")
	configContent := `{
		"build": {
			"nlrc": {
				"includePath": ["test/include"]
			}
		}
	}`
	err := os.WriteFile(configFile, []byte(configContent), 0o644)
	suite.Require().NoError(err)

	// Change to temp directory
	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd)

	err = os.Chdir(suite.tempDir)
	suite.Require().NoError(err)

	// Test loading local config
	result, err := loadLocalConfig()
	suite.Require().NoError(err)
	assert.NotNil(suite.T(), result.Config)
	assert.True(suite.T(), result.Found)
	assert.Contains(suite.T(), result.Config.Build.NLRC.IncludePath, "test/include")
}

// TestLoadLocalConfigFindUp tests the find-up functionality
func (suite *OptionsTestSuite) TestLoadLocalConfigFindUp() {
	// Create directory structure: temp/parent/child/
	parentDir := filepath.Join(suite.tempDir, "parent")
	childDir := filepath.Join(parentDir, "child")

	err := os.MkdirAll(childDir, 0o755)
	suite.Require().NoError(err)

	// Create config file in parent directory
	configFile := filepath.Join(parentDir, ".genlinxrc.json")
	configContent := `{
		"build": {
			"nlrc": {
				"modulePath": ["findup/module"]
			}
		}
	}`
	err = os.WriteFile(configFile, []byte(configContent), 0o644)
	suite.Require().NoError(err)

	// Change to child directory
	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd)

	err = os.Chdir(childDir)
	suite.Require().NoError(err)

	// Test that config is found via find-up
	result, err := loadLocalConfig()
	suite.Require().NoError(err)
	assert.NotNil(suite.T(), result.Config)
	assert.True(suite.T(), result.Found)
	assert.Contains(suite.T(), result.Config.Build.NLRC.ModulePath, "findup/module")
}

// TestLoadLocalConfigNoConfig tests when no local config is found
func (suite *OptionsTestSuite) TestLoadLocalConfigNoConfig() {
	// Change to temp directory (no config files)
	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd)

	err = os.Chdir(suite.tempDir)
	suite.Require().NoError(err)

	// Test loading when no config exists
	result, err := loadLocalConfig()
	suite.Require().NoError(err)
	assert.NotNil(suite.T(), result.Config)
	assert.False(suite.T(), result.Found)
	// Should return empty config
	assert.Empty(suite.T(), result.Config.Build.NLRC.IncludePath)
}

// TestMergeConfigs tests configuration merging
func (suite *OptionsTestSuite) TestMergeConfigs() {
	config1 := &config.Config{
		Build: config.BuildConfig{
			NLRC: config.NLRCConfig{
				IncludePath: []string{"config1/include"},
				Path:        "config1.exe",
			},
		},
	}

	config2 := &config.Config{
		Build: config.BuildConfig{
			NLRC: config.NLRCConfig{
				IncludePath: []string{"config2/include"},
				ModulePath:  []string{"config2/module"},
			},
		},
	}

	merged := mergeConfigs(config1, config2)

	// Should contain items from both configs
	assert.Contains(suite.T(), merged.Build.NLRC.IncludePath, "config1/include")
	assert.Contains(suite.T(), merged.Build.NLRC.IncludePath, "config2/include")
	assert.Contains(suite.T(), merged.Build.NLRC.ModulePath, "config2/module")
	assert.Equal(suite.T(), "config1.exe", merged.Build.NLRC.Path)
}

// TestBuildOptionsStruct tests the BuildOptions struct
func (suite *OptionsTestSuite) TestBuildOptionsStruct() {
	opts := &BuildOptions{
		SourceFiles: []string{"file1.axs", "file2.axi"},
		CFGFiles:    []string{"config.cfg"},
		IncludePath: []string{"include1", "include2"},
		ModulePath:  []string{"module1", "module2"},
		LibraryPath: []string{"lib1", "lib2"},
		OutputPath:  "output/",
		NLRCPath:    "nlrc.exe",
		ShellPath:   "shell.exe",
		All:         true,
		Verbose:     true,
	}

	assert.Len(suite.T(), opts.SourceFiles, 2)
	assert.Len(suite.T(), opts.IncludePath, 2)
	assert.Len(suite.T(), opts.ModulePath, 2)
	assert.Len(suite.T(), opts.LibraryPath, 2)
	assert.Equal(suite.T(), "output/", opts.OutputPath)
	assert.Equal(suite.T(), "nlrc.exe", opts.NLRCPath)
	assert.Equal(suite.T(), "shell.exe", opts.ShellPath)
	assert.True(suite.T(), opts.All)
	assert.True(suite.T(), opts.Verbose)
}

// TestCLIOptionsStruct tests the CLIOptions struct
func (suite *OptionsTestSuite) TestCLIOptionsStruct() {
	cliOpts := &CLIOptions{
		SourceFiles: []string{"cli1.axs", "cli2.axi"},
		CFGFiles:    []string{"cli.cfg"},
		IncludePath: []string{"cli/include"},
		ModulePath:  []string{"cli/module"},
		LibraryPath: []string{"cli/lib"},
		OutputPath:  "cli/output",
		All:         false,
		Verbose:     true,
	}

	assert.Len(suite.T(), cliOpts.SourceFiles, 2)
	assert.Len(suite.T(), cliOpts.IncludePath, 1)
	assert.Equal(suite.T(), "cli/output", cliOpts.OutputPath)
	assert.False(suite.T(), cliOpts.All)
	assert.True(suite.T(), cliOpts.Verbose)
}

// TestMergeConfigs_LaterConfigWins verifies that a later (higher-precedence)
// config's non-empty Path field overrides earlier configs, mirroring the
// default < global < local precedence used in LoadBuildOptions.
func (suite *OptionsTestSuite) TestMergeConfigs_LaterConfigWins() {
	defaultCfg := &config.Config{
		Build: config.BuildConfig{
			NLRC: config.NLRCConfig{Path: "default.exe"},
		},
	}
	globalCfg := &config.Config{
		Build: config.BuildConfig{
			NLRC: config.NLRCConfig{Path: "global.exe"},
		},
	}
	localCfg := &config.Config{
		Build: config.BuildConfig{
			NLRC: config.NLRCConfig{Path: "local.exe"},
		},
	}

	merged := mergeConfigs(defaultCfg, globalCfg, localCfg)

	// local is last — it must win
	assert.Equal(suite.T(), "local.exe", merged.Build.NLRC.Path)
}

// TestMergeConfigs_EmptyLaterDoesNotOverride verifies that an empty string field
// in a later config does not wipe out a value set by an earlier config.
func (suite *OptionsTestSuite) TestMergeConfigs_EmptyLaterDoesNotOverride() {
	config1 := &config.Config{
		Build: config.BuildConfig{
			NLRC: config.NLRCConfig{Path: "config1.exe"},
		},
	}
	config2 := &config.Config{} // Path is "" (zero value)

	merged := mergeConfigs(config1, config2)

	assert.Equal(suite.T(), "config1.exe", merged.Build.NLRC.Path)
}

// TestMergeConfigs_IncludePathAccumulates verifies that include paths from
// multiple configs are accumulated (not overwritten) in precedence order.
func (suite *OptionsTestSuite) TestMergeConfigs_IncludePathAccumulates() {
	config1 := &config.Config{
		Build: config.BuildConfig{
			NLRC: config.NLRCConfig{IncludePath: []string{"global/include"}},
		},
	}
	config2 := &config.Config{
		Build: config.BuildConfig{
			NLRC: config.NLRCConfig{IncludePath: []string{"local/include"}},
		},
	}

	merged := mergeConfigs(config1, config2)

	assert.Contains(suite.T(), merged.Build.NLRC.IncludePath, "global/include")
	assert.Contains(suite.T(), merged.Build.NLRC.IncludePath, "local/include")
	assert.Len(suite.T(), merged.Build.NLRC.IncludePath, 2)
}

// TestLoadBuildOptions_CLIPathsPrependedBeforeConfig verifies that CLI-supplied
// include paths appear before config-file paths in the merged result.
func (suite *OptionsTestSuite) TestLoadBuildOptions_CLIPathsPrependedBeforeConfig() {
	// Prevent any real global config from being discovered
	noGlobalDir := filepath.Join(suite.tempDir, "no_global")
	suite.Require().NoError(os.MkdirAll(noGlobalDir, 0o755))
	suite.T().Setenv("GENLINX_CONFIG_DIR", noGlobalDir)

	// Write a local config with an include path
	configContent := `{"build":{"nlrc":{"includePath":["config/include"]}}}`
	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(configContent),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadBuildOptions(&CLIOptions{IncludePath: []string{"cli/include"}})
	suite.Require().NoError(err)

	// Find positions — paths will be absolute after resolvePaths()
	cliIdx, cfgIdx := -1, -1
	for i, p := range opts.IncludePath {
		if strings.Contains(p, "cli") {
			cliIdx = i
		}
		if strings.Contains(p, "config") {
			cfgIdx = i
		}
	}

	assert.NotEqual(suite.T(), -1, cliIdx, "CLI include path should be present")
	assert.NotEqual(suite.T(), -1, cfgIdx, "config include path should be present")
	assert.Less(suite.T(), cliIdx, cfgIdx, "CLI path should appear before config path")
}

// TestOptionsTestSuite runs the test suite
func TestOptionsTestSuite(t *testing.T) {
	suite.Run(t, new(OptionsTestSuite))
}
