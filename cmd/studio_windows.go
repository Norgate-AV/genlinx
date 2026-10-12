//go:build windows

package cmd

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/Norgate-AV/genlinx/internal/prompt"
	"github.com/Norgate-AV/genlinx/internal/studio"
)

func studioBackupRunE(cmd *cobra.Command, _ []string) error {
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
}

func studioRestoreRunE(cmd *cobra.Command, args []string) error {
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
}

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
