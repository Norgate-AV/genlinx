package cmd

import (
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "find NetLinx devices on a local broadcast subnet",
	Long:  `Find NetLinx devices on a local broadcast subnet`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement find logic
		return nil
	},
}

func init() {
	findCmd.Flags().BoolP("json", "j", false, "output as JSON")
}
