package cmd

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config [key]",
	Short: "view/edit configuration properties for genlinx",
	Long:  `View/edit configuration properties for genlinx`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement config logic
		return nil
	},
}

func init() {
	configCmd.Flags().Bool("global", false, "use global configuration")
	configCmd.Flags().Bool("local", false, "use local configuration")
	configCmd.Flags().BoolP("list", "l", false, "display the configuration in stdout")
	configCmd.Flags().BoolP("edit", "e", false, "edit the configuration with default text editor")
}
