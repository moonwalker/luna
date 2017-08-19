package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var svcListCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Luna",

	Run: func(cmd *cobra.Command, args []string) {
		// TODO
		fmt.Println("v0.1.0")
	},
}

func init() {
	RootCmd.AddCommand(svcListCmd)
}
