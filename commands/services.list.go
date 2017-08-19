package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "list",
	Short: "Print list of services",

	Run: func(cmd *cobra.Command, args []string) {
		// TODO
		fmt.Println("services...")
	},
}

func init() {
	servicesCmd.AddCommand(versionCmd)
}
