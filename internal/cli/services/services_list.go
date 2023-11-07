package services

import (
	"fmt"
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
		servicesListTable(support.ServicesSorted())
	},
}

func init() {
	servicesCmd.AddCommand(svcListCmd)
}

func servicesListTable(svcList []*support.Service) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(append([]string{"#"}, support.ServicesKeys()...))
	for i, svc := range svcList {
		n := fmt.Sprintf("%d", i+1)
		table.Append(append([]string{n}, svc.Fields()...))
	}
	table.Render()
}
