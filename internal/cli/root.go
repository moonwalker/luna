package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/moonwalker/luna/internal/cli/services"
	"github.com/moonwalker/luna/internal/support"
	"github.com/moonwalker/luna/internal/tasks"
)

const (
	name = "luna"
	desc = "Command line tool for microservices in monorepos"

	defaultFile = "Lunafile"
	defaultYaml = "luna.yaml"
)

var (
	file    string
	nofile  bool
	rootCmd = &cobra.Command{
		Use:               name,
		Short:             desc,
		SilenceErrors:     true,
		SilenceUsage:      true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
)

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})
	rootCmd.PersistentFlags().StringVarP(&file, "filename", "f", "", fmt.Sprintf("file to execute (default: %s or %s)", defaultFile, defaultYaml))
	rootCmd.ParseFlags(os.Args)
}

func Run(version, commit, date string) error {
	rootCmd.Version = fmt.Sprintf("%s, build %s", version, commit)

	// no file specified
	nofile = len(file) == 0

	// no file specified, try defaults
	if nofile {
		// lunafile
		if support.FileExists(defaultFile) {
			if err := tasks.Load(defaultFile, rootCmd); err != nil {
				return err
			}
		} else {
			// lunayaml
			if support.FileExists(defaultYaml) {
				if err := support.LoadYaml(defaultYaml); err != nil {
					return err
				}
			}
		}
	}

	// file specified
	if !nofile {
		// but not exists
		if !support.FileExists(file) {
			rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
				return fmt.Errorf("%s not found", file)
			}
		} else {
			// file exists
			// check if yaml
			if support.IsYaml(file) {
				if err := support.LoadYaml(file); err != nil {
					return err
				}
			} else {
				if err := tasks.Load(file, rootCmd); err != nil {
					return err
				}
			}
		}
	}

	if len(support.Services()) > 0 {
		services.EnableServicesCmd(rootCmd)
	}

	if len(rootCmd.Commands()) == 0 {
		rootCmd.Run = func(cmd *cobra.Command, args []string) {}
	}

	return rootCmd.Execute()
}
