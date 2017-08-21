package commands

import (
	"os"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/moonwalker/luna/support"
)

const (
	svcDirRx = "services/[^/]*"
)

var compareURL string

var svcListCmd = &cobra.Command{
	Use:   "list",
	Short: "Print list of services",

	Run: func(cmd *cobra.Command, args []string) {
		gitDiff := support.GetGitDiff(compareURL)
		changedServices := getChangedServices(gitDiff)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Dir", "*"})

		for name, svc := range cfg.Services {
			changed := support.StringInSlice(svc.Chdir, changedServices)
			table.Append([]string{name, svc.Chdir, support.BoolTostring(changed, "✔")})
		}

		table.Render()
	},
}

func init() {
	svcListCmd.Flags().StringVarP(&compareURL, "compare", "c", "", "Git compare URL")
	servicesCmd.AddCommand(svcListCmd)
}

func getChangedServices(gitDiff string) []string {
	var services []string
	diffs := strings.Split(gitDiff, "\n")
	re := regexp.MustCompile(svcDirRx)
	for _, diff := range diffs {
		svc := re.FindString(diff)
		services = support.AppendUnique(services, svc)
	}
	return services
}
