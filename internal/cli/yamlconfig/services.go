package yamlconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/moonwalker/luna/internal/support"
)

var (
	cfgFile string
	cfg     support.Config
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Working with services",
}

func UseYamlConfig(rootCmd *cobra.Command, file string) {
	cfgFile = file
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(servicesCmd)
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	// enable ability to specify config file via flag
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	cfgName := strings.TrimSuffix(cfgFile, filepath.Ext(cfgFile))

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
