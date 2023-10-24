package starlark

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/moonwalker/luna/internal/support"
)

var (
	servicesCmd = &cobra.Command{
		Use:   "services",
		Short: "Working with services",
	}

	svcRunCmd = &cobra.Command{
		Use:   "run",
		Short: "Run services",
		Run: func(cmd *cobra.Command, args []string) {
			runServices(args)
		},
	}
)

func UseStarlark(rootCmd *cobra.Command) {
	rootCmd.AddCommand(servicesCmd)
}

func init() {
	servicesCmd.AddCommand(svcRunCmd)
}

func runServices(serviceNames []string) {
	services := support.FindServices(serviceNames...)

	for _, s := range services {
		fmt.Println(s)
	}
}
