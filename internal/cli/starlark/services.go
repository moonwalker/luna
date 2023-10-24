package starlark

import (
	"github.com/spf13/cobra"
)

var (
	servicesCmd = &cobra.Command{
		Use:   "services",
		Short: "Working with services",
	}
)

func UseStarlark(rootCmd *cobra.Command) {
	rootCmd.AddCommand(servicesCmd)
}

func init() {
	servicesCmd.AddCommand(svcRunCmd)
}
