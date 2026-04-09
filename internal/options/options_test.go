package options

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/Norgate-AV/genlinx/internal/config"
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
	resolved, err := filepath.EvalSymlinks(tempDir)
	suite.Require().NoError(err)
	suite.tempDir = resolved
}

// TearDownTest cleans up the test environment
func (suite *OptionsTestSuite) TearDownTest() {
	_ = os.RemoveAll(suite.tempDir)
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
		SourceFiles:  []string{"relative/file.axs"},
		CFGFiles:     []string{"./config.cfg"},
		IncludePath:  []string{"./include"},
		ModulePath:   []string{"./module"},
		LibraryPath:  []string{"./lib"},
		OutputPath:   "./output",
		CompilerPath: "./nlrc.exe",
	}

	// Change to temp directory for relative path testing
	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck

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
	assert.True(suite.T(), filepath.IsAbs(opts.CompilerPath), "compiler path should be absolute: %s", opts.CompilerPath)
}

// TestLoadGlobalConfig tests global configuration loading
func (suite *OptionsTestSuite) TestLoadGlobalConfig() {
	// Redirect the global config search to an empty temp dir so the test is
	// not affected by any real config file present on the developer's machine.
	suite.T().Setenv("GENLINX_CONFIG_DIR", suite.tempDir)

	// Test loading global config
	result, err := LoadGlobalConfig()
	suite.Require().NoError(err)
	// Config should be empty since no global config exists in test environment
	assert.NotNil(suite.T(), result.Config)
	assert.False(suite.T(), result.Found)
}

// TestLoadGlobalConfig_YAML verifies that a global config.yaml file is recognised
// and its values are loaded correctly.
func (suite *OptionsTestSuite) TestLoadGlobalConfig_YAML() {
	globalDir := filepath.Join(suite.tempDir, "global_yaml")
	suite.Require().NoError(os.MkdirAll(globalDir, 0o755))
	suite.T().Setenv("GENLINX_CONFIG_DIR", globalDir)

	content := "compiler:\n  includePath:\n    - yaml/include\n"
	suite.Require().NoError(os.WriteFile(
		filepath.Join(globalDir, "config.yaml"),
		[]byte(content), 0o644,
	))

	result, err := LoadGlobalConfig()
	suite.Require().NoError(err)
	assert.True(suite.T(), result.Found)
	assert.Contains(suite.T(), result.Config.Compiler.IncludePath, filepath.FromSlash("yaml/include"))
}

// TestLoadGlobalConfig_YML verifies that a global config.yml file is recognised
// and its values are loaded correctly.
func (suite *OptionsTestSuite) TestLoadGlobalConfig_YML() {
	globalDir := filepath.Join(suite.tempDir, "global_yml")
	suite.Require().NoError(os.MkdirAll(globalDir, 0o755))
	suite.T().Setenv("GENLINX_CONFIG_DIR", globalDir)

	content := "compiler:\n  includePath:\n    - yml/include\n"
	suite.Require().NoError(os.WriteFile(
		filepath.Join(globalDir, "config.yml"),
		[]byte(content), 0o644,
	))

	result, err := LoadGlobalConfig()
	suite.Require().NoError(err)
	assert.True(suite.T(), result.Found)
	assert.Contains(suite.T(), result.Config.Compiler.IncludePath, filepath.FromSlash("yml/include"))
}

// TestLoadLocalConfig tests local configuration loading with find-up
func (suite *OptionsTestSuite) TestLoadLocalConfig() {
	// Create a config file in the temp directory
	configFile := filepath.Join(suite.tempDir, ".genlinxrc.json")
	configContent := `{
		"compiler": {
			"includePath": ["test/include"]
		}
	}`

	err := os.WriteFile(configFile, []byte(configContent), 0o644)
	suite.Require().NoError(err)

	// Change to temp directory
	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer func() { _ = os.Chdir(oldWd) }()

	err = os.Chdir(suite.tempDir)
	suite.Require().NoError(err)

	// Test loading local config
	result, err := LoadLocalConfig()
	suite.Require().NoError(err)
	assert.NotNil(suite.T(), result.Config)
	assert.True(suite.T(), result.Found)
	assert.Contains(suite.T(), result.Config.Compiler.IncludePath, filepath.FromSlash("test/include"))
}

// TestLoadLocalConfig_Formats verifies that every supported filename and format
// variant for the local config is discovered and correctly parsed.
func (suite *OptionsTestSuite) TestLoadLocalConfig_Formats() {
	tests := []struct {
		filename string
		content  string
		wantPath string
	}{
		{
			filename: ".genlinxrc.yaml",
			content:  "compiler:\n  includePath:\n    - yaml/include\n",
			wantPath: "yaml/include",
		},
		{
			filename: ".genlinxrc.yml",
			content:  "compiler:\n  includePath:\n    - yml/include\n",
			wantPath: "yml/include",
		},
		{
			filename: ".genlinx.json",
			content:  `{"compiler":{"includePath":["genlinx/include"]}}`,
			wantPath: "genlinx/include",
		},
		{
			filename: ".genlinx.yaml",
			content:  "compiler:\n  includePath:\n    - genlinxyaml/include\n",
			wantPath: "genlinxyaml/include",
		},
		{
			filename: ".genlinx.yml",
			content:  "compiler:\n  includePath:\n    - genlinxyml/include\n",
			wantPath: "genlinxyml/include",
		},
	}

	for _, tc := range tests {
		suite.Run(tc.filename, func() {
			testDir := filepath.Join(suite.tempDir, tc.filename+"_test")
			suite.Require().NoError(os.MkdirAll(testDir, 0o755))
			suite.Require().NoError(os.WriteFile(
				filepath.Join(testDir, tc.filename),
				[]byte(tc.content), 0o644,
			))

			oldWd, err := os.Getwd()
			suite.Require().NoError(err)
			defer os.Chdir(oldWd) //nolint:errcheck
			suite.Require().NoError(os.Chdir(testDir))

			result, err := LoadLocalConfig()
			suite.Require().NoError(err)
			assert.True(suite.T(), result.Found, "config not found for %s", tc.filename)
			assert.Contains(suite.T(), result.Config.Compiler.IncludePath,
				filepath.FromSlash(tc.wantPath), "wrong include path for %s", tc.filename)
		})
	}
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
		"compiler": {
			"modulePath": ["findup/module"]
		}
	}`

	err = os.WriteFile(configFile, []byte(configContent), 0o644)
	suite.Require().NoError(err)

	// Change to child directory
	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck

	err = os.Chdir(childDir)
	suite.Require().NoError(err)

	// Test that config is found via find-up
	result, err := LoadLocalConfig()
	suite.Require().NoError(err)
	assert.NotNil(suite.T(), result.Config)
	assert.True(suite.T(), result.Found)
	assert.Contains(suite.T(), result.Config.Compiler.ModulePath, filepath.FromSlash("findup/module"))
}

// TestLoadLocalConfigNoConfig tests when no local config is found
func (suite *OptionsTestSuite) TestLoadLocalConfigNoConfig() {
	// Change to temp directory (no config files)
	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck

	err = os.Chdir(suite.tempDir)
	suite.Require().NoError(err)

	// Test loading when no config exists
	result, err := LoadLocalConfig()
	suite.Require().NoError(err)
	assert.NotNil(suite.T(), result.Config)
	assert.False(suite.T(), result.Found)
	// Should return empty config
	assert.Empty(suite.T(), result.Config.Compiler.IncludePath)
}

// TestMergeConfigs tests configuration merging
func (suite *OptionsTestSuite) TestMergeConfigs() {
	config1 := &config.Config{
		Compiler: config.CompilerConfig{
			IncludePath: []string{"config1/include"},
			Path:        "config1.exe",
		},
	}

	config2 := &config.Config{
		Compiler: config.CompilerConfig{
			IncludePath: []string{"config2/include"},
			ModulePath:  []string{"config2/module"},
		},
	}

	merged := mergeConfigs(config1, config2)

	// Should contain items from both configs
	assert.Contains(suite.T(), merged.Compiler.IncludePath, "config1/include")
	assert.Contains(suite.T(), merged.Compiler.IncludePath, "config2/include")
	assert.Contains(suite.T(), merged.Compiler.ModulePath, "config2/module")
	assert.Equal(suite.T(), "config1.exe", merged.Compiler.Path)
}

// TestBuildOptionsStruct tests the BuildOptions struct
func (suite *OptionsTestSuite) TestBuildOptionsStruct() {
	opts := &BuildOptions{
		SourceFiles:  []string{"file1.axs", "file2.axi"},
		CFGFiles:     []string{"config.cfg"},
		IncludePath:  []string{"include1", "include2"},
		ModulePath:   []string{"module1", "module2"},
		LibraryPath:  []string{"lib1", "lib2"},
		OutputPath:   "output/",
		CompilerPath: "nlrc.exe",
		All:          true,
		Verbose:      true,
	}

	assert.Len(suite.T(), opts.SourceFiles, 2)
	assert.Len(suite.T(), opts.IncludePath, 2)
	assert.Len(suite.T(), opts.ModulePath, 2)
	assert.Len(suite.T(), opts.LibraryPath, 2)
	assert.Equal(suite.T(), "output/", opts.OutputPath)
	assert.Equal(suite.T(), "nlrc.exe", opts.CompilerPath)
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
		Compiler: config.CompilerConfig{Path: "default.exe"},
	}

	globalCfg := &config.Config{
		Compiler: config.CompilerConfig{Path: "global.exe"},
	}

	localCfg := &config.Config{
		Compiler: config.CompilerConfig{Path: "local.exe"},
	}

	merged := mergeConfigs(defaultCfg, globalCfg, localCfg)

	// local is last — it must win
	assert.Equal(suite.T(), "local.exe", merged.Compiler.Path)
}

// TestMergeConfigs_EmptyLaterDoesNotOverride verifies that an empty string field
// in a later config does not wipe out a value set by an earlier config.
func (suite *OptionsTestSuite) TestMergeConfigs_EmptyLaterDoesNotOverride() {
	config1 := &config.Config{
		Compiler: config.CompilerConfig{Path: "config1.exe"},
	}

	config2 := &config.Config{} // Path is "" (zero value)

	merged := mergeConfigs(config1, config2)

	assert.Equal(suite.T(), "config1.exe", merged.Compiler.Path)
}

// TestMergeConfigs_IncludePathAccumulates verifies that include paths from
// multiple configs are accumulated (not overwritten) in precedence order.
func (suite *OptionsTestSuite) TestMergeConfigs_IncludePathAccumulates() {
	config1 := &config.Config{
		Compiler: config.CompilerConfig{IncludePath: []string{"global/include"}},
	}

	config2 := &config.Config{
		Compiler: config.CompilerConfig{IncludePath: []string{"local/include"}},
	}

	merged := mergeConfigs(config1, config2)

	assert.Contains(suite.T(), merged.Compiler.IncludePath, "global/include")
	assert.Contains(suite.T(), merged.Compiler.IncludePath, "local/include")
	assert.Len(suite.T(), merged.Compiler.IncludePath, 2)
	// Later config (config2 / local) paths must be prepended before earlier config (config1 / global)
	assert.Equal(suite.T(), "local/include", merged.Compiler.IncludePath[0], "later config path should come first")
	assert.Equal(suite.T(), "global/include", merged.Compiler.IncludePath[1])
}

// TestMergeConfigs_DeduplicatesPaths verifies that identical paths present in
// multiple configs only appear once in the merged result, mirroring
// mergician's dedupArrays: true behaviour.
func (suite *OptionsTestSuite) TestMergeConfigs_DeduplicatesPaths() {
	sharedPath := "C:/Common/AMXShare/AXIs"
	config1 := &config.Config{
		Compiler: config.CompilerConfig{IncludePath: []string{sharedPath}},
	}

	config2 := &config.Config{
		Compiler: config.CompilerConfig{IncludePath: []string{sharedPath}},
	}

	merged := mergeConfigs(config1, config2)

	assert.Len(suite.T(), merged.Compiler.IncludePath, 1, "shared path should only appear once")
	assert.Equal(suite.T(), sharedPath, merged.Compiler.IncludePath[0])
}

// TestMergeConfigs_DeduplicatesNormalized verifies that the same physical
// path written with different separators (e.g. forward vs backslash) is
// treated as a duplicate after NormalizeConfigPaths is applied by the loaders.
// This catches the bug where global config (forward-slash) and default config
// (OS-normalized backslash) would both survive dedup and produce duplicate
// compiler flags.
func (suite *OptionsTestSuite) TestMergeConfigs_DeduplicatesNormalized() {
	rawPath := "C:/Program Files (x86)/Common Files/AMXShare/AXIs"

	// Simulate the default config — already normalized (as LoadDefaultConfig produces).
	defaultCfg := &config.Config{
		Compiler: config.CompilerConfig{IncludePath: []string{rawPath}},
	}

	config.NormalizeConfigPaths(defaultCfg)

	// Simulate a global config loaded from a JSON file with forward slashes,
	// then normalized by loadGlobalConfig.
	globalCfg := &config.Config{
		Compiler: config.CompilerConfig{IncludePath: []string{rawPath}},
	}

	config.NormalizeConfigPaths(globalCfg)

	// After normalization both have the same OS-canonical string, so merging
	// must deduplicate down to a single entry.
	merged := mergeConfigs(defaultCfg, globalCfg)

	assert.Len(suite.T(), merged.Compiler.IncludePath, 1,
		"same path must deduplicate after normalization")
}

// TestMergeConfigsWithPresence_BoolsRespected verifies that boolean fields are
// only applied from an override when the corresponding presence flag is set,
// preventing JSON zero-value ambiguity.
func (suite *OptionsTestSuite) TestMergeConfigsWithPresence_BoolsRespected() {
	base := &config.Config{
		CFG: config.CFGConfig{
			OutputLogConsoleOption: true,  // default: true
			BuildWithSource:        false, // default: false
		},
		Archive: config.ArchiveConfig{
			IncludeCompiledSourceFiles: true, // default: true
		},
	}

	override := &config.Config{
		CFG: config.CFGConfig{
			OutputLogConsoleOption: false, // wants to override to false
			BuildWithSource:        true,  // wants to override to true
		},
		Archive: config.ArchiveConfig{
			IncludeCompiledSourceFiles: false, // wants to override to false
		},
	}

	// Without presence flags — boolean overrides must NOT be applied.
	merged := mergeConfigsWithPresence(base, configMergeInput{cfg: override, presence: boolPresence{}})
	suite.True(merged.CFG.OutputLogConsoleOption, "bool without presence flag should not be overridden")
	suite.False(merged.CFG.BuildWithSource, "bool without presence flag should not be overridden")
	suite.True(merged.Archive.IncludeCompiledSourceFiles, "bool without presence flag should not be overridden")

	// With presence flags — boolean overrides MUST be applied.
	merged = mergeConfigsWithPresence(base, configMergeInput{
		cfg: override,
		presence: boolPresence{
			cfgOutputLogConsoleOption: true,
			cfgBuildWithSource:        true,
			archiveIncludeCompiledSrc: true,
		},
	})
	suite.False(merged.CFG.OutputLogConsoleOption, "bool with presence flag should be applied")
	suite.True(merged.CFG.BuildWithSource, "bool with presence flag should be applied")
	suite.False(merged.Archive.IncludeCompiledSourceFiles, "bool with presence flag should be applied")
}

// TestLoadMergedConfig_NoConfigs verifies that when neither a global nor a
// local config file is present, loadMergedConfig returns the built-in
// defaults and ConfigLoadInfo reflects both as not-found.
func (suite *OptionsTestSuite) TestLoadMergedConfig_NoConfigs() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global"), 0o755))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	merged, info, err := LoadMergedConfig()

	suite.Require().NoError(err)
	suite.Require().NotNil(merged)
	suite.Require().NotNil(info)

	assert.True(suite.T(), info.DefaultLoaded)
	assert.False(suite.T(), info.GlobalResult.Found)
	assert.False(suite.T(), info.LocalResult.Found)

	// Default compiler path must be present
	assert.NotEmpty(suite.T(), merged.Compiler.Path)
	// Default include / module / library paths must be non-empty
	assert.NotEmpty(suite.T(), merged.Compiler.IncludePath)
	assert.NotEmpty(suite.T(), merged.Compiler.ModulePath)
	assert.NotEmpty(suite.T(), merged.Compiler.LibraryPath)
}

// TestLoadMergedConfig_LocalPrependsBeforeDefault verifies that a local config's
// paths appear before the default paths in the merged result and that no paths
// are duplicated.
func (suite *OptionsTestSuite) TestLoadMergedConfig_LocalPrependsBeforeDefault() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"compiler":{"includePath":["local/include"]}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	merged, info, err := LoadMergedConfig()
	suite.Require().NoError(err)

	assert.True(suite.T(), info.LocalResult.Found)

	// local/include must appear first
	assert.True(suite.T(), strings.Contains(merged.Compiler.IncludePath[0], "local"),
		"local path must be first, got: %v", merged.Compiler.IncludePath)

	// No duplicates
	seen := make(map[string]int)
	for _, p := range merged.Compiler.IncludePath {
		seen[p]++
	}

	for p, count := range seen {
		assert.Equal(suite.T(), 1, count, "path %q appears %d times", p, count)
	}
}

// TestLoadMergedConfig_GlobalPrependsBeforeDefault verifies that a global
// config's paths appear before the default paths and are deduplicated.
func (suite *OptionsTestSuite) TestLoadMergedConfig_GlobalPrependsBeforeDefault() {
	globalDir := filepath.Join(suite.tempDir, "global")
	suite.Require().NoError(os.MkdirAll(globalDir, 0o755))
	suite.T().Setenv("GENLINX_CONFIG_DIR", globalDir)

	suite.Require().NoError(os.WriteFile(
		filepath.Join(globalDir, "config.json"),
		[]byte(`{"compiler":{"includePath":["global/include"]}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	merged, info, err := LoadMergedConfig()
	suite.Require().NoError(err)

	assert.True(suite.T(), info.GlobalResult.Found)
	assert.False(suite.T(), info.LocalResult.Found)

	// global/include must appear first
	assert.True(suite.T(), strings.Contains(merged.Compiler.IncludePath[0], "global"),
		"global path must be first, got: %v", merged.Compiler.IncludePath)

	// No duplicates
	seen := make(map[string]int)
	for _, p := range merged.Compiler.IncludePath {
		seen[p]++
	}

	for p, count := range seen {
		assert.Equal(suite.T(), 1, count, "path %q appears %d times", p, count)
	}
}

// TestLoadMergedConfig_AllThreeLayers verifies the full default < global < local
// precedence chain: local paths appear first, then global, then default — with
// no duplicates across any layer.
func (suite *OptionsTestSuite) TestLoadMergedConfig_AllThreeLayers() {
	globalDir := filepath.Join(suite.tempDir, "global")
	suite.Require().NoError(os.MkdirAll(globalDir, 0o755))
	suite.T().Setenv("GENLINX_CONFIG_DIR", globalDir)

	suite.Require().NoError(os.WriteFile(
		filepath.Join(globalDir, "config.json"),
		[]byte(`{"compiler":{"includePath":["global/include"]}}`),
		0o644,
	))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"compiler":{"includePath":["local/include"]}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	merged, info, err := LoadMergedConfig()
	suite.Require().NoError(err)

	assert.True(suite.T(), info.GlobalResult.Found)
	assert.True(suite.T(), info.LocalResult.Found)

	paths := merged.Compiler.IncludePath
	localIdx, globalIdx, defaultIdx := -1, -1, -1
	for i, p := range paths {
		switch {
		case strings.Contains(p, "local"):
			localIdx = i
		case strings.Contains(p, "global"):
			globalIdx = i
		case strings.Contains(p, "AMXShare"):
			defaultIdx = i
		}
	}

	assert.NotEqual(suite.T(), -1, localIdx, "local path should be present")
	assert.NotEqual(suite.T(), -1, globalIdx, "global path should be present")
	assert.NotEqual(suite.T(), -1, defaultIdx, "default path should be present")
	assert.Less(suite.T(), localIdx, globalIdx, "local must come before global")
	assert.Less(suite.T(), globalIdx, defaultIdx, "global must come before default")

	// No duplicates
	seen := make(map[string]int)
	for _, p := range paths {
		seen[p]++
	}

	for p, count := range seen {
		assert.Equal(suite.T(), 1, count, "path %q appears %d times", p, count)
	}
}

// TestLoadMergedConfig_ResolvesRelativePaths verifies that relative paths
// specified in a local config (e.g. "./include") are resolved to absolute
// paths against the CWD at merge time.
func (suite *OptionsTestSuite) TestLoadMergedConfig_ResolvesRelativePaths() {
	noGlobalDir := filepath.Join(suite.tempDir, "no_global_rel")
	suite.Require().NoError(os.MkdirAll(noGlobalDir, 0o755))
	suite.T().Setenv("GENLINX_CONFIG_DIR", noGlobalDir)

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"compiler":{"includePath":["./include"],"modulePath":["./module"]}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	merged, _, err := LoadMergedConfig()
	suite.Require().NoError(err)

	wantInclude := filepath.Join(suite.tempDir, "include")
	wantModule := filepath.Join(suite.tempDir, "module")

	assert.Contains(suite.T(), merged.Compiler.IncludePath, wantInclude,
		"relative ./include must resolve to absolute path")
	assert.Contains(suite.T(), merged.Compiler.ModulePath, wantModule,
		"relative ./module must resolve to absolute path")

	// Every path in the merged result must be absolute
	for _, p := range merged.Compiler.IncludePath {
		assert.True(suite.T(), filepath.IsAbs(p), "IncludePath must be absolute: %s", p)
	}

	for _, p := range merged.Compiler.ModulePath {
		assert.True(suite.T(), filepath.IsAbs(p), "ModulePath must be absolute: %s", p)
	}
}

// TestLoadMergedConfig_AllPathsAbsolute verifies that even when no local or
// global config exists, the default paths returned by LoadMergedConfig are
// all absolute.
func (suite *OptionsTestSuite) TestLoadMergedConfig_AllPathsAbsolute() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_abs"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_abs"), 0o755))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	merged, _, err := LoadMergedConfig()
	suite.Require().NoError(err)

	allPaths := []string{merged.Compiler.Path}
	allPaths = append(allPaths, merged.Compiler.IncludePath...)
	allPaths = append(allPaths, merged.Compiler.ModulePath...)
	allPaths = append(allPaths, merged.Compiler.LibraryPath...)
	allPaths = append(allPaths, merged.Archive.ExtraFileSearchLocations...)

	for _, p := range allPaths {
		if p == "" {
			continue
		}

		assert.True(suite.T(), filepath.IsAbs(p), "all merged paths must be absolute: %s", p)
	}
}

// TestLoadBuildOptions_CLIPathsPrependedBeforeConfig verifies that CLI-supplied
// include paths appear before config-file paths in the merged result.
func (suite *OptionsTestSuite) TestLoadBuildOptions_CLIPathsPrependedBeforeConfig() {
	// Prevent any real global config from being discovered
	noGlobalDir := filepath.Join(suite.tempDir, "no_global")
	suite.Require().NoError(os.MkdirAll(noGlobalDir, 0o755))
	suite.T().Setenv("GENLINX_CONFIG_DIR", noGlobalDir)

	// Write a local config with an include path
	configContent := `{"compiler":{"includePath":["config/include"]}}`
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

// TestLoadBuildOptions_NilCLI verifies that nil CLI options return a result
// built purely from the merged (default + global + local) config.
func (suite *OptionsTestSuite) TestLoadBuildOptions_NilCLI() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_build_nil"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_build_nil"), 0o755))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, info, err := LoadBuildOptions(nil)
	suite.Require().NoError(err)
	suite.Require().NotNil(opts)
	suite.Require().NotNil(info)
	suite.True(info.DefaultLoaded)
	suite.NotEmpty(opts.CompilerPath)
	suite.NotEmpty(opts.IncludePath)
}

// TestLoadBuildOptions_CompilerPathFromConfig verifies that a local config's
// compiler.path overrides the built-in default.
func (suite *OptionsTestSuite) TestLoadBuildOptions_CompilerPathFromConfig() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_build_nlrc"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_build_nlrc"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"compiler":{"path":"C:/custom/NLRC.exe"}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadBuildOptions(nil)
	suite.Require().NoError(err)
	suite.Contains(opts.CompilerPath, "custom",
		"local config compiler.path must override the built-in default")
}

// TestLoadBuildOptions_ModulePathPrependedBeforeConfig verifies that CLI-supplied
// module paths appear before config-file paths in the merged result.
func (suite *OptionsTestSuite) TestLoadBuildOptions_ModulePathPrependedBeforeConfig() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_build_mod"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_build_mod"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"compiler":{"modulePath":["config/module"]}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadBuildOptions(&CLIOptions{ModulePath: []string{"cli/module"}})
	suite.Require().NoError(err)

	cliIdx, cfgIdx := -1, -1
	for i, p := range opts.ModulePath {
		if strings.Contains(p, "cli") {
			cliIdx = i
		}

		if strings.Contains(p, "config") {
			cfgIdx = i
		}
	}

	suite.NotEqual(-1, cliIdx, "CLI module path should be present")
	suite.NotEqual(-1, cfgIdx, "config module path should be present")
	suite.Less(cliIdx, cfgIdx, "CLI module path must appear before config path")
}

// TestLoadBuildOptions_LibraryPathPrependedBeforeConfig verifies that CLI-supplied
// library paths appear before config-file paths in the merged result.
func (suite *OptionsTestSuite) TestLoadBuildOptions_LibraryPathPrependedBeforeConfig() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_build_lib"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_build_lib"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"compiler":{"libraryPath":["config/lib"]}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadBuildOptions(&CLIOptions{LibraryPath: []string{"cli/lib"}})
	suite.Require().NoError(err)

	cliIdx, cfgIdx := -1, -1
	for i, p := range opts.LibraryPath {
		if strings.Contains(p, "cli") {
			cliIdx = i
		}

		if strings.Contains(p, "config") {
			cfgIdx = i
		}
	}

	suite.NotEqual(-1, cliIdx, "CLI library path should be present")
	suite.NotEqual(-1, cfgIdx, "config library path should be present")
	suite.Less(cliIdx, cfgIdx, "CLI library path must appear before config path")
}

// TestLoadBuildOptions_OutputPathOverride verifies that a CLI output path
// overrides the (empty) config value and is resolved to an absolute path.
func (suite *OptionsTestSuite) TestLoadBuildOptions_OutputPathOverride() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_build_out"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_build_out"), 0o755))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadBuildOptions(&CLIOptions{OutputPath: "cli/output"})
	suite.Require().NoError(err)
	suite.Contains(opts.OutputPath, "cli")
	suite.True(filepath.IsAbs(opts.OutputPath), "output path must be resolved to an absolute path")
}

// TestLoadBuildOptions_CLIAllOverridesConfig verifies that CLI All:true is
// applied even when the config does not set it.
func (suite *OptionsTestSuite) TestLoadBuildOptions_CLIAllOverridesConfig() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_build_cliall"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_build_cliall"), 0o755))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadBuildOptions(&CLIOptions{All: true})
	suite.Require().NoError(err)
	suite.True(opts.All)
}

// ---------------------------------------------------------------------------
// LoadArchiveOptions
// ---------------------------------------------------------------------------

// TestLoadArchiveOptions_NilCLI verifies that passing nil CLI options returns
// a result built purely from the merged (default + global + local) config.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_NilCLI() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_arch"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_arch"), 0o755))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, info, err := LoadArchiveOptions(nil)
	suite.Require().NoError(err)
	suite.Require().NotNil(opts)
	suite.Require().NotNil(info)
	suite.True(info.DefaultLoaded)
}

// TestLoadArchiveOptions_OutputFileSuffix verifies that a CLI OutputFileSuffix
// overrides the config value.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_OutputFileSuffix() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_arch2"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_arch2"), 0o755))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	cliOpts := &ArchiveCLIOptions{
		OutputFileSuffix:  ".custom",
		ExplicitBoolFlags: map[string]bool{},
	}

	opts, _, err := LoadArchiveOptions(cliOpts)
	suite.Require().NoError(err)
	suite.Equal(".custom", opts.OutputFileSuffix)
}

// TestLoadArchiveOptions_BoolFlagOverride verifies that a ExplicitBoolFlags flag
// overrides the config default for the corresponding boolean option.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_BoolFlagOverride() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_arch3"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_arch3"), 0o755))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	cliOpts := &ArchiveCLIOptions{
		ExplicitBoolFlags: map[string]bool{
			"no-include-compiled-source-files":  true,
			"include-compiled-module-files":     true,
			"no-include-files-not-in-workspace": true,
		},
	}

	opts, _, err := LoadArchiveOptions(cliOpts)
	suite.Require().NoError(err)
	suite.False(opts.IncludeCompiledSourceFiles)
	suite.True(opts.IncludeCompiledModuleFiles)
	suite.False(opts.IncludeFilesNotInWorkspace)
}

// TestLoadArchiveOptions_ExtraSearchLocationsPrepended verifies that CLI
// ExtraFileSearchLocations are prepended before config paths.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_ExtraSearchLocationsPrepended() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_arch4"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_arch4"), 0o755))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	cliOpts := &ArchiveCLIOptions{
		ExtraFileSearchLocations: []string{"cli/search"},
		ExplicitBoolFlags:        map[string]bool{},
	}

	opts, _, err := LoadArchiveOptions(cliOpts)
	suite.Require().NoError(err)
	suite.Require().NotEmpty(opts.ExtraFileSearchLocations)
	suite.Equal("cli/search", opts.ExtraFileSearchLocations[0])
}

// TestLoadArchiveOptions_ConfigFileBoolPreservedWhenNoFlag verifies that boolean
// values set in a local config file are preserved when no CLI flag is passed —
// an empty ExplicitBoolFlags map must not reset config-file values back to built-in defaults.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_ConfigFileBoolPreservedWhenNoFlag() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_arch_preserve"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_arch_preserve"), 0o755))

	// Flip all three flags to the opposite of the built-in defaults
	// (defaults: includeCompiledSourceFiles=true, includeCompiledModuleFiles=true, includeFilesNotInWorkspace=true).
	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"archive":{"includeCompiledSourceFiles":false,"includeCompiledModuleFiles":false,"includeFilesNotInWorkspace":false}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadArchiveOptions(&ArchiveCLIOptions{ExplicitBoolFlags: map[string]bool{}})
	suite.Require().NoError(err)
	suite.False(opts.IncludeCompiledSourceFiles,
		"config file includeCompiledSourceFiles:false must be preserved when no CLI flag is set")
	suite.False(opts.IncludeCompiledModuleFiles,
		"config file includeCompiledModuleFiles:false must be preserved when no CLI flag is set")
	suite.False(opts.IncludeFilesNotInWorkspace,
		"config file includeFilesNotInWorkspace:false must be preserved when no CLI flag is set")
}

// TestLoadArchiveOptions_CLIOverridesConfigFile_IncludeCompiledSource verifies
// that --include-compiled-source-files overrides a config file that sets
// includeCompiledSourceFiles:false.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_CLIOverridesConfigFile_IncludeCompiledSource() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_arch_ics"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_arch_ics"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"archive":{"includeCompiledSourceFiles":false}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadArchiveOptions(&ArchiveCLIOptions{
		ExplicitBoolFlags: map[string]bool{"include-compiled-source-files": true},
	})
	suite.Require().NoError(err)
	suite.True(opts.IncludeCompiledSourceFiles,
		"--include-compiled-source-files must override config file includeCompiledSourceFiles:false")
}

// TestLoadArchiveOptions_CLIOverridesConfigFile_IncludeCompiledModule verifies
// that --include-compiled-module-files overrides a config file that sets
// includeCompiledModuleFiles:false.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_CLIOverridesConfigFile_IncludeCompiledModule() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_arch_icm"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_arch_icm"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"archive":{"includeCompiledModuleFiles":false}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadArchiveOptions(&ArchiveCLIOptions{
		ExplicitBoolFlags: map[string]bool{"include-compiled-module-files": true},
	})
	suite.Require().NoError(err)
	suite.True(opts.IncludeCompiledModuleFiles,
		"--include-compiled-module-files must override config file includeCompiledModuleFiles:false")
}

// TestLoadArchiveOptions_CLIOverridesConfigFile_IncludeFilesNotInWorkspace
// verifies that --include-files-not-in-workspace overrides a config file that
// sets includeFilesNotInWorkspace:false.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_CLIOverridesConfigFile_IncludeFilesNotInWorkspace() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_arch_ifw"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_arch_ifw"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"archive":{"includeFilesNotInWorkspace":false}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadArchiveOptions(&ArchiveCLIOptions{
		ExplicitBoolFlags: map[string]bool{"include-files-not-in-workspace": true},
	})
	suite.Require().NoError(err)
	suite.True(opts.IncludeFilesNotInWorkspace,
		"--include-files-not-in-workspace must override config file includeFilesNotInWorkspace:false")
}

// TestLoadArchiveOptions_CLIOverridesConfigFile_NoIncludeFilesNotInWorkspace
// verifies that --no-include-files-not-in-workspace overrides a config file
// that sets includeFilesNotInWorkspace:true.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_CLIOverridesConfigFile_NoIncludeFilesNotInWorkspace() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_arch_nofw"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_arch_nofw"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"archive":{"includeFilesNotInWorkspace":true}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadArchiveOptions(&ArchiveCLIOptions{
		ExplicitBoolFlags: map[string]bool{"no-include-files-not-in-workspace": true},
	})
	suite.Require().NoError(err)
	suite.False(opts.IncludeFilesNotInWorkspace,
		"--no-include-files-not-in-workspace must override config file includeFilesNotInWorkspace:true")
}

// ---------------------------------------------------------------------------
// LoadCfgOptions
// ---------------------------------------------------------------------------

// TestLoadCfgOptions_NilCLI verifies that passing nil CLI options returns the
// merged config default.
func (suite *OptionsTestSuite) TestLoadCfgOptions_NilCLI() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_cfg"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_cfg"), 0o755))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, info, err := LoadCfgOptions(nil)
	suite.Require().NoError(err)
	suite.Require().NotNil(opts)
	suite.Require().NotNil(info)
	suite.True(info.DefaultLoaded)
}

// TestLoadCfgOptions_CLIOverrides verifies that string CLI fields override
// the config values and ExplicitBoolFlags boolean flags flip the corresponding booleans.
func (suite *OptionsTestSuite) TestLoadCfgOptions_CLIOverrides() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_cfg2"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_cfg2"), 0o755))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	cliOpts := &CfgCLIOptions{
		RootDirectory:       "/my/root",
		OutputFileSuffix:    ".cfg_custom",
		OutputLogFileSuffix: ".log_custom",
		OutputLogFileOption: "file",
		ExplicitBoolFlags: map[string]bool{
			"output-log-console-option":    true,
			"build-with-debug-information": true,
			"no-build-with-source":         true,
		},
	}

	opts, _, err := LoadCfgOptions(cliOpts)
	suite.Require().NoError(err)
	suite.Equal("/my/root", opts.RootDirectory)
	suite.Equal(".cfg_custom", opts.OutputFileSuffix)
	suite.Equal(".log_custom", opts.OutputLogFileSuffix)
	suite.Equal("file", opts.OutputLogFileOption)
	suite.True(opts.OutputLogConsoleOption)
	suite.True(opts.BuildWithDebugInformation)
	suite.False(opts.BuildWithSource)
}

// TestLoadCfgOptions_ConfigFileBoolPreservedWhenNoFlag verifies that boolean
// values set in a local config file are preserved when no CLI flag is passed —
// an empty ExplicitBoolFlags map must not reset config-file values back to built-in defaults.
func (suite *OptionsTestSuite) TestLoadCfgOptions_ConfigFileBoolPreservedWhenNoFlag() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_cfg_preserve"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_cfg_preserve"), 0o755))

	// Flip all three booleans to the opposite of the built-in defaults
	// (defaults: buildWithSource=false, buildWithDebugInformation=false, outputLogConsoleOption=true).
	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"cfg":{"buildWithSource":true,"buildWithDebugInformation":true,"outputLogConsoleOption":false}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadCfgOptions(&CfgCLIOptions{ExplicitBoolFlags: map[string]bool{}})
	suite.Require().NoError(err)
	suite.True(opts.BuildWithSource,
		"config file buildWithSource:true must be preserved when no CLI flag is set")
	suite.True(opts.BuildWithDebugInformation,
		"config file buildWithDebugInformation:true must be preserved when no CLI flag is set")
	suite.False(opts.OutputLogConsoleOption,
		"config file outputLogConsoleOption:false must be preserved when no CLI flag is set")
}

// TestLoadCfgOptions_CLIOverridesConfigFile_NoBuildWithSource verifies that
// --no-build-with-source overrides a config file that sets buildWithSource:true.
func (suite *OptionsTestSuite) TestLoadCfgOptions_CLIOverridesConfigFile_NoBuildWithSource() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_cfg_bws"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_cfg_bws"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"cfg":{"buildWithSource":true}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadCfgOptions(&CfgCLIOptions{
		ExplicitBoolFlags: map[string]bool{"no-build-with-source": true},
	})
	suite.Require().NoError(err)
	suite.False(opts.BuildWithSource,
		"--no-build-with-source must override config file buildWithSource:true")
}

// TestLoadCfgOptions_CLIOverridesConfigFile_BuildWithDebug verifies that
// --build-with-debug-information overrides a config file that sets
// buildWithDebugInformation:false.
func (suite *OptionsTestSuite) TestLoadCfgOptions_CLIOverridesConfigFile_BuildWithDebug() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_cfg_bwd"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_cfg_bwd"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"cfg":{"buildWithDebugInformation":false}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadCfgOptions(&CfgCLIOptions{
		ExplicitBoolFlags: map[string]bool{"build-with-debug-information": true},
	})
	suite.Require().NoError(err)
	suite.True(opts.BuildWithDebugInformation,
		"--build-with-debug-information must override config file buildWithDebugInformation:false")
}

// TestLoadCfgOptions_CLIOverridesConfigFile_OutputLogConsole verifies that
// --output-log-console-option overrides a config file that sets
// outputLogConsoleOption:false.
func (suite *OptionsTestSuite) TestLoadCfgOptions_CLIOverridesConfigFile_OutputLogConsole() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_cfg_olco"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_cfg_olco"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"cfg":{"outputLogConsoleOption":false}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadCfgOptions(&CfgCLIOptions{
		ExplicitBoolFlags: map[string]bool{"output-log-console-option": true},
	})
	suite.Require().NoError(err)
	suite.True(opts.OutputLogConsoleOption,
		"--output-log-console-option must override config file outputLogConsoleOption:false")
}

// TestLoadCfgOptions_CLIIncludePathPrependedBeforeConfig verifies that
// CLI include paths appear before config-file paths in the merged result.
func (suite *OptionsTestSuite) TestLoadCfgOptions_CLIIncludePathPrependedBeforeConfig() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_cfg_inc"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_cfg_inc"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"compiler":{"includePath":["config/include"]}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadCfgOptions(&CfgCLIOptions{
		IncludePath:       []string{"cli/include"},
		ExplicitBoolFlags: map[string]bool{},
	})
	suite.Require().NoError(err)

	cliIdx, cfgIdx := -1, -1
	for i, p := range opts.IncludePath {
		if strings.Contains(p, "cli") {
			cliIdx = i
		}

		if strings.Contains(p, "config") {
			cfgIdx = i
		}
	}

	suite.NotEqual(-1, cliIdx, "CLI include path should be present")
	suite.NotEqual(-1, cfgIdx, "config include path should be present")
	suite.Less(cliIdx, cfgIdx, "CLI include path must appear before config path")
}

// TestLoadCfgOptions_CLIModulePathPrependedBeforeConfig verifies that
// CLI module paths appear before config-file paths in the merged result.
func (suite *OptionsTestSuite) TestLoadCfgOptions_CLIModulePathPrependedBeforeConfig() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_cfg_mod"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_cfg_mod"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"compiler":{"modulePath":["config/module"]}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadCfgOptions(&CfgCLIOptions{
		ModulePath:        []string{"cli/module"},
		ExplicitBoolFlags: map[string]bool{},
	})
	suite.Require().NoError(err)

	cliIdx, cfgIdx := -1, -1
	for i, p := range opts.ModulePath {
		if strings.Contains(p, "cli") {
			cliIdx = i
		}

		if strings.Contains(p, "config") {
			cfgIdx = i
		}
	}

	suite.NotEqual(-1, cliIdx, "CLI module path should be present")
	suite.NotEqual(-1, cfgIdx, "config module path should be present")
	suite.Less(cliIdx, cfgIdx, "CLI module path must appear before config path")
}

// TestLoadCfgOptions_CLILibraryPathPrependedBeforeConfig verifies that
// CLI library paths appear before config-file paths in the merged result.
func (suite *OptionsTestSuite) TestLoadCfgOptions_CLILibraryPathPrependedBeforeConfig() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_cfg_lib"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_cfg_lib"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"compiler":{"libraryPath":["config/lib"]}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadCfgOptions(&CfgCLIOptions{
		LibraryPath:       []string{"cli/lib"},
		ExplicitBoolFlags: map[string]bool{},
	})
	suite.Require().NoError(err)

	cliIdx, cfgIdx := -1, -1
	for i, p := range opts.LibraryPath {
		if strings.Contains(p, "cli") {
			cliIdx = i
		}

		if strings.Contains(p, "config") {
			cfgIdx = i
		}
	}

	suite.NotEqual(-1, cliIdx, "CLI library path should be present")
	suite.NotEqual(-1, cfgIdx, "config library path should be present")
	suite.Less(cliIdx, cfgIdx, "CLI library path must appear before config path")
}

// ---------------------------------------------------------------------------
// GlobalConfigDir / GlobalConfigPath
// ---------------------------------------------------------------------------

// TestGlobalConfigDir_UsesEnvVar verifies that $GENLINX_CONFIG_DIR is
// returned as-is when set.
func TestGlobalConfigDir_UsesEnvVar(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GENLINX_CONFIG_DIR", dir)

	got, err := GlobalConfigDir()
	assert.NoError(t, err)
	assert.Equal(t, dir, got)
}

// TestGlobalConfigDir_ReturnsNonEmpty verifies that, without the env var, a
// non-empty platform-appropriate directory is returned.
func TestGlobalConfigDir_ReturnsNonEmpty(t *testing.T) {
	t.Setenv("GENLINX_CONFIG_DIR", "")

	got, err := GlobalConfigDir()
	assert.NoError(t, err)
	assert.NotEmpty(t, got)
}

// TestGlobalConfigPath_EndsInConfigJSON verifies that the returned path
// always ends in "config.json".
func TestGlobalConfigPath_EndsInConfigJSON(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GENLINX_CONFIG_DIR", dir)

	got, err := GlobalConfigPath()
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(got, "config.json"),
		"GlobalConfigPath must end in config.json, got: %s", got)
}

// ---------------------------------------------------------------------------
// LoadArchiveOptions – ExtraGlobPatterns
// ---------------------------------------------------------------------------

// TestLoadArchiveOptions_ExtraGlobPatterns_NilCLI_UsesConfig verifies that
// ExtraGlobPatterns declared in the config file are loaded when no CLI opts
// are provided.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_ExtraGlobPatterns_NilCLI_UsesConfig() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_agp_nil"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_agp_nil"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"archive":{"extraGlobPatterns":["dist/*.jar","icons/**"]}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadArchiveOptions(nil)
	suite.Require().NoError(err)
	suite.Contains(opts.ExtraGlobPatterns, "dist/*.jar")
	suite.Contains(opts.ExtraGlobPatterns, "icons/**")
}

// TestLoadArchiveOptions_ExtraGlobPatterns_CLIOnlyWithNoCfg verifies that
// patterns supplied on the CLI are stored when no config file is present.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_ExtraGlobPatterns_CLIOnlyWithNoCfg() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_agp_cli"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_agp_cli"), 0o755))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	cliOpts := &ArchiveCLIOptions{
		ExtraGlobPatterns: []string{"plugins/**"},
		ExplicitBoolFlags: map[string]bool{},
	}

	opts, _, err := LoadArchiveOptions(cliOpts)
	suite.Require().NoError(err)
	suite.Contains(opts.ExtraGlobPatterns, "plugins/**")
}

// TestLoadArchiveOptions_ExtraGlobPatterns_CLIPrependedBeforeConfig verifies
// that CLI-supplied patterns are prepended before config-file patterns in the
// merged result, matching the precedence behaviour of ExtraFileSearchLocations.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_ExtraGlobPatterns_CLIPrependedBeforeConfig() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_agp_prepend"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_agp_prepend"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"archive":{"extraGlobPatterns":["config/*.cfg"]}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	cliOpts := &ArchiveCLIOptions{
		ExtraGlobPatterns: []string{"cli/**"},
		ExplicitBoolFlags: map[string]bool{},
	}

	opts, _, err := LoadArchiveOptions(cliOpts)
	suite.Require().NoError(err)

	cliIdx, cfgIdx := -1, -1
	for i, p := range opts.ExtraGlobPatterns {
		if strings.Contains(p, "cli") {
			cliIdx = i
		}

		if strings.Contains(p, "config") {
			cfgIdx = i
		}
	}

	suite.NotEqual(-1, cliIdx, "CLI glob pattern should be present in merged result")
	suite.NotEqual(-1, cfgIdx, "config-file glob pattern should be present in merged result")
	suite.Less(cliIdx, cfgIdx, "CLI glob pattern must appear before config-file pattern")
}

// TestLoadArchiveOptions_ExtraGlobPatterns_Deduplicated verifies that identical
// patterns appearing in both the CLI and the config file are deduplicated.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_ExtraGlobPatterns_Deduplicated() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_agp_dedup"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_agp_dedup"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"archive":{"extraGlobPatterns":["shared/**"]}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	cliOpts := &ArchiveCLIOptions{
		ExtraGlobPatterns: []string{"shared/**"}, // same as config
		ExplicitBoolFlags: map[string]bool{},
	}

	opts, _, err := LoadArchiveOptions(cliOpts)
	suite.Require().NoError(err)

	count := 0
	for _, p := range opts.ExtraGlobPatterns {
		if p == "shared/**" {
			count++
		}
	}

	suite.Equal(1, count, "duplicate pattern shared/** must appear only once after merge")
}

// TestLoadArchiveOptions_ExtraGlobPatterns_MultipleCLIPatterns verifies that
// multiple patterns passed on the CLI are all stored.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_ExtraGlobPatterns_MultipleCLIPatterns() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_agp_multi"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_agp_multi"), 0o755))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	cliOpts := &ArchiveCLIOptions{
		ExtraGlobPatterns: []string{"dist/*.jar", "**/*.json", "docs"},
		ExplicitBoolFlags: map[string]bool{},
	}

	opts, _, err := LoadArchiveOptions(cliOpts)
	suite.Require().NoError(err)
	suite.Contains(opts.ExtraGlobPatterns, "dist/*.jar")
	suite.Contains(opts.ExtraGlobPatterns, "**/*.json")
	suite.Contains(opts.ExtraGlobPatterns, "docs")
}

// TestLoadArchiveOptions_ExtraGlobPatterns_EmptyCLI_NoOverride verifies that
// an empty CLI ExtraGlobPatterns slice does not wipe out config-file patterns.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_ExtraGlobPatterns_EmptyCLI_NoOverride() {
	suite.T().Setenv("GENLINX_CONFIG_DIR", filepath.Join(suite.tempDir, "no_global_agp_empty"))
	suite.Require().NoError(os.MkdirAll(filepath.Join(suite.tempDir, "no_global_agp_empty"), 0o755))

	suite.Require().NoError(os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(`{"archive":{"extraGlobPatterns":["assets/**"]}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	// CLI provides no patterns (empty slice — equivalent to flag not set).
	cliOpts := &ArchiveCLIOptions{
		ExtraGlobPatterns: []string{},
		ExplicitBoolFlags: map[string]bool{},
	}

	opts, _, err := LoadArchiveOptions(cliOpts)
	suite.Require().NoError(err)
	suite.Contains(opts.ExtraGlobPatterns, "assets/**",
		"config-file pattern must be preserved when CLI provides no patterns")
}

// TestLoadArchiveOptions_ExtraGlobPatterns_GlobalConfigLoaded verifies that
// ExtraGlobPatterns declared in the global config are surfaced in the merged
// result when no local config is present.
func (suite *OptionsTestSuite) TestLoadArchiveOptions_ExtraGlobPatterns_GlobalConfigLoaded() {
	globalDir := filepath.Join(suite.tempDir, "global_agp")
	suite.Require().NoError(os.MkdirAll(globalDir, 0o755))
	suite.T().Setenv("GENLINX_CONFIG_DIR", globalDir)

	suite.Require().NoError(os.WriteFile(
		filepath.Join(globalDir, "config.json"),
		[]byte(`{"archive":{"extraGlobPatterns":["global/**"]}}`),
		0o644,
	))

	oldWd, err := os.Getwd()
	suite.Require().NoError(err)
	defer os.Chdir(oldWd) //nolint:errcheck
	suite.Require().NoError(os.Chdir(suite.tempDir))

	opts, _, err := LoadArchiveOptions(nil)
	suite.Require().NoError(err)
	suite.Contains(opts.ExtraGlobPatterns, "global/**",
		"ExtraGlobPatterns from global config must be present in merged result")
}

// TestGlobalConfigPath_UnderConfigDir verifies that the config.json resides
// inside GlobalConfigDir.
func TestGlobalConfigPath_UnderConfigDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GENLINX_CONFIG_DIR", dir)

	wantDir, err := GlobalConfigDir()
	assert.NoError(t, err)

	gotPath, err := GlobalConfigPath()
	assert.NoError(t, err)

	assert.Equal(t, filepath.Join(wantDir, "config.json"), gotPath)
}

// TestOptionsTestSuite runs the test suite
func TestOptionsTestSuite(t *testing.T) {
	suite.Run(t, new(OptionsTestSuite))
}
