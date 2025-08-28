package cmd

import (
	"github.com/spf13/cobra"
)

var cfgCmd = &cobra.Command{
	Use:   "cfg",
	Short: "generate NetLinx build CFG files",
	Long:  `Generate NetLinx build configuration files`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement cfg logic
		return nil
	},
}

func init() {
	cfgCmd.Flags().StringSliceP("workspace-files", "w", []string{}, "workspace file(s) to generate a CFG for")
	cfgCmd.Flags().StringP("root-directory", "r", "", "root directory reference")
	cfgCmd.Flags().StringP("output-file-suffix", "o", "", "output file suffix")
	cfgCmd.Flags().StringP("output-log-file-suffix", "f", "", "output log file suffix")
	cfgCmd.Flags().StringP("output-log-file-option", "k", "N", "output log file option (A - append, N - overwrite)")
	cfgCmd.Flags().BoolP("output-log-console-option", "c", true, "output log to console")
	cfgCmd.Flags().BoolP("build-with-debug-information", "d", false, "build with debug information")
	cfgCmd.Flags().BoolP("build-with-source", "s", false, "build with source")
	cfgCmd.Flags().StringSliceP("include-path", "i", []string{}, "add additional include paths")
	cfgCmd.Flags().StringSliceP("module-path", "m", []string{}, "add additional module paths")
	cfgCmd.Flags().StringSliceP("library-path", "l", []string{}, "add additional library paths")
	cfgCmd.Flags().BoolP("all", "a", false, "select all workspace files without prompting")
}
