package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Norgate-AV/genlinx-go/internal/utils"
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	CFG     CFGConfig     `mapstructure:"cfg"`
	Archive ArchiveConfig `mapstructure:"archive"`
	Build   BuildConfig   `mapstructure:"build"`
}

// CFGConfig represents CFG command configuration
type CFGConfig struct {
	OutputFile                string   `mapstructure:"outputFile"`
	OutputLogFile             string   `mapstructure:"outputLogFile"`
	OutputLogFileOption       string   `mapstructure:"outputLogFileOption"`
	OutputLogConsoleOption    bool     `mapstructure:"outputLogConsoleOption"`
	BuildWithDebugInformation bool     `mapstructure:"buildWithDebugInformation"`
	BuildWithSource           bool     `mapstructure:"buildWithSource"`
	IncludePath               []string `mapstructure:"includePath"`
	ModulePath                []string `mapstructure:"modulePath"`
	LibraryPath               []string `mapstructure:"libraryPath"`
	All                       bool     `mapstructure:"all"`
}

// ArchiveConfig represents archive command configuration
type ArchiveConfig struct {
	OutputFile                 string   `mapstructure:"outputFile"`
	IncludeCompiledSourceFiles bool     `mapstructure:"includeCompiledSourceFiles"`
	IncludeCompiledModuleFiles bool     `mapstructure:"includeCompiledModuleFiles"`
	IncludeFilesNotInWorkspace bool     `mapstructure:"includeFilesNotInWorkspace"`
	ExtraFileSearchLocations   []string `mapstructure:"extraFileSearchLocations"`
	ExtraFileArchiveLocation   string   `mapstructure:"extraFileArchiveLocation"`
	All                        bool     `mapstructure:"all"`
	IgnoredFiles               []string `mapstructure:"ignoredFiles"`
}

// BuildConfig represents build command configuration
type BuildConfig struct {
	NLRC  NLRCConfig  `mapstructure:"nlrc"`
	Shell ShellConfig `mapstructure:"shell"`
	All   bool        `mapstructure:"all"`
}

// NLRCConfig represents NetLinx compiler configuration
type NLRCConfig struct {
	Path        string   `mapstructure:"path"`
	IncludePath []string `mapstructure:"includePath"`
	ModulePath  []string `mapstructure:"modulePath"`
	LibraryPath []string `mapstructure:"libraryPath"`
}

// ShellConfig represents shell configuration
type ShellConfig struct {
	Path string `mapstructure:"path"`
}

var defaultConfig = Config{
	CFG: CFGConfig{
		OutputFile:                "build.cfg",
		OutputLogFile:             "build.log",
		OutputLogFileOption:       "N",
		OutputLogConsoleOption:    true,
		BuildWithDebugInformation: false,
		BuildWithSource:           false,
		IncludePath:               utils.NormalizePaths([]string{"C:/Program Files (x86)/Common Files/AMXShare/AXIs"}),
		ModulePath: utils.NormalizePaths([]string{
			"C:/Program Files (x86)/Common Files/AMXShare/Duet/bundle",
			"C:/Program Files (x86)/Common Files/AMXShare/Duet/lib",
			"C:/Program Files (x86)/Common Files/AMXShare/Duet/module",
		}),
		LibraryPath: utils.NormalizePaths([]string{"C:/Program Files (x86)/Common Files/AMXShare/SYCs"}),
		All:         false,
	},
	Archive: ArchiveConfig{
		OutputFile:                 "archive.zip",
		IncludeCompiledSourceFiles: true,
		IncludeCompiledModuleFiles: true,
		IncludeFilesNotInWorkspace: true,
		ExtraFileSearchLocations:   utils.NormalizePaths([]string{"C:/Program Files (x86)/Common Files/AMXShare"}),
		ExtraFileArchiveLocation:   utils.NormalizePath(".genlinx"),
		All:                        false,
		IgnoredFiles: []string{
			"G4API.axi",
			"NetLinx.axi",
			"SNAPI.axi",
			"UnicodeLib.axi",
			"componentssdk.jar",
			"componentssdkrt.jar",
			"DeviceDriverEngine.jar",
			"devicesdkrt.jar",
			"jregex1.2_01-bundle.jar",
			"js-14-bundle.jar",
			"json-bundle.jar",
			"picocontainer-1.3-bundle.jar",
			"snapirouter.jar",
			"snapirouter2.jar",
		},
	},
	Build: BuildConfig{
		NLRC: NLRCConfig{
			Path: utils.NormalizePath("C:/Program Files (x86)/Common Files/AMXShare/COM/NLRC.exe"),
			IncludePath: utils.NormalizePaths([]string{
				"C:/Program Files (x86)/Common Files/AMXShare/AXIs",
			}),
			ModulePath: utils.NormalizePaths([]string{
				"C:/Program Files (x86)/Common Files/AMXShare/Duet/bundle",
				"C:/Program Files (x86)/Common Files/AMXShare/Duet/lib",
				"C:/Program Files (x86)/Common Files/AMXShare/Duet/module",
			}),
			LibraryPath: utils.NormalizePaths([]string{"C:/Program Files (x86)/Common Files/AMXShare/SYCs"}),
		},
		Shell: ShellConfig{
			Path: utils.NormalizePath("C:/Windows/System32/cmd.exe"),
		},
		All: false,
	},
}

// LoadConfig loads the configuration from various sources
func LoadConfig() (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("cfg", defaultConfig.CFG)
	v.SetDefault("archive", defaultConfig.Archive)
	v.SetDefault("build", defaultConfig.Build)

	// Configuration file search paths
	configPaths := []string{
		".", // Current directory
		filepath.Join(os.Getenv("HOME"), ".config", "genlinx"), // User config
		"/etc/genlinx", // System config
	}

	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	v.SetConfigName("genlinx")
	v.SetConfigType("json")

	// Environment variables
	v.SetEnvPrefix("GENLINX")
	v.AutomaticEnv()

	// Read configuration
	if err := v.ReadInConfig(); err != nil {
		// If config file doesn't exist, use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Normalize all paths in the configuration
	normalizeConfigPaths(&config)

	return &config, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *Config, filename string) error {
	v := viper.New()
	v.Set("cfg", config.CFG)
	v.Set("archive", config.Archive)
	v.Set("build", config.Build)

	return v.WriteConfigAs(filename)
}

// normalizeConfigPaths normalizes all paths in the configuration to OS-specific format
func normalizeConfigPaths(config *Config) {
	// Normalize CFG paths
	config.CFG.IncludePath = utils.NormalizePaths(config.CFG.IncludePath)
	config.CFG.ModulePath = utils.NormalizePaths(config.CFG.ModulePath)
	config.CFG.LibraryPath = utils.NormalizePaths(config.CFG.LibraryPath)

	// Normalize Archive paths
	config.Archive.ExtraFileSearchLocations = utils.NormalizePaths(config.Archive.ExtraFileSearchLocations)
	config.Archive.ExtraFileArchiveLocation = utils.NormalizePath(config.Archive.ExtraFileArchiveLocation)

	// Normalize Build paths
	config.Build.NLRC.Path = utils.NormalizePath(config.Build.NLRC.Path)
	config.Build.NLRC.IncludePath = utils.NormalizePaths(config.Build.NLRC.IncludePath)
	config.Build.NLRC.ModulePath = utils.NormalizePaths(config.Build.NLRC.ModulePath)
	config.Build.NLRC.LibraryPath = utils.NormalizePaths(config.Build.NLRC.LibraryPath)
	config.Build.Shell.Path = utils.NormalizePath(config.Build.Shell.Path)
}

// LoadDefaultConfig returns the default configuration
func LoadDefaultConfig() (*Config, error) {
	config := defaultConfig

	// Normalize all paths in the default configuration
	normalizeConfigPaths(&config)

	return &config, nil
}
