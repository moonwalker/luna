package commands

import (
	"github.com/spf13/cobra"
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Build and run services",
}

func init() {
	RootCmd.AddCommand(servicesCmd)
}
