package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"

	"github.com/Norgate-AV/genlinx/internal/studio"
)

var studioCmd = &cobra.Command{
	Use:   "studio",
	Short: "NetLinx Studio utilities",
	Long:  `Tools for interacting with NetLinx Studio settings and configuration.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if runtime.GOOS != "windows" {
			return fmt.Errorf("the studio command is only supported on Windows")
		}

		if !studio.IsInstalled() {
			return fmt.Errorf("NetLinx Studio does not appear to be installed on this machine")
		}

		return nil
	},
}

var studioBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup NetLinx Studio settings to an .epx file",
	Long: `Read all NetLinx Studio settings from the Windows registry and write a
compatible .epx preferences file. The output file can be re-imported directly
using NetLinx Studio's native Tools > Import Preferences... option.`,
	RunE: studioBackupRunE,
}

var studioRestoreCmd = &cobra.Command{
	Use:   "restore [file]",
	Short: "restore NetLinx Studio settings from an .epx file",
	Long: `Read a NetLinx Studio .epx preferences file and apply all settings
directly to the Windows registry. By default a timestamped backup of the
current registry state is written before any changes are made.

The input file may be provided as a flag (--input/-i) or as a positional
argument.`,
	Args: cobra.MaximumNArgs(1),
	RunE: studioRestoreRunE,
}

func init() {
	studioBackupCmd.Flags().StringP("output", "o", "", "output file path (default: netlinx-studio-backup-<timestamp>.epx in current directory)")
	studioCmd.AddCommand(studioBackupCmd)

	studioRestoreCmd.Flags().StringP("input", "i", "", "path to the .epx backup file to restore")
	studioRestoreCmd.Flags().BoolP("dry-run", "n", false, "print what would be written without touching the registry")
	studioRestoreCmd.Flags().BoolP("backup", "b", true, "write a registry backup before applying changes")
	studioRestoreCmd.Flags().BoolP("force", "f", false, "skip the confirmation prompt")
	studioCmd.AddCommand(studioRestoreCmd)
}

// runDryRun prints the pre-computed diff in a table without touching the registry.
func runDryRun(diff []studio.DiffEntry) error {
	fmt.Printf("Dry run — %d entries would be written:\n\n", len(diff))
	return printDiffTable(diff)
}

// printDiffTable renders diff as a table to stdout.
func printDiffTable(diff []studio.DiffEntry) error {
	green := renderer.Tint{FG: renderer.Colors{color.FgGreen}}

	r := renderer.NewColorized(renderer.ColorizedConfig{
		Border:    green,
		Separator: green,
	})

	table := tablewriter.NewTable(
		os.Stdout,
		tablewriter.WithRenderer(r),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignCenter},
			},
		}),
	)

	table.Header("Status", "Hive", "Subkey", "Name", "Old Value", "New Value")

	for _, d := range diff {
		if d.Status == studio.DiffReplaced {
			oldCount, _ := d.OldValue.(int)
			hive := "HKCU"
			if d.HiveLM {
				hive = "HKLM"
			}

			newCount := len(d.TCPIPReplace)
			if d.DirReplace != nil {
				newCount = len(d.DirReplace)
			}

			_ = table.Append(
				string(d.Status),
				hive,
				d.SubKey,
				"(all entries)",
				fmt.Sprintf("%d entries", oldCount),
				fmt.Sprintf("%d entries", newCount),
			)
			continue
		}

		hive := "HKCU"
		if d.HiveLM {
			hive = "HKLM"
		}

		oldVal := ""
		if d.OldValue != nil {
			oldVal = fmt.Sprintf("%v", d.OldValue)
		}

		_ = table.Append(string(d.Status), hive, d.SubKey, d.ValueName, oldVal, fmt.Sprintf("%v", d.Value))
	}

	return table.Render()
}
