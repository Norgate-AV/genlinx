package config

import (
	"github.com/Norgate-AV/genlinx/internal/utils"
)

// Config represents the application configuration
type Config struct {
	CFG     CFGConfig     `mapstructure:"cfg"     json:"cfg"`
	Archive ArchiveConfig `mapstructure:"archive" json:"archive"`
	Build   BuildConfig   `mapstructure:"build"   json:"build"`
}

// CFGConfig represents CFG command configuration
type CFGConfig struct {
	OutputFile                string   `mapstructure:"outputFile"                json:"outputFile"`
	OutputLogFile             string   `mapstructure:"outputLogFile"              json:"outputLogFile"`
	OutputLogFileOption       string   `mapstructure:"outputLogFileOption"        json:"outputLogFileOption"`
	OutputLogConsoleOption    bool     `mapstructure:"outputLogConsoleOption"     json:"outputLogConsoleOption"`
	BuildWithDebugInformation bool     `mapstructure:"buildWithDebugInformation" json:"buildWithDebugInformation"`
	BuildWithSource           bool     `mapstructure:"buildWithSource"            json:"buildWithSource"`
	IncludePath               []string `mapstructure:"includePath"                json:"includePath"`
	ModulePath                []string `mapstructure:"modulePath"                 json:"modulePath"`
	LibraryPath               []string `mapstructure:"libraryPath"                json:"libraryPath"`
	All                       bool     `mapstructure:"all"                        json:"all"`
}

// ArchiveConfig represents archive command configuration
type ArchiveConfig struct {
	OutputFile                 string   `mapstructure:"outputFile"                 json:"outputFile"`
	IncludeCompiledSourceFiles bool     `mapstructure:"includeCompiledSourceFiles" json:"includeCompiledSourceFiles"`
	IncludeCompiledModuleFiles bool     `mapstructure:"includeCompiledModuleFiles" json:"includeCompiledModuleFiles"`
	IncludeFilesNotInWorkspace bool     `mapstructure:"includeFilesNotInWorkspace" json:"includeFilesNotInWorkspace"`
	ExtraFileSearchLocations   []string `mapstructure:"extraFileSearchLocations"   json:"extraFileSearchLocations"`
	All                        bool     `mapstructure:"all"                        json:"all"`
	IgnoredFiles               []string `mapstructure:"ignoredFiles"               json:"ignoredFiles"`
}

// BuildConfig represents build command configuration
type BuildConfig struct {
	NLRC NLRCConfig `mapstructure:"nlrc" json:"nlrc"`
	All  bool       `mapstructure:"all"  json:"all"`
}

// NLRCConfig represents NetLinx compiler configuration
type NLRCConfig struct {
	Path        string   `mapstructure:"path"        json:"path"`
	IncludePath []string `mapstructure:"includePath" json:"includePath"`
	ModulePath  []string `mapstructure:"modulePath"  json:"modulePath"`
	LibraryPath []string `mapstructure:"libraryPath" json:"libraryPath"`
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
		All: false,
	},
}

// NormalizeConfigPaths normalizes all paths in the configuration to OS-specific format
func NormalizeConfigPaths(config *Config) {
	// Normalize CFG paths
	config.CFG.IncludePath = utils.NormalizePaths(config.CFG.IncludePath)
	config.CFG.ModulePath = utils.NormalizePaths(config.CFG.ModulePath)
	config.CFG.LibraryPath = utils.NormalizePaths(config.CFG.LibraryPath)

	// Normalize Archive paths
	config.Archive.ExtraFileSearchLocations = utils.NormalizePaths(config.Archive.ExtraFileSearchLocations)

	// Normalize Build paths
	config.Build.NLRC.Path = utils.NormalizePath(config.Build.NLRC.Path)
	config.Build.NLRC.IncludePath = utils.NormalizePaths(config.Build.NLRC.IncludePath)
	config.Build.NLRC.ModulePath = utils.NormalizePaths(config.Build.NLRC.ModulePath)
	config.Build.NLRC.LibraryPath = utils.NormalizePaths(config.Build.NLRC.LibraryPath)
}

// LoadDefaultConfig returns the default configuration
func LoadDefaultConfig() (*Config, error) {
	config := defaultConfig

	// Normalize all paths in the default configuration
	NormalizeConfigPaths(&config)

	return &config, nil
}
