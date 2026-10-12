package find

import (
	"os"

	"github.com/fatih/color"
	"github.com/neilotoole/jsoncolor"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

// PrintJSON outputs devices as colorized pretty-printed JSON.
func PrintJSON(devices []Device) error {
	enc := jsoncolor.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	if jsoncolor.IsColorTerminal(os.Stdout) {
		enc.SetColors(jsoncolor.DefaultColors())
	}

	return enc.Encode(devices)
}

// PrintTable outputs devices as a coloured table with green borders and a title.
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

	table.Caption(tw.Caption{Text: "Discovered Devices", Spot: tw.SpotTopCenter})
	table.Header("IP Address", "System", "Date", "Time", "MAC Address", "Hostname", "ID")

	for _, d := range devices {
		_ = table.Append(d.IP, d.System, d.Date.Text, d.Time, d.MAC, d.Hostname, d.ID)
	}

	_ = table.Render()
}
