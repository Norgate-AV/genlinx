package config

import (
	"github.com/Norgate-AV/genlinx/internal/utils"
)

// Config represents the application configuration
type Config struct {
	Core    CoreConfig    `mapstructure:"core"    json:"core"`
	NLRC    NLRCConfig    `mapstructure:"nlrc"    json:"nlrc"`
	CFG     CFGConfig     `mapstructure:"cfg"     json:"cfg"`
	Archive ArchiveConfig `mapstructure:"archive" json:"archive"`
	Build   BuildConfig   `mapstructure:"build"   json:"build"`
}

// CoreConfig holds global/editor settings
type CoreConfig struct {
	Editor string `mapstructure:"editor" json:"editor"`
}

// NLRCConfig represents NetLinx Run Compiler (NLRC) configuration shared by the build and cfg commands
type NLRCConfig struct {
	Path        string   `mapstructure:"path"        json:"path"`
	IncludePath []string `mapstructure:"includePath" json:"includePath"`
	ModulePath  []string `mapstructure:"modulePath"  json:"modulePath"`
	LibraryPath []string `mapstructure:"libraryPath" json:"libraryPath"`
}

// CFGConfig represents CFG command configuration
type CFGConfig struct {
	OutputFile                string `mapstructure:"outputFile"                json:"outputFile"`
	OutputLogFile             string `mapstructure:"outputLogFile"              json:"outputLogFile"`
	OutputLogFileOption       string `mapstructure:"outputLogFileOption"        json:"outputLogFileOption"`
	OutputLogConsoleOption    bool   `mapstructure:"outputLogConsoleOption"     json:"outputLogConsoleOption"`
	BuildWithDebugInformation bool   `mapstructure:"buildWithDebugInformation" json:"buildWithDebugInformation"`
	BuildWithSource           bool   `mapstructure:"buildWithSource"            json:"buildWithSource"`
	All                       bool   `mapstructure:"all"                        json:"all"`
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
	All bool `mapstructure:"all" json:"all"`
}

var defaultConfig = Config{
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
	CFG: CFGConfig{
		OutputFile:                "build.cfg",
		OutputLogFile:             "build.log",
		OutputLogFileOption:       "N",
		OutputLogConsoleOption:    true,
		BuildWithDebugInformation: false,
		BuildWithSource:           false,
		All:                       false,
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
		All: false,
	},
}

// NormalizeConfigPaths normalizes all paths in the configuration to OS-specific format
func NormalizeConfigPaths(config *Config) {
	// Normalize NLRC paths
	config.NLRC.Path = utils.NormalizePath(config.NLRC.Path)
	config.NLRC.IncludePath = utils.NormalizePaths(config.NLRC.IncludePath)
	config.NLRC.ModulePath = utils.NormalizePaths(config.NLRC.ModulePath)
	config.NLRC.LibraryPath = utils.NormalizePaths(config.NLRC.LibraryPath)

	// Normalize Archive paths
	config.Archive.ExtraFileSearchLocations = utils.NormalizePaths(config.Archive.ExtraFileSearchLocations)
}

// LoadDefaultConfig returns the default configuration
func LoadDefaultConfig() (*Config, error) {
	config := defaultConfig

	// Normalize all paths in the default configuration
	NormalizeConfigPaths(&config)

	return &config, nil
}
