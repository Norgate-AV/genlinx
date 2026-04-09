package options

import (
	"github.com/Norgate-AV/genlinx/internal/archive"
)

// ArchiveCLIOptions holds the raw values parsed from archive command flags.
type ArchiveCLIOptions struct {
	WorkspaceFiles   []string
	OutputFileSuffix string
	// Unused: boolean overrides are handled via the ExplicitBoolFlags map instead.
	// IncludeCompiledSourceFiles   bool
	// NoIncludeCompiledSourceFiles bool
	// IncludeCompiledModuleFiles   bool
	// NoIncludeCompiledModuleFiles bool
	// IncludeFilesNotInWorkspace   bool
	// NoIncludeFilesNotInWorkspace bool
	// ExplicitBoolFlags tracks which boolean flags were explicitly provided by the user,
	// so the CLI can override the merged config value only when the user asked.
	ExplicitBoolFlags        map[string]bool
	ExtraFileSearchLocations []string
	ExtraGlobPatterns        []string
	All                      bool
	Verbose                  bool
	// ProjectID, when non-empty, restricts the archive to a single project.
	ProjectID string
	// SystemID, when non-empty, further restricts the archive to a single system
	// within the named project (ProjectID must also be set).
	SystemID string
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
		ExtraGlobPatterns:          mergedCfg.Archive.ExtraGlobPatterns,
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
	explicitBoolFlags := cliOpts.ExplicitBoolFlags

	if explicitBoolFlags["no-include-compiled-source-files"] {
		opts.IncludeCompiledSourceFiles = false
	} else if explicitBoolFlags["include-compiled-source-files"] {
		opts.IncludeCompiledSourceFiles = true
	}

	if explicitBoolFlags["no-include-compiled-module-files"] {
		opts.IncludeCompiledModuleFiles = false
	} else if explicitBoolFlags["include-compiled-module-files"] {
		opts.IncludeCompiledModuleFiles = true
	}

	if explicitBoolFlags["no-include-files-not-in-workspace"] {
		opts.IncludeFilesNotInWorkspace = false
	} else if explicitBoolFlags["include-files-not-in-workspace"] {
		opts.IncludeFilesNotInWorkspace = true
	}

	if len(cliOpts.ExtraFileSearchLocations) > 0 {
		opts.ExtraFileSearchLocations = prependAndDeduplicate(
			opts.ExtraFileSearchLocations,
			cliOpts.ExtraFileSearchLocations,
		)
	}

	if len(cliOpts.ExtraGlobPatterns) > 0 {
		opts.ExtraGlobPatterns = prependAndDeduplicate(
			opts.ExtraGlobPatterns,
			cliOpts.ExtraGlobPatterns,
		)
	}

	if cliOpts.All {
		opts.All = true
	}

	if cliOpts.Verbose {
		opts.Verbose = true
	}

	if cliOpts.ProjectID != "" {
		opts.ProjectID = cliOpts.ProjectID
	}

	if cliOpts.SystemID != "" {
		opts.SystemID = cliOpts.SystemID
	}

	return opts, configInfo, nil
}
