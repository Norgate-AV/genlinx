package options

import (
	"github.com/Norgate-AV/genlinx/internal/cfg"
)

// CfgCLIOptions holds the raw values parsed from the cfg command flags.
type CfgCLIOptions struct {
	WorkspaceFiles              []string
	RootDirectory               string
	OutputFileSuffix            string
	OutputLogFileSuffix         string
	OutputLogFileOption         string
	OutputLogConsoleOption      bool
	NoOutputLogConsoleOption    bool
	BuildWithDebugInformation   bool
	NoBuildWithDebugInformation bool
	BuildWithSource             bool
	NoBuildWithSource           bool
	IncludePath                 []string
	ModulePath                  []string
	LibraryPath                 []string
	All                         bool
	Verbose                     bool
	// Changed tracks which boolean flags were explicitly provided by the user.
	Changed map[string]bool
}

// LoadCfgOptions loads and merges CFG options from defaults, global config,
// local config, and CLI flags – in that order of precedence.
func LoadCfgOptions(cliOpts *CfgCLIOptions) (*cfg.Options, *ConfigLoadInfo, error) {
	mergedCfg, configInfo, err := LoadMergedConfig()
	if err != nil {
		return nil, nil, err
	}

	opts := &cfg.Options{
		OutputFileSuffix:          mergedCfg.CFG.OutputFile,
		OutputLogFileSuffix:       mergedCfg.CFG.OutputLogFile,
		OutputLogFileOption:       mergedCfg.CFG.OutputLogFileOption,
		OutputLogConsoleOption:    mergedCfg.CFG.OutputLogConsoleOption,
		BuildWithDebugInformation: mergedCfg.CFG.BuildWithDebugInformation,
		BuildWithSource:           mergedCfg.CFG.BuildWithSource,
		IncludePath:               mergedCfg.CFG.IncludePath,
		ModulePath:                mergedCfg.CFG.ModulePath,
		LibraryPath:               mergedCfg.CFG.LibraryPath,
		All:                       mergedCfg.CFG.All,
	}

	if cliOpts == nil {
		return opts, configInfo, nil
	}

	// Apply CLI overrides (highest precedence).

	if cliOpts.RootDirectory != "" {
		opts.RootDirectory = cliOpts.RootDirectory
	}

	if cliOpts.OutputFileSuffix != "" {
		opts.OutputFileSuffix = cliOpts.OutputFileSuffix
	}

	if cliOpts.OutputLogFileSuffix != "" {
		opts.OutputLogFileSuffix = cliOpts.OutputLogFileSuffix
	}

	if cliOpts.OutputLogFileOption != "" {
		opts.OutputLogFileOption = cliOpts.OutputLogFileOption
	}

	changed := cliOpts.Changed

	if changed["no-output-log-console-option"] {
		opts.OutputLogConsoleOption = false
	} else if changed["output-log-console-option"] {
		opts.OutputLogConsoleOption = true
	}

	if changed["no-build-with-debug-information"] {
		opts.BuildWithDebugInformation = false
	} else if changed["build-with-debug-information"] {
		opts.BuildWithDebugInformation = true
	}

	if changed["no-build-with-source"] {
		opts.BuildWithSource = false
	} else if changed["build-with-source"] {
		opts.BuildWithSource = true
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

	if cliOpts.All {
		opts.All = true
	}

	if cliOpts.Verbose {
		opts.Verbose = true
	}

	return opts, configInfo, nil
}
