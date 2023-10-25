package services

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/luna/internal/pm"
	"github.com/moonwalker/luna/internal/support"
)

var (
	svcStopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop services",

		Run: func(cmd *cobra.Command, args []string) {
			pm := pm.NewPM(support.Services(), args)
			pm.Stop()
		},
	}
)

func init() {
	servicesCmd.AddCommand(svcStopCmd)
}