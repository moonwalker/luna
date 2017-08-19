package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/moonwalker/luna/support"
)

var svcRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run services specified in config",

	Run: func(cmd *cobra.Command, args []string) {
		var cfg support.Config

		err := viper.Unmarshal(&cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		pm := support.NewPM(cfg)
		pm.Run()
	},
}

func init() {
	servicesCmd.AddCommand(svcRunCmd)
}
