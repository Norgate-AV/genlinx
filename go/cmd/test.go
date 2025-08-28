package cmd

import (
	"fmt"

	"github.com/Norgate-AV/genlinx-go/internal/utils"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test-paths",
	Short: "Test path normalization functionality",
	Long:  `Test command to demonstrate path normalization across different operating systems`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🧪 Testing Path Normalization")
		fmt.Println("=============================")

		// Test various path formats
		testPaths := []string{
			"C:/Program Files (x86)/Common Files/AMXShare/AXIs",
			"C:\\Program Files (x86)\\Common Files\\AMXShare\\AXIs",
			"/usr/local/share/amx/axis",
			"./relative/path",
			"../parent/path",
			"simple-file.txt",
		}

		fmt.Println("Original Paths → Normalized Paths:")
		fmt.Println("----------------------------------")
		for _, path := range testPaths {
			normalized := utils.NormalizePath(path)
			fmt.Printf("  %-50s → %s\n", path, normalized)
		}

		fmt.Println("\n✅ Path normalization test completed!")
		return nil
	},
}

func init() {
	// Add test command to root
	rootCmd.AddCommand(testCmd)
}
