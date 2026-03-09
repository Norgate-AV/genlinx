package cmd

import (
	"strconv"
	"time"

	"github.com/damienbutt/figlet"
	"github.com/spf13/cobra"

	"github.com/Norgate-AV/genlinx/internal/version"
)

const AppName = "genlinx"

var rootCmd = &cobra.Command{
	Use:           AppName,
	Short:         "CLI utility for NetLinx projects 🚀🚀🚀",
	Long:          getBanner(),
	Version:       version.Version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func getBanner() string {
	year := strconv.Itoa(time.Now().Year())

	banner, err := figlet.Text(AppName)
	if err != nil {
		banner = AppName
	}

	banner += `

` + version.Version + `
Open source CLI tool for NetLinx projects
Copyright (c) 2010-` + year + `, Norgate AV
https://github.com/Norgate-AV/genlinx

===================================================`

	return banner
}

func getFooter() string {
	return `
===================================================

For more help, make sure to check out the man page:
    $ man ` + AppName
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

	// Append afterAll footer to the default help output.
	rootCmd.SetHelpTemplate(`{{with .Long}}{{. | trimRightSpace}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}` + getFooter() + "\n")
}
