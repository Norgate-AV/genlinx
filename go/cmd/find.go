package cmd

import (
	"fmt"
	"time"

	"github.com/Norgate-AV/genlinx-go/internal/find"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "find NetLinx devices on a local broadcast subnet",
	Long:  `Find NetLinx devices on a local broadcast subnet`,
	RunE:  runFind,
}

func runFind(cmd *cobra.Command, args []string) error {
	timeoutMs, _ := cmd.Flags().GetInt("timeout")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	fmt.Println("Listening for NetLinx devices...")

	devices, err := find.Discover(time.Duration(timeoutMs) * time.Millisecond)
	if err != nil {
		return err
	}

	if len(devices) == 0 {
		fmt.Println("No devices found")
		return nil
	}

	if jsonOutput {
		return find.PrintJSON(devices)
	}

	find.PrintTable(devices)

	return nil
}

func init() {
	findCmd.Flags().IntP("timeout", "t", 6000, "timeout in milliseconds")
	findCmd.Flags().BoolP("json", "j", false, "output as JSON")
}
