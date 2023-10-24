package starlark

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/moonwalker/luna/internal/support"
)

var (
	svcRunCmd = &cobra.Command{
		Use:   "run",
		Short: "Run services",

		Run: func(cmd *cobra.Command, args []string) {
			runServices(args)
		},
	}
)

func init() {
	servicesCmd.AddCommand(svcRunCmd)
}
func runServices(serviceNames []string) {
	services := []*support.Service{}

	if len(serviceNames) == 0 {
		services = support.AllServices()
	} else {
		services = support.FindServices(serviceNames...)
	}

	for _, svc := range services {
		fmt.Println(svc.Name)
	}
}
