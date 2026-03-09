package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Norgate-AV/genlinx/internal/apw"
	"github.com/Norgate-AV/genlinx/internal/cfg"
	"github.com/Norgate-AV/genlinx/internal/options"
	"github.com/Norgate-AV/genlinx/internal/prompt"
	"github.com/Norgate-AV/genlinx/internal/utils"
)

var cfgCmd = &cobra.Command{
	Use:   "cfg",
	Short: "generate NetLinx build CFG files",
	Long: `Generate NetLinx build configuration files.

The CFG file is used to configure the NetLinx Studio build process.

If an option is omitted, genlinx will:
  1. look for a .genlinxrc.json file in the root directory of the project;
     if one is found, any options defined will be set to the value in the file
  2. use the default values from the global config file for any remaining options

The output file will be named the same as the workspace ID combined with the suffix.
The default suffix is "build.cfg".

For example, if the workspace ID is "SomeAwesomeProject", the output file will be
"SomeAwesomeProject.build.cfg", unless the -o option is used to specify a suffix.

Examples:
  genlinx cfg -w workspace.apw                                     generate CFG for workspace.apw
  genlinx cfg -w workspace.apw -s                                  build with source
  genlinx cfg -w workspace.apw -D                                  do not build with debug information
  genlinx cfg -a                                                   search for and automatically select all workspace files
  genlinx cfg -A                                                   search for and prompt to select workspace files
  genlinx cfg -w workspace.apw -i \\path\\to\\includes             add additional include paths
  genlinx cfg -w workspace.apw -m \\path\\to\\modules              add additional module paths
  genlinx cfg -w workspace.apw -l \\path\\to\\libraries            add additional library paths`,
	RunE: runCfg,
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
	noAll, _ := cmd.Flags().GetBool("no-all")

	if outputLogFileOption != "" && outputLogFileOption != "A" && outputLogFileOption != "N" {
		return fmt.Errorf("invalid value %q for --output-log-file-option: must be A or N", outputLogFileOption)
	}

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

	// -A/--no-all explicitly overrides a config-file all:true.
	if noAll {
		opts.All = false
	}

	if verbose {
		configInfo.Print()
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
		selected, err := prompt.SelectFiles(workspaceFiles)
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

		data, err := os.ReadFile(workspaceFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading workspace %s: %v\n", workspaceFile, err)
			continue
		}

		workspace, err := apw.Parse(workspaceFile, data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading workspace %s: %v\n", workspaceFile, err)
			continue
		}

		builder := cfg.NewBuilder(workspace, opts)
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
	cfgCmd.Flags().BoolP("no-all", "A", false, "prompt to select workspace files even when multiple are found")
}
