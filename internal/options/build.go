package options

import (
	"fmt"
	"path/filepath"
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

	// Compiler path
	CompilerPath string
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

// LoadBuildOptions loads and merges build options from all sources
func LoadBuildOptions(cliOpts *CLIOptions) (*BuildOptions, *ConfigLoadInfo, error) {
	mergedCfg, configInfo, err := LoadMergedConfig()
	if err != nil {
		return nil, nil, err
	}

	// Create build options from merged config
	opts := &BuildOptions{
		IncludePath:  mergedCfg.Compiler.IncludePath,
		ModulePath:   mergedCfg.Compiler.ModulePath,
		LibraryPath:  mergedCfg.Compiler.LibraryPath,
		CompilerPath: mergedCfg.Compiler.Path,
		All:          mergedCfg.Build.All,
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

	// Resolve compiler path
	if o.CompilerPath != "" {
		if o.CompilerPath, err = filepath.Abs(o.CompilerPath); err != nil {
			return fmt.Errorf("failed to resolve compiler path %s: %w", o.CompilerPath, err)
		}
	}

	return nil
}
