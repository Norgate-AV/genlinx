package cmd

import (
	"github.com/Norgate-AV/genlinx-go/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "genlinx",
	Short: "CLI utility for NetLinx projects 🚀🚀🚀",
	Long: `                   _ _
   __ _  ___ _ __ | (_)_ __ __  __
  / _` + "`" + ` |/ _ \ '_ \| | | '_ / \/ /
 | (_| |  __/ | | | | | | | |>  <
  \__, |\___|_| |_|_|_|_| |_/_/\_\
  |___/

Open source CLI tool for NetLinx projects
Copyright (c) 2025, Norgate AV
https://github.com/Norgate-AV/genlinx`,
	Version: version.Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(cfgCmd)
	rootCmd.AddCommand(archiveCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(findCmd)

	// Global flags
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")
}
