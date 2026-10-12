//go:build !windows

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func studioBackupRunE(_ *cobra.Command, _ []string) error {
	return fmt.Errorf("the studio command is only supported on Windows")
}

func studioRestoreRunE(_ *cobra.Command, _ []string) error {
	return fmt.Errorf("the studio command is only supported on Windows")
}

func writeBackup(_ string) error {
	return fmt.Errorf("the studio command is only supported on Windows")
}
