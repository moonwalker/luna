package services

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/moonwalker/luna/internal/support"
)

var svcListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List services",
	Aliases: []string{"ls"},

	Run: func(cmd *cobra.Command, args []string) {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Directory"})
		for _, svc := range support.Services() {
			table.Append([]string{svc.Name, svc.Dir})
		}
		table.Render()
	},
}

func init() {
	servicesCmd.AddCommand(svcListCmd)
}
