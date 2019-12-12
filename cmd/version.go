package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type versionInfo struct {
	version string
	commit  string
	date    string
}

var (
	verInfo *versionInfo
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s version %s, build %s\n", cliName, verInfo.version, verInfo.commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
