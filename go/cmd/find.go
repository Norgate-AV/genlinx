package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/Norgate-AV/genlinx-go/internal/find"
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
	watch, _ := cmd.Flags().GetBool("watch")

	var devices []find.Device
	var err error

	if watch {
		fmt.Fprintln(os.Stderr, "Listening for NetLinx devices... (press Ctrl+C to stop)")

		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		devices, err = find.DiscoverWithContext(ctx)
	} else {
		fmt.Fprintln(os.Stderr, "Listening for NetLinx devices...")
		devices, err = find.Discover(time.Duration(timeoutMs) * time.Millisecond)
	}

	if err != nil {
		return err
	}

	if len(devices) == 0 {
		fmt.Fprintln(os.Stderr, "No devices found")
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
	findCmd.Flags().BoolP("watch", "w", false, "listen until Ctrl+C, then output results")
}
