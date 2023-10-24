package yamlconfig

import (
	"os"

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
		for name, svc := range cfg.Services {
			table.Append([]string{name, svc.Dir})
		}
		table.Render()
	},
}

func init() {
	servicesCmd.AddCommand(svcListCmd)
}
