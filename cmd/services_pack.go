package cmd

import (
	"github.com/moonwalker/luna/support"
	"github.com/spf13/cobra"
)

var (
	packCompareRange string
)

var svcPackCmd = &cobra.Command{
	Use:   "pack",
	Short: "Package services",

	Run: func(cmd *cobra.Command, args []string) {
		cfg.PopulateChanges(packCompareRange)
		support.PackageServices(cfg)
	},
}

func init() {
	svcPackCmd.Flags().StringVarP(&packCompareRange, "compare", "c", "", "Git compare range or URL")
	servicesCmd.AddCommand(svcPackCmd)
}
