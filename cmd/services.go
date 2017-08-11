package cmd

import "github.com/spf13/cobra"

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	RootCmd.AddCommand(servicesCmd)
}
