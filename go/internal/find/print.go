package find

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

// PrintJSON outputs devices as pretty-printed JSON.
func PrintJSON(devices []Device) error {
	b, err := json.MarshalIndent(devices, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

// PrintTable outputs devices as a coloured table with green borders.
func PrintTable(devices []Device) {
	green := renderer.Tint{FG: renderer.Colors{color.FgGreen}}

	r := renderer.NewColorized(renderer.ColorizedConfig{
		Border:    green,
		Separator: green,
	})

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(r),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{Alignment: tw.AlignCenter},
			},
		}),
	)

	table.Header("IP Address", "System", "Date", "Time", "MAC Address", "Hostname", "ID")

	for _, d := range devices {
		table.Append(d.IP, d.System, d.Date.Text, d.Time, d.MAC, d.Hostname, d.ID)
	}

	table.Render()
}
