package options

import (
	"github.com/Norgate-AV/genlinx/internal/archive"
)

// ArchiveCLIOptions holds the raw values parsed from archive command flags.
type ArchiveCLIOptions struct {
	WorkspaceFiles               []string
	OutputFileSuffix             string
	IncludeCompiledSourceFiles   bool
	NoIncludeCompiledSourceFiles bool
	IncludeCompiledModuleFiles   bool
	NoIncludeCompiledModuleFiles bool
	IncludeFilesNotInWorkspace   bool
	NoIncludeFilesNotInWorkspace bool
	// Changed tracks which boolean flags were explicitly provided by the user,
	// so the CLI can override the merged config value only when the user asked.
	Changed                  map[string]bool
	ExtraFileSearchLocations []string
	ExtraFileArchiveLocation string
	All                      bool
	Verbose                  bool
}

// LoadArchiveOptions loads and merges archive options from defaults, global
// config, local config, and CLI flags – in that order of precedence.
func LoadArchiveOptions(cliOpts *ArchiveCLIOptions) (*archive.Options, *ConfigLoadInfo, error) {
	mergedCfg, configInfo, err := LoadMergedConfig()
	if err != nil {
		return nil, nil, err
	}

	opts := &archive.Options{
		OutputFileSuffix:           mergedCfg.Archive.OutputFile,
		IncludeCompiledSourceFiles: mergedCfg.Archive.IncludeCompiledSourceFiles,
		IncludeCompiledModuleFiles: mergedCfg.Archive.IncludeCompiledModuleFiles,
		IncludeFilesNotInWorkspace: mergedCfg.Archive.IncludeFilesNotInWorkspace,
		ExtraFileSearchLocations:   mergedCfg.Archive.ExtraFileSearchLocations,
		ExtraFileArchiveLocation:   mergedCfg.Archive.ExtraFileArchiveLocation,
		All:                        mergedCfg.Archive.All,
		IgnoredFiles:               mergedCfg.Archive.IgnoredFiles,
	}

	if cliOpts == nil {
		return opts, configInfo, nil
	}

	// Apply CLI overrides (highest precedence).

	if cliOpts.OutputFileSuffix != "" {
		opts.OutputFileSuffix = cliOpts.OutputFileSuffix
	}

	// Boolean flags: only override config when user explicitly passed the flag.
	changed := cliOpts.Changed

	if changed["no-include-compiled-source-files"] {
		opts.IncludeCompiledSourceFiles = false
	} else if changed["include-compiled-source-files"] {
		opts.IncludeCompiledSourceFiles = true
	}

	if changed["no-include-compiled-module-files"] {
		opts.IncludeCompiledModuleFiles = false
	} else if changed["include-compiled-module-files"] {
		opts.IncludeCompiledModuleFiles = true
	}

	if changed["no-include-files-not-in-workspace"] {
		opts.IncludeFilesNotInWorkspace = false
	} else if changed["include-files-not-in-workspace"] {
		opts.IncludeFilesNotInWorkspace = true
	}

	if len(cliOpts.ExtraFileSearchLocations) > 0 {
		opts.ExtraFileSearchLocations = prependAndDeduplicate(
			opts.ExtraFileSearchLocations,
			cliOpts.ExtraFileSearchLocations,
		)
	}

	if cliOpts.ExtraFileArchiveLocation != "" {
		opts.ExtraFileArchiveLocation = cliOpts.ExtraFileArchiveLocation
	}

	if cliOpts.All {
		opts.All = true
	}

	if cliOpts.Verbose {
		opts.Verbose = true
	}

	return opts, configInfo, nil
}
