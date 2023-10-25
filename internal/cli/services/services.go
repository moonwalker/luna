package services

import (
	"github.com/spf13/cobra"
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Working with services",
}

func EnableServicesCmd(rootCmd *cobra.Command) {
	rootCmd.AddCommand(servicesCmd)
}
