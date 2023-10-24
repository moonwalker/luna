package builtins

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
)

func task(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name string
	var dir, cmds starlark.Value

	err := starlark.UnpackArgs(b.Name(), args, kwargs,
		"name", &name,
		"cmds?", &cmds,
		"dir?", &dir,
	)
	if err != nil {
		return nil, err
	}

	if !Config.RootCommand.ContainsGroup(Config.TasksGroupID) {
		Config.RootCommand.AddGroup(&cobra.Group{ID: Config.TasksGroupID, Title: Config.TasksGroupTitle})
	}

	Config.RootCommand.AddCommand(&cobra.Command{
		Use:     name,
		GroupID: Config.TasksGroupID,
		Run: func(c *cobra.Command, a []string) {
			src := strings.Join(stringArray(cmds), "\n")
			err = run(src, nil)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		},
	})

	return starlark.None, nil
}
