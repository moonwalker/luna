package starlark

import (
	"os"

	"github.com/moonwalker/luna/internal/support"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var svcListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List services",
	Aliases: []string{"ls"},

	Run: func(cmd *cobra.Command, args []string) {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Directory"})
		for _, svc := range support.AllServices() {
			table.Append([]string{svc.Name, svc.Dir})
		}
		table.Render()
	},
}

func init() {
	servicesCmd.AddCommand(svcListCmd)
}
