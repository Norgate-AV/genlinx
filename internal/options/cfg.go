package options

import (
	"github.com/Norgate-AV/genlinx/internal/cfg"
)

// CfgCLIOptions holds the raw values parsed from the cfg command flags.
type CfgCLIOptions struct {
	WorkspaceFiles      []string
	RootDirectory       string
	OutputFileSuffix    string
	OutputLogFileSuffix string
	OutputLogFileOption string
	// Unused: boolean overrides are handled via the ExplicitBoolFlags map instead.
	// OutputLogConsoleOption      bool
	// NoOutputLogConsoleOption    bool
	// BuildWithDebugInformation   bool
	// NoBuildWithDebugInformation bool
	// BuildWithSource             bool
	// NoBuildWithSource           bool
	IncludePath []string
	ModulePath  []string
	LibraryPath []string
	All         bool
	Verbose     bool
	// ExplicitBoolFlags tracks which boolean flags were explicitly provided by the user.
	ExplicitBoolFlags map[string]bool
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
		IncludePath:               mergedCfg.Compiler.IncludePath,
		ModulePath:                mergedCfg.Compiler.ModulePath,
		LibraryPath:               mergedCfg.Compiler.LibraryPath,
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

	explicitBoolFlags := cliOpts.ExplicitBoolFlags

	if explicitBoolFlags["no-output-log-console-option"] {
		opts.OutputLogConsoleOption = false
	} else if explicitBoolFlags["output-log-console-option"] {
		opts.OutputLogConsoleOption = true
	}

	if explicitBoolFlags["no-build-with-debug-information"] {
		opts.BuildWithDebugInformation = false
	} else if explicitBoolFlags["build-with-debug-information"] {
		opts.BuildWithDebugInformation = true
	}

	if explicitBoolFlags["no-build-with-source"] {
		opts.BuildWithSource = false
	} else if explicitBoolFlags["build-with-source"] {
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
