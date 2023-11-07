package services

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/luna/internal/config"
	"github.com/moonwalker/luna/internal/pm"
	"github.com/moonwalker/luna/internal/support"
)

var (
	dryRun    bool
	svcRunCmd = &cobra.Command{
		Use:   "run",
		Short: "Run services",
		Run: func(cmd *cobra.Command, args []string) {
			pm := pm.NewPM(config.GetConfig(), support.Services(), args)
			if dryRun {
				servicesListTable(pm.Runnables)
			} else {
				pm.Run()
			}
		},
	}
)

func init() {
	svcRunCmd.PersistentFlags().BoolVarP(&dryRun, "dryrun", "d", false, "simply display the list of services to run")
	servicesCmd.AddCommand(svcRunCmd)
}
