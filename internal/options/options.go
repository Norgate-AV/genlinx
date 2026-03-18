package options

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"

	"github.com/Norgate-AV/genlinx/internal/config"
	"github.com/Norgate-AV/genlinx/internal/utils"
)

// ConfigLoadResult wraps a config with loading information
type ConfigLoadResult struct {
	Config   *config.Config
	Path     string
	Found    bool
	Presence boolPresence
}

// boolPresence records which boolean config keys were explicitly set in a
// loaded file. JSON/YAML unmarshal cannot distinguish an absent field from
// one explicitly set to its zero value, so Viper's IsSet is used when
// loading to populate this struct.
type boolPresence struct {
	cfgOutputLogConsoleOption    bool
	cfgBuildWithDebugInformation bool
	cfgBuildWithSource           bool
	archiveIncludeCompiledSrc    bool
	archiveIncludeCompiledMod    bool
	archiveIncludeFilesNotInWs   bool
}

// configMergeInput packages a config with its boolean-presence information
// for use by mergeConfigsWithPresence.
type configMergeInput struct {
	cfg      *config.Config
	presence boolPresence
}

// viperBoolPresence probes which boolean config keys were explicitly present
// in the Viper instance after reading a config file.
func viperBoolPresence(v *viper.Viper) boolPresence {
	return boolPresence{
		cfgOutputLogConsoleOption:    v.IsSet("cfg.outputLogConsoleOption"),
		cfgBuildWithDebugInformation: v.IsSet("cfg.buildWithDebugInformation"),
		cfgBuildWithSource:           v.IsSet("cfg.buildWithSource"),
		archiveIncludeCompiledSrc:    v.IsSet("archive.includeCompiledSourceFiles"),
		archiveIncludeCompiledMod:    v.IsSet("archive.includeCompiledModuleFiles"),
		archiveIncludeFilesNotInWs:   v.IsSet("archive.includeFilesNotInWorkspace"),
	}
}

// ConfigLoadInfo holds information about configuration loading for verbose output
type ConfigLoadInfo struct {
	DefaultLoaded bool
	GlobalResult  ConfigLoadResult
	LocalResult   ConfigLoadResult
}

// Print writes a human-readable summary of which config files were loaded.
func (info *ConfigLoadInfo) Print() {
	fmt.Println("Configuration loading:")
	fmt.Println("  Default config: Built-in defaults loaded")

	if info.GlobalResult.Found {
		fmt.Printf("  Global config: Loaded from %s\n", info.GlobalResult.Path)
	}

	if info.LocalResult.Found {
		fmt.Printf("  Local config: Loaded from %s\n", info.LocalResult.Path)
	}
}

// LoadMergedConfig is the central config-merge entry point. It loads the default,
// global, and local configs, merges them in order of increasing precedence
// (default < global < local) and returns the fully-merged config together
// with load metadata.  Every command-specific Load*Options function calls
// this rather than duplicating the load/merge boilerplate.
func LoadMergedConfig() (*config.Config, *ConfigLoadInfo, error) {
	defaultCfg, err := config.LoadDefaultConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load default config: %w", err)
	}

	globalResult, err := LoadGlobalConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load global config: %w", err)
	}

	localResult, err := LoadLocalConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load local config: %w", err)
	}

	configInfo := &ConfigLoadInfo{
		DefaultLoaded: true,
		GlobalResult:  globalResult,
		LocalResult:   localResult,
	}

	mergedCfg := mergeConfigsWithPresence(defaultCfg,
		configMergeInput{globalResult.Config, globalResult.Presence},
		configMergeInput{localResult.Config, localResult.Presence},
	)

	if err := resolveConfigPaths(mergedCfg); err != nil {
		return nil, nil, fmt.Errorf("failed to resolve config paths: %w", err)
	}

	return mergedCfg, configInfo, nil
}

// globalConfigDirs returns an ordered list of candidate directories to search
// for the global config file. The first existing config file found wins.
//
//   - If $GENLINX_CONFIG_DIR is set it is the only candidate.
//   - On Windows: %APPDATA%\genlinx is checked first, then ~/.config/genlinx.
//   - On Linux/macOS: only ~/.config/genlinx is checked.
func globalConfigDirs() ([]string, error) {
	if dir := os.Getenv("GENLINX_CONFIG_DIR"); dir != "" {
		return []string{dir}, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	dotConfig := filepath.Join(homeDir, ".config", "genlinx")

	if runtime.GOOS == "windows" {
		var dirs []string

		if appData := os.Getenv("APPDATA"); appData != "" {
			dirs = append(dirs, filepath.Join(appData, "genlinx"))
		}

		dirs = append(dirs, dotConfig)
		return dirs, nil
	}

	return []string{dotConfig}, nil
}

// GlobalConfigDir returns the primary (canonical) directory used when creating
// or editing the global config file. On Windows this is %APPDATA%\genlinx;
// elsewhere it is ~/.config/genlinx. $GENLINX_CONFIG_DIR overrides both.
func GlobalConfigDir() (string, error) {
	dirs, err := globalConfigDirs()
	if err != nil {
		return "", err
	}

	return dirs[0], nil
}

// GlobalConfigPath returns the canonical path for the global JSON config file
// (used for editing / creating a new global config).
func GlobalConfigPath() (string, error) {
	dir, err := GlobalConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "config.json"), nil
}

// LoadGlobalConfig loads the global configuration file.
// It searches each candidate directory (from globalConfigDirs) for config.json,
// config.yaml, or config.yml and returns the first one successfully loaded.
func LoadGlobalConfig() (ConfigLoadResult, error) {
	v := viper.New()

	dirs, err := globalConfigDirs()
	if err != nil {
		return ConfigLoadResult{Config: &config.Config{}, Path: "", Found: false}, err
	}

	fileNames := []string{"config.json", "config.yaml", "config.yml"}

	var cfg config.Config

	for _, dir := range dirs {
		for _, name := range fileNames {
			configFile := filepath.Join(dir, name)
			if !utils.FileExists(configFile) {
				continue
			}

			v.SetConfigFile(configFile)
			if err := v.ReadInConfig(); err != nil {
				continue
			}
			if err := v.Unmarshal(&cfg); err != nil {
				continue
			}

			config.NormalizeConfigPaths(&cfg)

			return ConfigLoadResult{
				Config:   &cfg,
				Path:     configFile,
				Found:    true,
				Presence: viperBoolPresence(v),
			}, nil
		}
	}

	return ConfigLoadResult{Config: &config.Config{}, Path: "", Found: false}, nil
}

// LoadLocalConfig loads the local configuration file using find-up approach
func LoadLocalConfig() (ConfigLoadResult, error) {
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

				readErr := v.ReadInConfig()

				// An empty file isn't an error — treat it as an empty config.
				if readErr != nil {
					// Check if it's just an empty file.
					info, statErr := os.Stat(configFilePath)
					if statErr != nil || info.Size() > 0 {
						continue // Genuinely unreadable or malformed — skip.
					}
				}

				var cfg config.Config
				if err := v.Unmarshal(&cfg); err != nil {
					continue // Try next file
				}

				config.NormalizeConfigPaths(&cfg)
				// Successfully loaded config
				return ConfigLoadResult{
					Config:   &cfg,
					Path:     configFilePath,
					Found:    true,
					Presence: viperBoolPresence(v),
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

// resolveConfigPaths resolves all relative paths in cfg to absolute paths
// using the current working directory. This resolves paths at merge time
// so that config --list and the build commands always show/use fully-qualified paths.
func resolveConfigPaths(cfg *config.Config) error {
	abs := func(p string) (string, error) {
		if p == "" {
			return p, nil
		}

		return filepath.Abs(p)
	}

	absSlice := func(paths []string) error {
		for i, p := range paths {
			resolved, err := abs(p)
			if err != nil {
				return err
			}

			paths[i] = resolved
		}

		return nil
	}

	var err error

	if cfg.Compiler.Path, err = abs(cfg.Compiler.Path); err != nil {
		return fmt.Errorf("compiler.path: %w", err)
	}

	if err = absSlice(cfg.Compiler.IncludePath); err != nil {
		return fmt.Errorf("compiler.includePath: %w", err)
	}

	if err = absSlice(cfg.Compiler.ModulePath); err != nil {
		return fmt.Errorf("compiler.modulePath: %w", err)
	}

	if err = absSlice(cfg.Compiler.LibraryPath); err != nil {
		return fmt.Errorf("compiler.libraryPath: %w", err)
	}

	if err = absSlice(cfg.Archive.ExtraFileSearchLocations); err != nil {
		return fmt.Errorf("archive.extraFileSearchLocations: %w", err)
	}

	return nil
}

// mergeConfigs is a backward-compatible wrapper around mergeConfigsWithPresence
// with empty presence info. Used by tests that only exercise path/string merging.
func mergeConfigs(base *config.Config, overrides ...*config.Config) *config.Config {
	inputs := make([]configMergeInput, len(overrides))
	for i, cfg := range overrides {
		inputs[i] = configMergeInput{cfg: cfg}
	}

	return mergeConfigsWithPresence(base, inputs...)
}

// mergeConfigsWithPresence merges base with each override in order of increasing
// precedence (later overrides win). String/slice fields are merged when non-empty.
// Boolean fields are only applied when the corresponding presence flag is set,
// preventing JSON zero-value ambiguity from silently discarding a deliberate false.
func mergeConfigsWithPresence(base *config.Config, overrides ...configMergeInput) *config.Config {
	if base == nil {
		return &config.Config{}
	}

	result := *base // value copy — do not mutate the caller's struct

	for _, ov := range overrides {
		cfg := ov.cfg
		p := ov.presence

		if cfg == nil {
			continue
		}

		// Merge compiler config
		if cfg.Compiler.Path != "" {
			result.Compiler.Path = cfg.Compiler.Path
		}

		if len(cfg.Compiler.IncludePath) > 0 {
			result.Compiler.IncludePath = prependAndDeduplicate(result.Compiler.IncludePath, cfg.Compiler.IncludePath)
		}

		if len(cfg.Compiler.ModulePath) > 0 {
			result.Compiler.ModulePath = prependAndDeduplicate(result.Compiler.ModulePath, cfg.Compiler.ModulePath)
		}

		if len(cfg.Compiler.LibraryPath) > 0 {
			result.Compiler.LibraryPath = prependAndDeduplicate(result.Compiler.LibraryPath, cfg.Compiler.LibraryPath)
		}

		// Merge CFG config
		if cfg.CFG.OutputFile != "" {
			result.CFG.OutputFile = cfg.CFG.OutputFile
		}

		if cfg.CFG.OutputLogFile != "" {
			result.CFG.OutputLogFile = cfg.CFG.OutputLogFile
		}

		if cfg.CFG.OutputLogFileOption != "" {
			result.CFG.OutputLogFileOption = cfg.CFG.OutputLogFileOption
		}

		if p.cfgOutputLogConsoleOption {
			result.CFG.OutputLogConsoleOption = cfg.CFG.OutputLogConsoleOption
		}

		if p.cfgBuildWithDebugInformation {
			result.CFG.BuildWithDebugInformation = cfg.CFG.BuildWithDebugInformation
		}

		if p.cfgBuildWithSource {
			result.CFG.BuildWithSource = cfg.CFG.BuildWithSource
		}

		// Merge Archive config
		if cfg.Archive.OutputFile != "" {
			result.Archive.OutputFile = cfg.Archive.OutputFile
		}

		if p.archiveIncludeCompiledSrc {
			result.Archive.IncludeCompiledSourceFiles = cfg.Archive.IncludeCompiledSourceFiles
		}

		if p.archiveIncludeCompiledMod {
			result.Archive.IncludeCompiledModuleFiles = cfg.Archive.IncludeCompiledModuleFiles
		}

		if p.archiveIncludeFilesNotInWs {
			result.Archive.IncludeFilesNotInWorkspace = cfg.Archive.IncludeFilesNotInWorkspace
		}

		if len(cfg.Archive.ExtraFileSearchLocations) > 0 {
			result.Archive.ExtraFileSearchLocations = prependAndDeduplicate(result.Archive.ExtraFileSearchLocations, cfg.Archive.ExtraFileSearchLocations)
		}

		// IgnoredFiles: nil-check distinguishes "absent" from "explicitly empty".
		if cfg.Archive.IgnoredFiles != nil {
			result.Archive.IgnoredFiles = cfg.Archive.IgnoredFiles
		}
	}

	return &result
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
