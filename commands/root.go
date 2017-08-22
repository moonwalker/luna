package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/moonwalker/luna/support"
)

const (
	cliName = "luna"
	cfgName = "luna"
	cfgExt  = "yml"
)

var (
	version string
	commit  string
	cfgFile string
	cfg     support.Config
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   cliName,
	Short: "Command line tool for microservices in monorepos",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v, c string) {
	version = v
	commit = c
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default: %s.%s)", cfgName, cfgExt))
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	// enable ability to specify config file via flag
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(cfgName)
	viper.AddConfigPath(".")

	// read in environment variables that match
	viper.AutomaticEnv()

	// if a config file is found, read it in
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file not found.")
		os.Exit(1)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		fmt.Println("Config file not valid.")
		os.Exit(1)
	}
}
