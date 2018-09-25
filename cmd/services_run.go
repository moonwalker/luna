package cmd

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/luna/support"
)

var (
	detach bool

	svcRunCmd = &cobra.Command{
		Use:   "run",
		Short: "Run services",

		Run: func(cmd *cobra.Command, args []string) {
			pm := support.NewPM(cfg)
			pm.Run(args, detach)
		},
	}
)

func init() {
	svcRunCmd.Flags().BoolVarP(&detach, "detach", "d", false, "run services in the background")
	servicesCmd.AddCommand(svcRunCmd)
}
