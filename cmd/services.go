package cmd

import (
	"github.com/spf13/cobra"
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Work with services",
}

func init() {
	rootCmd.AddCommand(servicesCmd)
}
