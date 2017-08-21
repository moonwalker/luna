package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s version %s, build %s", cliName, version, commit)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
