package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/Norgate-AV/genlinx/internal/compiler"
	"github.com/Norgate-AV/genlinx/internal/options"
	"github.com/Norgate-AV/genlinx/internal/prompt"
	"github.com/Norgate-AV/genlinx/internal/utils"
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
		verbose, _ := rootCmd.PersistentFlags().GetBool("verbose")

		// Positional args are treated as additional source files
		if len(args) > 0 {
			sourceFiles = append(sourceFiles, args...)
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

		if verbose {
			configInfo.Print()
		}

		c := compiler.NewNetLinxCompiler(opts.CompilerPath)

		if _, err := os.Stat(opts.CompilerPath); err != nil {
			return fmt.Errorf("compiler not found at %q — is NetLinx Studio installed?\n%w", opts.CompilerPath, err)
		}

		if len(opts.SourceFiles) > 0 {
			return executeSourceBuild(opts.SourceFiles, c, opts)
		}

		return executeCfgBuild(opts.CFGFiles, c, opts)
	},
}

// executeSourceBuild compiles each source file with a separate compiler invocation.
func executeSourceBuild(files []string, c compiler.Compiler, opts *options.BuildOptions) error {
	for _, file := range files {
		if opts.Verbose {
			color.Blue("Executing build for %s...", file)
		}

		compileOpts := buildCompileOpts(opts)
		compileOpts.SourceFiles = []string{file}
		compileOpts.CFGFiles = nil

		result, err := c.Compile(compileOpts)
		if err != nil {
			return fmt.Errorf("compilation failed for %s: %w", file, err)
		}

		if err := printBuildResult(result); err != nil {
			return err
		}
	}

	return nil
}

// executeCfgBuild compiles each CFG file with a separate compiler invocation.
// When no files are given it auto-discovers .cfg files in the current working
// directory. If multiple files are found and --all is not set the user is prompted
// to select which ones to build.
func executeCfgBuild(files []string, c compiler.Compiler, opts *options.BuildOptions) error {
	if len(files) == 0 {
		if opts.Verbose {
			color.Blue("Searching for CFG files...")
		}

		located, err := utils.FindFilesByExtension(".", ".cfg")
		if err != nil {
			return fmt.Errorf("failed to search for CFG files: %w", err)
		}

		files = located
	}

	if len(files) == 0 {
		color.Red("No CFG files found.")
		return nil
	}

	if !opts.All && len(files) > 1 {
		selected, err := prompt.SelectFiles(files)
		if err != nil {
			return err
		}

		files = selected
	}

	for _, file := range files {
		if opts.Verbose {
			color.Blue("Executing build for %s...", file)
		}

		compileOpts := buildCompileOpts(opts)
		compileOpts.CFGFiles = []string{file}
		compileOpts.SourceFiles = nil

		result, err := c.Compile(compileOpts)
		if err != nil {
			return fmt.Errorf("compilation failed for %s: %w", file, err)
		}

		if err := printBuildResult(result); err != nil {
			return err
		}
	}

	return nil
}

// printBuildResult prints warnings and errors from a compile result and returns
// an error if any errors were present.
func printBuildResult(result *compiler.CompileResult) error {
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

	buildCmd.MarkFlagsMutuallyExclusive("cfg-files", "source-files")
	buildCmd.MarkFlagsMutuallyExclusive("cfg-files", "include-path")
	buildCmd.MarkFlagsMutuallyExclusive("cfg-files", "module-path")
	buildCmd.MarkFlagsMutuallyExclusive("cfg-files", "library-path")
}
