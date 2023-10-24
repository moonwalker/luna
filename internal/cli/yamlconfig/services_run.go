package yamlconfig

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/luna/internal/support"
)

var (
	svcRunCmd = &cobra.Command{
		Use:   "run",
		Short: "Run services",

		Run: func(cmd *cobra.Command, args []string) {
			pm := support.NewPM(cfg, args)
			pm.Run()
		},
	}
)

func init() {
	servicesCmd.AddCommand(svcRunCmd)
}
