package cmd

import (
	"encoding/xml"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"

	"github.com/Norgate-AV/genlinx/internal/studio"
)

var studioCmd = &cobra.Command{
	Use:   "studio",
	Short: "NetLinx Studio utilities",
	Long:  `Tools for interacting with NetLinx Studio settings and configuration.`,
}

var studioBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup NetLinx Studio settings to an .epx file",
	Long: `Read all NetLinx Studio settings from the Windows registry and write a
compatible .epx preferences file. The output file can be re-imported directly
using NetLinx Studio's native File > Import Preferences option.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if runtime.GOOS != "windows" {
			return fmt.Errorf("the studio command is only supported on Windows")
		}

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
	studioBackupCmd.Flags().StringP("output", "o", "", "output file path (default: <timestamp>.epx in current directory)")
	studioCmd.AddCommand(studioBackupCmd)
}
