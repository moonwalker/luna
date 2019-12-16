package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	servicesCmd.AddCommand(svcPackCmd)
}

var svcPackCmd = &cobra.Command{
	Use:   "build",
	Short: "Build services",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("build...")
	},
}
