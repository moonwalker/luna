package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// version should be updated manually on release
	Version = "0.1.0"
	// git commit will be overwritten during build
	GitCommit = "HEAD"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s version %s, build %s", cliName, Version, GitCommit)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
