package cmd

import (
	"fmt"
	"os"

	"github.com/Norgate-AV/genlinx-go/internal/apw"
	gocfg "github.com/Norgate-AV/genlinx-go/internal/cfg"
	"github.com/Norgate-AV/genlinx-go/internal/options"
	"github.com/Norgate-AV/genlinx-go/internal/utils"
	"github.com/spf13/cobra"
)

var cfgCmd = &cobra.Command{
	Use:   "cfg",
	Short: "generate NetLinx build CFG files",
	Long:  `Generate NetLinx build configuration files`,
	RunE:  runCfg,
}

func runCfg(cmd *cobra.Command, _ []string) error {
	verbose, _ := rootCmd.PersistentFlags().GetBool("verbose")

	// Collect CLI flag values.
	workspaceFiles, _ := cmd.Flags().GetStringSlice("workspace-files")
	rootDirectory, _ := cmd.Flags().GetString("root-directory")
	outputFileSuffix, _ := cmd.Flags().GetString("output-file-suffix")
	outputLogFileSuffix, _ := cmd.Flags().GetString("output-log-file-suffix")
	outputLogFileOption, _ := cmd.Flags().GetString("output-log-file-option")
	includePath, _ := cmd.Flags().GetStringSlice("include-path")
	modulePath, _ := cmd.Flags().GetStringSlice("module-path")
	libraryPath, _ := cmd.Flags().GetStringSlice("library-path")
	all, _ := cmd.Flags().GetBool("all")

	// Build a Changed map so LoadCfgOptions can distinguish user-set flags
	// from Cobra zero-value defaults.
	changed := map[string]bool{
		"output-log-console-option":       cmd.Flags().Changed("output-log-console-option"),
		"no-output-log-console-option":    cmd.Flags().Changed("no-output-log-console-option"),
		"build-with-debug-information":    cmd.Flags().Changed("build-with-debug-information"),
		"no-build-with-debug-information": cmd.Flags().Changed("no-build-with-debug-information"),
		"build-with-source":               cmd.Flags().Changed("build-with-source"),
		"no-build-with-source":            cmd.Flags().Changed("no-build-with-source"),
	}

	cliOpts := &options.CfgCLIOptions{
		WorkspaceFiles:      workspaceFiles,
		RootDirectory:       rootDirectory,
		OutputFileSuffix:    outputFileSuffix,
		OutputLogFileSuffix: outputLogFileSuffix,
		OutputLogFileOption: outputLogFileOption,
		IncludePath:         includePath,
		ModulePath:          modulePath,
		LibraryPath:         libraryPath,
		All:                 all,
		Verbose:             verbose,
		Changed:             changed,
	}

	opts, configInfo, err := options.LoadCfgOptions(cliOpts)
	if err != nil {
		return fmt.Errorf("failed to load cfg options: %w", err)
	}

	if verbose {
		fmt.Println("Configuration loading:")
		fmt.Println("  Default config: Built-in defaults loaded")

		if configInfo.GlobalResult.Found {
			fmt.Printf("  Global config: Loaded from %s\n", configInfo.GlobalResult.Path)
		} else {
			fmt.Println("  Global config: No config found")
		}

		if configInfo.LocalResult.Found {
			fmt.Printf("  Local config: Loaded from %s\n", configInfo.LocalResult.Path)
		} else {
			fmt.Println("  Local config: No config found")
		}
	}

	// -----------------------------------------------------------------------
	// Workspace file discovery
	// -----------------------------------------------------------------------

	if len(workspaceFiles) == 0 {
		if verbose {
			fmt.Println("Searching for workspace files...")
		}

		found, err := utils.FindFilesByExtension(".", apw.AmxExtensions[apw.FileTypeWorkspace])
		if err != nil {
			return fmt.Errorf("error searching for workspace files: %w", err)
		}

		workspaceFiles = append(workspaceFiles, found...)

		if verbose && len(workspaceFiles) > 0 {
			for _, f := range workspaceFiles {
				fmt.Printf("  Found: %s\n", f)
			}
		}
	}

	if len(workspaceFiles) == 0 {
		fmt.Fprintln(os.Stderr, "No workspace files found.")
		os.Exit(0)
	}

	// -----------------------------------------------------------------------
	// Interactive file selection (mirrors shouldPromptUser / selectFiles)
	// -----------------------------------------------------------------------

	if !opts.All && len(workspaceFiles) > 1 {
		selected, err := selectWorkspaceFiles(workspaceFiles)
		if err != nil {
			return fmt.Errorf("file selection failed: %w", err)
		}

		workspaceFiles = selected
	}

	// -----------------------------------------------------------------------
	// Build CFG files
	// -----------------------------------------------------------------------

	for _, workspaceFile := range workspaceFiles {
		if verbose {
			fmt.Printf("Generating CFG for %s...\n", workspaceFile)
		}

		workspace := apw.New(workspaceFile)
		if err := workspace.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading workspace %s: %v\n", workspaceFile, err)
			continue
		}

		builder := gocfg.NewBuilder(workspace, opts)
		content := builder.Build()

		outputFile := fmt.Sprintf("%s.%s", workspace.ID(), opts.OutputFileSuffix)

		if err := os.WriteFile(outputFile, []byte(content), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing CFG file %s: %v\n", outputFile, err)
			continue
		}

		fmt.Printf("Created CFG file: %s\n", outputFile)
	}

	return nil
}

func init() {
	cfgCmd.Flags().StringSliceP("workspace-files", "w", []string{}, "workspace file(s) to generate a CFG for (default: search current directory)")
	cfgCmd.Flags().StringP("root-directory", "r", "", "root directory reference (default: parent directory of workspace file)")
	cfgCmd.Flags().StringP("output-file-suffix", "o", "", "output file suffix")
	cfgCmd.Flags().StringP("output-log-file-suffix", "f", "", "output log file suffix")
	cfgCmd.Flags().StringP("output-log-file-option", "k", "", "output log file option (A = append, N = overwrite)")
	cfgCmd.Flags().BoolP("output-log-console-option", "c", false, "output log to console")
	cfgCmd.Flags().BoolP("no-output-log-console-option", "C", false, "do not output log to console")
	cfgCmd.Flags().BoolP("build-with-debug-information", "d", false, "build with debug information")
	cfgCmd.Flags().BoolP("no-build-with-debug-information", "D", false, "do not build with debug information")
	cfgCmd.Flags().BoolP("build-with-source", "s", false, "build with source")
	cfgCmd.Flags().BoolP("no-build-with-source", "S", false, "do not build with source")
	cfgCmd.Flags().StringSliceP("include-path", "i", []string{}, "add additional include paths")
	cfgCmd.Flags().StringSliceP("module-path", "m", []string{}, "add additional module paths")
	cfgCmd.Flags().StringSliceP("library-path", "l", []string{}, "add additional library paths")
	cfgCmd.Flags().BoolP("all", "a", false, "process all found workspace files without prompting")
}
