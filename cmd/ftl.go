package cmd

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/Norgate-AV/genlinx/internal/apw"
	"github.com/Norgate-AV/genlinx/internal/ftl"
	"github.com/Norgate-AV/genlinx/internal/prompt"
	"github.com/Norgate-AV/genlinx/internal/utils"
)

var ftlCmd = &cobra.Command{
	Use:   "ftl [workspaceFiles...]",
	Short: "generate a File Transfer List (.ftl) from a NetLinx workspace",
	Long: `Generate a File Transfer List (.ftl) from a NetLinx workspace (.apw) file.

An FTL file enumerates every transferable item — compiled master firmware
(.tkn) and panel binaries (.TP4, .TP5, .KPB) — across the selected scope of
the workspace.  The resulting file can be imported directly into NetLinx Studio
(Tools > Open Transfer List...) or the standalone AMX File Transfer Tool.

Scope can be narrowed with --project and/or --system.  Without those flags the
entire workspace is included.`,
	Args: cobra.ArbitraryArgs,
	RunE: runFTL,
}

func runFTL(cmd *cobra.Command, args []string) error {
	verbose, _ := rootCmd.PersistentFlags().GetBool("verbose")

	workspaceFiles, _ := cmd.Flags().GetStringSlice("workspace-files")
	output, _ := cmd.Flags().GetString("output")
	all, _ := cmd.Flags().GetBool("all")
	projectID, _ := cmd.Flags().GetString("project")
	systemID, _ := cmd.Flags().GetString("system")

	// Positional args are treated as additional workspace files.
	if len(args) > 0 {
		workspaceFiles = append(workspaceFiles, args...)
	}

	// Auto-discover .apw files in the current directory when none specified.
	if len(workspaceFiles) == 0 {
		if verbose {
			color.Blue("Searching for workspace files...")
		}

		found, err := utils.FindFilesByExtension(".", apw.AmxExtensions[apw.FileTypeWorkspace])
		if err != nil {
			return fmt.Errorf("error searching for workspace files: %w", err)
		}

		workspaceFiles = append(workspaceFiles, found...)

		if verbose {
			for _, f := range workspaceFiles {
				color.Cyan("  Found: %s", f)
			}
		}
	}

	if len(workspaceFiles) == 0 {
		color.Red("No workspace files found.")
		return nil
	}

	// Prompt when multiple files found and --all not set.
	if !all && len(workspaceFiles) > 1 {
		selected, err := prompt.SelectFiles(workspaceFiles)
		if err != nil {
			return fmt.Errorf("file selection failed: %w", err)
		}

		workspaceFiles = selected
	}

	for _, workspaceFile := range workspaceFiles {
		if err := generateFTL(workspaceFile, output, projectID, systemID, verbose); err != nil {
			color.Red("Error generating FTL for %s: %v", workspaceFile, err)
		}
	}

	return nil
}

func generateFTL(workspaceFile, outputPath, projectID, systemID string, verbose bool) error {
	data, err := os.ReadFile(workspaceFile)
	if err != nil {
		return fmt.Errorf("read workspace: %w", err)
	}

	workspace, err := apw.Parse(workspaceFile, data)
	if err != nil {
		return fmt.Errorf("parse workspace: %w", err)
	}

	// Resolve --system ambiguity when --project was not provided.
	if systemID != "" && projectID == "" {
		matches := workspace.FindSystemAcrossProjects(systemID)
		switch len(matches) {
		case 0:
			return fmt.Errorf("system %q not found in workspace %s", systemID, workspaceFile)
		case 1:
			projectID = matches[0].ProjectID
		default:
			names := make([]string, len(matches))
			for i, m := range matches {
				names[i] = m.ProjectID
			}

			return fmt.Errorf(
				"system %q exists in multiple projects in %s: %s — use --project to disambiguate",
				systemID, workspaceFile, strings.Join(names, ", "),
			)
		}
	}

	// Build the FTL for the requested scope.
	var list *ftl.FileTransferList

	switch {
	case projectID != "" && systemID != "":
		list = ftl.FromAPWSystem(workspace, projectID, systemID)
	case projectID != "":
		list = ftl.FromAPWProject(workspace, projectID)
	default:
		list = ftl.FromAPW(workspace)
	}

	yellow := color.New(color.FgYellow)
	for _, w := range list.Warnings {
		_, _ = yellow.Fprintln(os.Stderr, "WARNING:", w)
	}

	if len(list.Items) == 0 {
		color.Yellow("No transferable items found in %s for the selected scope.", workspaceFile)
		return nil
	}

	// Determine output filename.
	outFile := outputPath
	if outFile == "" {
		stem := strings.TrimSuffix(filepath.Base(workspaceFile), filepath.Ext(workspaceFile))
		outFile = stem + ".ftl"
	}

	if err := writeFTL(list, outFile); err != nil {
		return fmt.Errorf("write FTL: %w", err)
	}

	fmt.Printf("Generated %s (%d item(s))\n", outFile, len(list.Items))

	if verbose {
		for _, item := range list.Items {
			color.Cyan("  [%s] %s → %s", item.SystemName, item.ProjectName, filepath.Base(item.SourceFile))
		}
	}

	return nil
}

// writeFTL serialises list to outFile as an indented XML document.
func writeFTL(list *ftl.FileTransferList, outFile string) error {
	f, err := os.Create(outFile)
	if err != nil {
		return err
	}

	defer func() { _ = f.Close() }()

	if _, err := f.WriteString(xml.Header); err != nil {
		return err
	}

	enc := xml.NewEncoder(f)
	enc.Indent("", "    ")

	if err := enc.Encode(list); err != nil {
		return err
	}

	return enc.Close()
}

func init() {
	ftlCmd.Flags().StringSliceP("workspace-files", "w", []string{}, "workspace file(s) to process")
	ftlCmd.Flags().StringP("output", "o", "", "output filename (default: <workspace-name>.ftl)")
	ftlCmd.Flags().BoolP("all", "a", false, "process all discovered workspace files without prompting")
	ftlCmd.Flags().StringP("project", "p", "", "include only systems from this project")
	ftlCmd.Flags().StringP("system", "s", "", "include only this system")
}
