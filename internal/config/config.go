package config

import (
	"github.com/Norgate-AV/genlinx/internal/utils"
)

// Config represents the application configuration
type Config struct {
	Core     CoreConfig     `mapstructure:"core"     json:"core"`
	Compiler CompilerConfig `mapstructure:"compiler" json:"compiler"`
	CFG      CFGConfig      `mapstructure:"cfg"      json:"cfg"`
	Archive  ArchiveConfig  `mapstructure:"archive"  json:"archive"`
}

// CoreConfig holds global/editor settings
type CoreConfig struct {
	Editor string `mapstructure:"editor" json:"editor"`
}

// CompilerConfig represents NetLinx compiler configuration shared by the build and cfg commands
type CompilerConfig struct {
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
}

// ArchiveConfig represents archive command configuration
type ArchiveConfig struct {
	OutputFile                 string   `mapstructure:"outputFile"                 json:"outputFile"`
	IncludeCompiledSourceFiles bool     `mapstructure:"includeCompiledSourceFiles" json:"includeCompiledSourceFiles"`
	IncludeCompiledModuleFiles bool     `mapstructure:"includeCompiledModuleFiles" json:"includeCompiledModuleFiles"`
	IncludeFilesNotInWorkspace bool     `mapstructure:"includeFilesNotInWorkspace" json:"includeFilesNotInWorkspace"`
	ExtraFileSearchLocations   []string `mapstructure:"extraFileSearchLocations"   json:"extraFileSearchLocations"`
	IgnoredFiles               []string `mapstructure:"ignoredFiles"               json:"ignoredFiles"`
}

var defaultConfig = Config{
	Compiler: CompilerConfig{
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
	},
	Archive: ArchiveConfig{
		OutputFile:                 "archive.zip",
		IncludeCompiledSourceFiles: true,
		IncludeCompiledModuleFiles: true,
		IncludeFilesNotInWorkspace: true,
		ExtraFileSearchLocations:   utils.NormalizePaths([]string{"C:/Program Files (x86)/Common Files/AMXShare"}),
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
}

// NormalizeConfigPaths normalizes all paths in the configuration to OS-specific format
func NormalizeConfigPaths(config *Config) {
	// Normalize compiler paths
	config.Compiler.Path = utils.NormalizePath(config.Compiler.Path)
	config.Compiler.IncludePath = utils.NormalizePaths(config.Compiler.IncludePath)
	config.Compiler.ModulePath = utils.NormalizePaths(config.Compiler.ModulePath)
	config.Compiler.LibraryPath = utils.NormalizePaths(config.Compiler.LibraryPath)

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
