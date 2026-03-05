package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Norgate-AV/genlinx-go/internal/apw"
	"github.com/Norgate-AV/genlinx-go/internal/archive"
	"github.com/Norgate-AV/genlinx-go/internal/options"
	"github.com/Norgate-AV/genlinx-go/internal/utils"
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
	extraArchiveLocation, _ := cmd.Flags().GetString("extra-file-archive-location")
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
		ExtraFileArchiveLocation: extraArchiveLocation,
		All:                      all,
		Verbose:                  verbose,
		Changed:                  changed,
	}

	opts, configInfo, err := options.LoadArchiveOptions(cliOpts)
	if err != nil {
		return fmt.Errorf("failed to load archive options: %w", err)
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
	// Build archives
	// -----------------------------------------------------------------------

	for _, workspaceFile := range workspaceFiles {
		if verbose {
			fmt.Printf("Generating archive for %s...\n", workspaceFile)
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

		builder := archive.NewBuilder(workspace, opts)
		if err := builder.Build(); err != nil {
			fmt.Fprintf(os.Stderr, "Error building archive for %s: %v\n", workspaceFile, err)
			continue
		}
	}

	return nil
}

// selectWorkspaceFiles presents a numbered list of workspace files and asks the
// user to pick one or more by index (comma-separated). Pressing Enter without
// input selects all files. This mirrors the TypeScript @inquirer/prompts
// checkbox behaviour.
func selectWorkspaceFiles(files []string) ([]string, error) {
	fmt.Println("Select workspace file(s):")

	for i, f := range files {
		fmt.Printf("  [%d] %s\n", i, f)
	}

	fmt.Print("Enter comma-separated indices (or press Enter to select all): ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	line := strings.TrimSpace(scanner.Text())

	if line == "" {
		return files, nil
	}

	var selected []string

	for _, part := range strings.Split(line, ",") {
		part = strings.TrimSpace(part)

		idx, err := strconv.Atoi(part)
		if err != nil || idx < 0 || idx >= len(files) {
			return nil, fmt.Errorf("invalid selection %q", part)
		}

		selected = append(selected, files[idx])
	}

	if len(selected) == 0 {
		return nil, fmt.Errorf("you must choose at least one file")
	}

	return selected, nil
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
	archiveCmd.Flags().StringP("extra-file-archive-location", "p", "", "location to place extra files in the archive")
	archiveCmd.Flags().BoolP("all", "a", false, "process all found workspace files without prompting")
}
