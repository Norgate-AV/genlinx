package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/Norgate-AV/genlinx/internal/apw"
	"github.com/Norgate-AV/genlinx/internal/archive"
	"github.com/Norgate-AV/genlinx/internal/options"
	"github.com/Norgate-AV/genlinx/internal/prompt"
	"github.com/Norgate-AV/genlinx/internal/utils"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "generate a NetLinx workspace zip archive",
	Long:  `Generate a NetLinx workspace zip archive`,
	RunE:  runArchive,
}

func runArchive(cmd *cobra.Command, _ []string) error {
	verbose, _ := rootCmd.PersistentFlags().GetBool("verbose")

	// Collect CLI flag values.
	workspaceFiles, _ := cmd.Flags().GetStringSlice("workspace-files")
	outputFileSuffix, _ := cmd.Flags().GetString("output-file-suffix")
	extraSearchLocations, _ := cmd.Flags().GetStringSlice("extra-file-search-locations")
	all, _ := cmd.Flags().GetBool("all")

	// Build a Changed map so LoadArchiveOptions can distinguish "user set this
	// flag" from "Cobra zero-value default".
	changed := map[string]bool{
		"include-compiled-source-files":     cmd.Flags().Changed("include-compiled-source-files"),
		"no-include-compiled-source-files":  cmd.Flags().Changed("no-include-compiled-source-files"),
		"include-compiled-module-files":     cmd.Flags().Changed("include-compiled-module-files"),
		"no-include-compiled-module-files":  cmd.Flags().Changed("no-include-compiled-module-files"),
		"include-files-not-in-workspace":    cmd.Flags().Changed("include-files-not-in-workspace"),
		"no-include-files-not-in-workspace": cmd.Flags().Changed("no-include-files-not-in-workspace"),
	}

	cliOpts := &options.ArchiveCLIOptions{
		WorkspaceFiles:           workspaceFiles,
		OutputFileSuffix:         outputFileSuffix,
		ExtraFileSearchLocations: extraSearchLocations,
		All:                      all,
		Verbose:                  verbose,
		Changed:                  changed,
	}

	opts, configInfo, err := options.LoadArchiveOptions(cliOpts)
	if err != nil {
		return fmt.Errorf("failed to load archive options: %w", err)
	}

	if verbose {
		configInfo.Print()
	}

	// -----------------------------------------------------------------------
	// Workspace file discovery
	// -----------------------------------------------------------------------

	if len(workspaceFiles) == 0 {
		if verbose {
			color.Blue("Searching for workspace files...")
		}

		found, err := utils.FindFilesByExtension(".", apw.AmxExtensions[apw.FileTypeWorkspace])
		if err != nil {
			return fmt.Errorf("error searching for workspace files: %w", err)
		}

		workspaceFiles = append(workspaceFiles, found...)

		if verbose && len(workspaceFiles) > 0 {
			for _, f := range workspaceFiles {
				color.Cyan("  Found: %s", f)
			}
		}
	}

	if len(workspaceFiles) == 0 {
		color.Red("No workspace files found.")
		return nil
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
	// Build archives
	// -----------------------------------------------------------------------

	for _, workspaceFile := range workspaceFiles {
		if verbose {
			color.Blue("Generating archive for %s...", workspaceFile)
		}

		data, err := os.ReadFile(workspaceFile)
		if err != nil {
			color.Red("Error loading workspace %s: %v", workspaceFile, err)
			continue
		}

		workspace, err := apw.Parse(workspaceFile, data)
		if err != nil {
			color.Red("Error loading workspace %s: %v", workspaceFile, err)
			continue
		}

		builder := archive.NewBuilder(workspace, opts)
		if err := builder.Build(); err != nil {
			color.Red("Error building archive for %s: %v", workspaceFile, err)
			continue
		}
	}

	return nil
}

func init() {
	archiveCmd.Flags().StringSliceP("workspace-files", "w", []string{}, "workspace file(s) to generate archive(s) for (default: search current directory)")
	archiveCmd.Flags().StringP("output-file-suffix", "o", "", "output file suffix")
	archiveCmd.Flags().BoolP("include-compiled-source-files", "s", false, "include compiled source files")
	archiveCmd.Flags().BoolP("no-include-compiled-source-files", "S", false, "do not include compiled source files")
	archiveCmd.Flags().BoolP("include-compiled-module-files", "m", false, "include compiled module files")
	archiveCmd.Flags().BoolP("no-include-compiled-module-files", "M", false, "do not include compiled module files")
	archiveCmd.Flags().BoolP("include-files-not-in-workspace", "n", false, "include files not in workspace")
	archiveCmd.Flags().BoolP("no-include-files-not-in-workspace", "N", false, "do not include files not in workspace")
	archiveCmd.Flags().StringSliceP("extra-file-search-locations", "l", []string{}, "extra file locations to search")
	archiveCmd.Flags().BoolP("all", "a", false, "process all found workspace files without prompting")

	archiveCmd.MarkFlagsMutuallyExclusive("include-compiled-source-files", "no-include-compiled-source-files")
	archiveCmd.MarkFlagsMutuallyExclusive("include-compiled-module-files", "no-include-compiled-module-files")
	archiveCmd.MarkFlagsMutuallyExclusive("include-files-not-in-workspace", "no-include-files-not-in-workspace")
}
