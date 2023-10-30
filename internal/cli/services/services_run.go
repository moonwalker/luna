package services

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/luna/internal/config"
	"github.com/moonwalker/luna/internal/pm"
	"github.com/moonwalker/luna/internal/support"
)

var (
	svcRunCmd = &cobra.Command{
		Use:   "run",
		Short: "Run services",

		Run: func(cmd *cobra.Command, args []string) {
			pm := pm.NewPM(config.GetConfig(), support.Services(), args)
			pm.Run()
		},
	}
)

func init() {
	servicesCmd.AddCommand(svcRunCmd)
}
