package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const cliName = "luna"

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   cliName,
	Short: "Command line tool for microservices in monorepos",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version, commit, date string) {
	verInfo = &versionInfo{version, commit, date}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
