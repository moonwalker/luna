package cli

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/moonwalker/luna/internal/cli/starlark"
	"github.com/moonwalker/luna/internal/cli/yamlconfig"
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

	// try default files
	if len(file) == 0 {
		if fileExists(defaultFile) {
			file = defaultFile
		} else if fileExists(defaultYaml) {
			file = defaultYaml
		} else {
			nofile = true
		}
	}

	if nofile {
		rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
			if len(file) == 0 {
				return nil
			}
			if nofile {
				return fmt.Errorf("%s not found", file)
			}
			return nil
		}
	} else {
		if isYaml(file) {
			yamlconfig.UseYamlConfig(rootCmd, file)
		} else {
			if err := tasks.Load(file, rootCmd); err != nil {
				return err
			}
		}
	}

	if len(support.AllServices()) > 0 {
		starlark.UseStarlark(rootCmd)
	}

	if len(rootCmd.Commands()) == 0 {
		rootCmd.Run = func(cmd *cobra.Command, args []string) {}
	}

	return rootCmd.Execute()
}

func fileExists(f string) bool {
	if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func isYaml(f string) bool {
	ext := path.Ext(f)
	return ext == ".yml" || ext == ".yaml" // TODO: check with yaml parser
}
