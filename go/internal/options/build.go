package options

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Norgate-AV/genlinx-go/internal/config"
	"github.com/Norgate-AV/genlinx-go/internal/utils"
	"github.com/spf13/viper"
)

// BuildOptions represents the merged options for the build command
type BuildOptions struct {
	// Source files to compile
	SourceFiles []string

	// CFG files to compile
	CFGFiles []string

	// Include paths for the compiler
	IncludePath []string

	// Module paths for the compiler
	ModulePath []string

	// Library paths for the compiler
	LibraryPath []string

	// Output path for compiled files
	OutputPath string

	// Whether to build all files
	All bool

	// Whether to enable verbose output
	Verbose bool

	// NLRC compiler path
	NLRCPath string

	// Shell path
	ShellPath string
}

// MergeOptions merges configuration from multiple sources with proper precedence
type MergeOptions struct {
	// CLI options (highest precedence)
	CLI *CLIOptions

	// Local config file path (optional)
	LocalConfigPath string

	// Global config file path (optional)
	GlobalConfigPath string
}

// CLIOptions represents command-line options
type CLIOptions struct {
	SourceFiles []string
	CFGFiles    []string
	IncludePath []string
	ModulePath  []string
	LibraryPath []string
	OutputPath  string
	All         bool
	Verbose     bool
}

// ConfigLoadResult wraps a config with loading information
type ConfigLoadResult struct {
	Config *config.Config
	Path   string
	Found  bool
}

// ConfigLoadInfo holds information about configuration loading for verbose output
type ConfigLoadInfo struct {
	DefaultLoaded bool
	GlobalResult  ConfigLoadResult
	LocalResult   ConfigLoadResult
}

// LoadBuildOptions loads and merges build options from all sources
func LoadBuildOptions(cliOpts *CLIOptions) (*BuildOptions, *ConfigLoadInfo, error) {
	// Load default configuration
	defaultCfg, err := config.LoadDefaultConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load default config: %w", err)
	}

	// Load global configuration
	globalResult, err := loadGlobalConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load global config: %w", err)
	}

	// Load local configuration
	localResult, err := loadLocalConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load local config: %w", err)
	}

	// Create config load info
	configInfo := &ConfigLoadInfo{
		DefaultLoaded: true,
		GlobalResult:  globalResult,
		LocalResult:   localResult,
	}

	// Merge configurations in order of precedence: default -> global -> local -> CLI
	mergedCfg := mergeConfigs(defaultCfg, globalResult.Config, localResult.Config)

	// Create build options from merged config
	opts := &BuildOptions{
		IncludePath: mergedCfg.Build.NLRC.IncludePath,
		ModulePath:  mergedCfg.Build.NLRC.ModulePath,
		LibraryPath: mergedCfg.Build.NLRC.LibraryPath,
		NLRCPath:    mergedCfg.Build.NLRC.Path,
		ShellPath:   mergedCfg.Build.Shell.Path,
		All:         mergedCfg.Build.All,
	}

	// Apply CLI options (highest precedence)
	if cliOpts != nil {
		if len(cliOpts.SourceFiles) > 0 {
			opts.SourceFiles = cliOpts.SourceFiles
		}
		if len(cliOpts.CFGFiles) > 0 {
			opts.CFGFiles = cliOpts.CFGFiles
		}
		if len(cliOpts.IncludePath) > 0 {
			opts.IncludePath = prependAndDeduplicate(opts.IncludePath, cliOpts.IncludePath)
		}
		if len(cliOpts.ModulePath) > 0 {
			opts.ModulePath = prependAndDeduplicate(opts.ModulePath, cliOpts.ModulePath)
		}
		if len(cliOpts.LibraryPath) > 0 {
			opts.LibraryPath = prependAndDeduplicate(opts.LibraryPath, cliOpts.LibraryPath)
		}
		if cliOpts.OutputPath != "" {
			opts.OutputPath = cliOpts.OutputPath
		}
		if cliOpts.All {
			opts.All = cliOpts.All
		}
		if cliOpts.Verbose {
			opts.Verbose = cliOpts.Verbose
		}
	}

	// Resolve all paths to absolute paths
	if err := opts.resolvePaths(); err != nil {
		return nil, nil, fmt.Errorf("failed to resolve paths: %w", err)
	}

	return opts, configInfo, nil
}

// loadGlobalConfig loads the global configuration file
func loadGlobalConfig() (ConfigLoadResult, error) {
	v := viper.New()

	// Determine global config directory based on OS
	var globalConfigDir string
	if utils.IsWindows() {
		// Windows: %APPDATA%\genlinx
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return ConfigLoadResult{Config: &config.Config{}, Path: "", Found: false}, fmt.Errorf("APPDATA environment variable not set")
		}
		globalConfigDir = filepath.Join(appData, "genlinx")
	} else {
		// Unix-like: ~/.config/genlinx
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ConfigLoadResult{Config: &config.Config{}, Path: "", Found: false}, fmt.Errorf("failed to get user home directory: %w", err)
		}
		globalConfigDir = filepath.Join(homeDir, ".config", "genlinx")
	}

	// Try different config file names and formats
	configFiles := []string{
		filepath.Join(globalConfigDir, "genlinx.json"),
		filepath.Join(globalConfigDir, "genlinx.yaml"),
		filepath.Join(globalConfigDir, "genlinx.yml"),
	}

	var cfg config.Config
	configFound := false
	var configPath string

	for _, configFile := range configFiles {
		if utils.FileExists(configFile) {
			v.SetConfigFile(configFile)
			if err := v.ReadInConfig(); err != nil {
				continue // Try next file
			}
			if err := v.Unmarshal(&cfg); err != nil {
				continue // Try next file
			}
			// Successfully loaded config
			configFound = true
			configPath = configFile
			break
		}
	}

	return ConfigLoadResult{
		Config: &cfg,
		Path:   configPath,
		Found:  configFound,
	}, nil
}

// loadLocalConfig loads the local configuration file using find-up approach
func loadLocalConfig() (ConfigLoadResult, error) {
	v := viper.New()

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return ConfigLoadResult{Config: &config.Config{}, Path: "", Found: false}, fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Possible config file names to look for
	configFileNames := []string{
		".genlinxrc.json",
		".genlinxrc.yaml",
		".genlinxrc.yml",
		".genlinx.json",
		".genlinx.yaml",
		".genlinx.yml",
	}

	// Walk up the directory tree from current directory to root
	currentDir := cwd
	for {
		// Check each possible config file name in the current directory
		for _, configFileName := range configFileNames {
			configFilePath := filepath.Join(currentDir, configFileName)
			if utils.FileExists(configFilePath) {
				v.SetConfigFile(configFilePath)
				if err := v.ReadInConfig(); err != nil {
					continue // Try next file
				}
				var cfg config.Config
				if err := v.Unmarshal(&cfg); err != nil {
					continue // Try next file
				}
				// Successfully loaded config
				return ConfigLoadResult{
					Config: &cfg,
					Path:   configFilePath,
					Found:  true,
				}, nil
			}
		}

		// Move up to parent directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached root directory, stop searching
			break
		}
		currentDir = parentDir
	}

	// No config file found, return empty config
	return ConfigLoadResult{
		Config: &config.Config{},
		Path:   "",
		Found:  false,
	}, nil
}

// mergeConfigs merges multiple configurations with precedence (later configs override earlier ones)
func mergeConfigs(configs ...*config.Config) *config.Config {
	if len(configs) == 0 {
		return &config.Config{}
	}

	result := configs[0]

	// Merge each subsequent config
	for i := 1; i < len(configs); i++ {
		cfg := configs[i]
		if cfg == nil {
			continue
		}

		// Merge Build config
		if cfg.Build.NLRC.Path != "" {
			result.Build.NLRC.Path = cfg.Build.NLRC.Path
		}
		if len(cfg.Build.NLRC.IncludePath) > 0 {
			result.Build.NLRC.IncludePath = append(result.Build.NLRC.IncludePath, cfg.Build.NLRC.IncludePath...)
		}
		if len(cfg.Build.NLRC.ModulePath) > 0 {
			result.Build.NLRC.ModulePath = append(result.Build.NLRC.ModulePath, cfg.Build.NLRC.ModulePath...)
		}
		if len(cfg.Build.NLRC.LibraryPath) > 0 {
			result.Build.NLRC.LibraryPath = append(result.Build.NLRC.LibraryPath, cfg.Build.NLRC.LibraryPath...)
		}
		if cfg.Build.Shell.Path != "" {
			result.Build.Shell.Path = cfg.Build.Shell.Path
		}
		if cfg.Build.All {
			result.Build.All = cfg.Build.All
		}

		// Merge CFG config
		if cfg.CFG.OutputFile != "" {
			result.CFG.OutputFile = cfg.CFG.OutputFile
		}
		if len(cfg.CFG.IncludePath) > 0 {
			result.CFG.IncludePath = append(result.CFG.IncludePath, cfg.CFG.IncludePath...)
		}
		if len(cfg.CFG.ModulePath) > 0 {
			result.CFG.ModulePath = append(result.CFG.ModulePath, cfg.CFG.ModulePath...)
		}
		if len(cfg.CFG.LibraryPath) > 0 {
			result.CFG.LibraryPath = append(result.CFG.LibraryPath, cfg.CFG.LibraryPath...)
		}

		// Merge Archive config
		if cfg.Archive.OutputFile != "" {
			result.Archive.OutputFile = cfg.Archive.OutputFile
		}
		if len(cfg.Archive.ExtraFileSearchLocations) > 0 {
			result.Archive.ExtraFileSearchLocations = append(result.Archive.ExtraFileSearchLocations, cfg.Archive.ExtraFileSearchLocations...)
		}
		if cfg.Archive.ExtraFileArchiveLocation != "" {
			result.Archive.ExtraFileArchiveLocation = cfg.Archive.ExtraFileArchiveLocation
		}
	}

	return result
}

// prependAndDeduplicate prepends new items to the existing slice and removes duplicates
func prependAndDeduplicate(existing, new []string) []string {
	// Create a map to track seen items
	seen := make(map[string]bool)

	// Add new items first (prepended)
	var result []string
	for _, item := range new {
		if !seen[item] {
			result = append(result, item)
			seen[item] = true
		}
	}

	// Add existing items that haven't been seen
	for _, item := range existing {
		if !seen[item] {
			result = append(result, item)
			seen[item] = true
		}
	}

	return result
}

// resolvePaths resolves all paths in the options to absolute paths
func (o *BuildOptions) resolvePaths() error {
	var err error

	// Resolve source files
	for i, file := range o.SourceFiles {
		if o.SourceFiles[i], err = filepath.Abs(file); err != nil {
			return fmt.Errorf("failed to resolve source file path %s: %w", file, err)
		}
	}

	// Resolve CFG files
	for i, file := range o.CFGFiles {
		if o.CFGFiles[i], err = filepath.Abs(file); err != nil {
			return fmt.Errorf("failed to resolve CFG file path %s: %w", file, err)
		}
	}

	// Resolve include paths
	for i, path := range o.IncludePath {
		if o.IncludePath[i], err = filepath.Abs(path); err != nil {
			return fmt.Errorf("failed to resolve include path %s: %w", path, err)
		}
	}

	// Resolve module paths
	for i, path := range o.ModulePath {
		if o.ModulePath[i], err = filepath.Abs(path); err != nil {
			return fmt.Errorf("failed to resolve module path %s: %w", path, err)
		}
	}

	// Resolve library paths
	for i, path := range o.LibraryPath {
		if o.LibraryPath[i], err = filepath.Abs(path); err != nil {
			return fmt.Errorf("failed to resolve library path %s: %w", path, err)
		}
	}

	// Resolve output path
	if o.OutputPath != "" {
		if o.OutputPath, err = filepath.Abs(o.OutputPath); err != nil {
			return fmt.Errorf("failed to resolve output path %s: %w", o.OutputPath, err)
		}
	}

	// Resolve NLRC path
	if o.NLRCPath != "" {
		if o.NLRCPath, err = filepath.Abs(o.NLRCPath); err != nil {
			return fmt.Errorf("failed to resolve NLRC path %s: %w", o.NLRCPath, err)
		}
	}

	// Resolve shell path
	if o.ShellPath != "" {
		if o.ShellPath, err = filepath.Abs(o.ShellPath); err != nil {
			return fmt.Errorf("failed to resolve shell path %s: %w", o.ShellPath, err)
		}
	}

	return nil
}
