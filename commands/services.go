package commands

import (
	"github.com/spf13/cobra"
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Work with services",
}

func init() {
	RootCmd.AddCommand(servicesCmd)
}
