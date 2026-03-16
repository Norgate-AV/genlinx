package cmd

import (
	"encoding/xml"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"

	"github.com/Norgate-AV/genlinx/internal/prompt"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = time.Now().Format("netlinx-studio-backup-2006-01-02-150405") + ".epx"
		}

		fmt.Println("Reading NetLinx Studio settings from registry...")

		settings, err := studio.ReadRegistry()
		if err != nil {
			return fmt.Errorf("read registry: %w", err)
		}

		prefs := studio.BuildPreferences(settings)

		xmlBytes, err := xml.MarshalIndent(prefs, "", "    ")
		if err != nil {
			return fmt.Errorf("marshal preferences: %w", err)
		}

		content := append([]byte(xml.Header), xmlBytes...)
		content = append(content, '\n')

		if err := os.WriteFile(output, content, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", output, err)
		}

		fmt.Printf("Backup written to: %s\n", output)
		return nil
	},
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

var studioRestoreCmd = &cobra.Command{
	Use:   "restore [file]",
	Short: "restore NetLinx Studio settings from an .epx file",
	Long: `Read a NetLinx Studio .epx preferences file and apply all settings
directly to the Windows registry. By default a timestamped backup of the
current registry state is written before any changes are made.

The input file may be provided as a flag (--input/-i) or as a positional
argument.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input, _ := cmd.Flags().GetString("input")
		if input == "" && len(args) > 0 {
			input = args[0]
		}

		if input == "" {
			return fmt.Errorf("no input file specified: use --input/-i or pass the file as an argument")
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		doBackup, _ := cmd.Flags().GetBool("backup")
		force, _ := cmd.Flags().GetBool("force")

		prefs, err := studio.ParsePreferencesFile(input)
		if err != nil {
			return fmt.Errorf("parse %s: %w", input, err)
		}

		diff, err := studio.DiffRegistry(prefs)
		if err != nil {
			return fmt.Errorf("compute diff: %w", err)
		}

		if len(diff) == 0 {
			fmt.Println("Nothing to restore — all settings are already up to date.")
			return nil
		}

		if dryRun {
			return runDryRun(diff)
		}

		if !force {
			if err := printDiffTable(diff); err != nil {
				return err
			}

			fmt.Println()

			confirmed, err := prompt.Confirm(
				"Restore NetLinx Studio settings?",
				fmt.Sprintf("%d settings differ from:\n%s\n\nThis will overwrite those registry values.", len(diff), input),
			)
			if err != nil {
				return fmt.Errorf("confirmation prompt: %w", err)
			}

			if !confirmed {
				fmt.Println("Aborted.")
				return nil
			}
		}

		if doBackup {
			backupPath := time.Now().Format("netlinx-studio-backup-2006-01-02-150405") + ".epx"
			if err := writeBackup(backupPath); err != nil {
				return fmt.Errorf("write pre-restore backup: %w", err)
			}

			fmt.Printf("Current settings backed up to: %s\n", backupPath)
		}

		fmt.Println("Writing settings to registry...")

		if err := studio.ApplyRegistryDiff(diff); err != nil {
			return fmt.Errorf("write registry: %w", err)
		}

		fmt.Printf("Registry updated: %d values written.\n", len(diff))
		return nil
	},
}

// writeBackup reads the current registry state and writes it to path.
func writeBackup(path string) error {
	settings, err := studio.ReadRegistry()
	if err != nil {
		return fmt.Errorf("read registry: %w", err)
	}

	prefs := studio.BuildPreferences(settings)

	xmlBytes, err := xml.MarshalIndent(prefs, "", "    ")
	if err != nil {
		return fmt.Errorf("marshal preferences: %w", err)
	}

	content := append([]byte(xml.Header), xmlBytes...)
	content = append(content, '\n')

	return os.WriteFile(path, content, 0o644)
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

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(r),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{Alignment: tw.AlignCenter},
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
