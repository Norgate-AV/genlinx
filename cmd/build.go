package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/Norgate-AV/genlinx/internal/compiler"
	"github.com/Norgate-AV/genlinx/internal/options"
)

var buildCmd = &cobra.Command{
	Use:   "build [sourceFiles...]",
	Short: "build NetLinx source file(s) or CFG file(s)",
	Long:  `Build NetLinx source files (.axs, .axi) or compile from CFG files`,
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get CLI flags
		cfgFiles, _ := cmd.Flags().GetStringSlice("cfg-files")
		sourceFiles, _ := cmd.Flags().GetStringSlice("source-files")
		includePath, _ := cmd.Flags().GetStringSlice("include-path")
		modulePath, _ := cmd.Flags().GetStringSlice("module-path")
		libraryPath, _ := cmd.Flags().GetStringSlice("library-path")
		outputPath, _ := cmd.Flags().GetString("output-path")
		all, _ := cmd.Flags().GetBool("all")
		verbose, _ := cmd.Flags().GetBool("verbose")

		// If source files provided as arguments, add them to sourceFiles
		if len(args) > 0 {
			sourceFiles = append(sourceFiles, args...)
		}

		if len(sourceFiles) == 0 && len(cfgFiles) == 0 {
			return fmt.Errorf("no source or CFG files specified")
		}

		if runtime.GOOS != "windows" {
			return fmt.Errorf("the build command is only supported on Windows")
		}

		// Create CLI options struct
		cliOpts := &options.CLIOptions{
			SourceFiles: sourceFiles,
			CFGFiles:    cfgFiles,
			IncludePath: includePath,
			ModulePath:  modulePath,
			LibraryPath: libraryPath,
			OutputPath:  outputPath,
			All:         all,
			Verbose:     verbose,
		}

		// Load and merge all configuration sources
		opts, configInfo, err := options.LoadBuildOptions(cliOpts)
		if err != nil {
			return fmt.Errorf("failed to load build options: %w", err)
		}

		// Print configuration loading information if verbose
		if verbose {
			configInfo.Print()
		}

		// TODO: Use the 'all' flag for file selection
		_ = all // Placeholder for future implementation

		// Create compiler
		nlrc := compiler.NewNLRCCompiler(opts.NLRCPath)

		// Compile
		result, err := nlrc.Compile(buildCompileOpts(opts))
		if err != nil {
			return fmt.Errorf("compilation failed: %w", err)
		}

		if len(result.Warnings) > 0 {
			fmt.Fprintln(os.Stderr, color.YellowString("A total of %d warning(s) occurred.", len(result.Warnings)))
			for _, w := range result.Warnings {
				fmt.Fprintln(os.Stderr, color.YellowString(w))
			}
		}

		if len(result.Errors) > 0 {
			fmt.Fprintln(os.Stderr, color.RedString("A total of %d error(s) occurred.", len(result.Errors)))
			for _, e := range result.Errors {
				fmt.Fprintln(os.Stderr, color.RedString(e))
			}

			return fmt.Errorf("the build process failed with a total of %d error(s)", len(result.Errors))
		}

		return nil
	},
}

// buildCompileOpts maps BuildOptions onto compiler.CompileOptions.
// Extracted so the field mapping can be unit-tested without a real compiler.
func buildCompileOpts(opts *options.BuildOptions) compiler.CompileOptions {
	return compiler.CompileOptions{
		SourceFiles: opts.SourceFiles,
		CFGFiles:    opts.CFGFiles,
		IncludePath: opts.IncludePath,
		ModulePath:  opts.ModulePath,
		LibraryPath: opts.LibraryPath,
		OutputPath:  opts.OutputPath,
		Verbose:     opts.Verbose,
	}
}

func init() {
	buildCmd.Flags().StringSliceP("cfg-files", "c", []string{}, "cfg file(s) to build from")
	buildCmd.Flags().StringSliceP("source-files", "s", []string{}, "source file(s) to build (*.axs, *.axi)")
	buildCmd.Flags().StringSliceP("include-path", "i", []string{}, "add additional include paths")
	buildCmd.Flags().StringSliceP("module-path", "m", []string{}, "add additional module paths")
	buildCmd.Flags().StringSliceP("library-path", "l", []string{}, "add additional library paths")
	buildCmd.Flags().StringP("output-path", "o", "", "set the output path for the compiled files")
	buildCmd.Flags().BoolP("all", "a", false, "select all cfg files without prompting")

	// Add conflicts (simplified for now)
	buildCmd.MarkFlagsMutuallyExclusive("cfg-files", "source-files")
	buildCmd.MarkFlagsMutuallyExclusive("cfg-files", "include-path")
	buildCmd.MarkFlagsMutuallyExclusive("cfg-files", "module-path")
	buildCmd.MarkFlagsMutuallyExclusive("cfg-files", "library-path")
}
