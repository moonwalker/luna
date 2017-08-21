package commands

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/moonwalker/luna/support"
)

var (
	listCompareRange string
)

var svcListCmd = &cobra.Command{
	Use:   "list",
	Short: "List services",

	Run: func(cmd *cobra.Command, args []string) {
		cfg.PopulateChanges(listCompareRange)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Directory", "Changed"})
		for name, svc := range cfg.Services {
			table.Append([]string{name, svc.Chdir, support.BoolTostring(svc.Changed, "✔", "")})
		}
		table.Render()
	},
}

func init() {
	svcListCmd.Flags().StringVarP(&listCompareRange, "compare", "c", "", "Git compare range or URL")
	servicesCmd.AddCommand(svcListCmd)
}
