package cmd

import (
	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "generate a NetLinx workspace zip archive",
	Long:  `Generate a NetLinx workspace zip archive`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement archive logic
		return nil
	},
}

func init() {
	archiveCmd.Flags().StringSliceP("workspace-files", "w", []string{}, "workspace file(s) to generate archive(s) for")
	archiveCmd.Flags().StringP("output-file-suffix", "o", "", "output file suffix")
	archiveCmd.Flags().BoolP("include-compiled-source-files", "S", true, "include compiled source files")
	archiveCmd.Flags().BoolP("include-compiled-module-files", "M", true, "include compiled module files")
	archiveCmd.Flags().BoolP("include-files-not-in-workspace", "N", true, "include files not in workspace")
	archiveCmd.Flags().StringSliceP("extra-file-search-locations", "l", []string{}, "extra file locations to search")
	archiveCmd.Flags().StringP("extra-file-archive-location", "p", ".genlinx", "location to place extra files in the archive")
}
