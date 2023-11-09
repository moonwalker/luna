package builtins

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type BuiltinsConfig struct {
	RootCommand     *cobra.Command
	TasksGroupID    string
	TasksGroupTitle string
}

var (
	Config = &BuiltinsConfig{}

	Thread = &starlark.Thread{
		Name:  "main",
		Print: print,
		Load:  load,
	}

	Predeclared = starlark.StringDict{
		"env":         starlark.NewBuiltin("env", env),
		"sh":          starlark.NewBuiltin("sh", sh),
		"task":        starlark.NewBuiltin("task", task),
		"service":     starlark.NewBuiltin("service", service),
		"go_service":  starlark.NewBuiltin("go_service", go_service),
		"docker_repo": starlark.NewBuiltin("docker_repo", docker_repo),
		"build_flags": starlark.NewBuiltin("build_flags", build_flags),
	}
)

func ConfigTasks(command *cobra.Command, tasksGroupID string, tasksGroupTitle string) {
	Config.RootCommand = command
	Config.TasksGroupID = tasksGroupID
	Config.TasksGroupTitle = tasksGroupTitle
}

func print(_ *starlark.Thread, msg string) {
	fmt.Println(msg)
}

func load(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	if starlark.Universe.Has(module) {
		return starlark.StringDict{module: starlark.Universe[module]}, nil
	}

	data, err := os.ReadFile(module)
	if err != nil {
		return nil, err
	}

	return starlark.ExecFile(thread, module, data, starlark.StringDict{
		"module": starlark.NewBuiltin("module", starlarkstruct.MakeModule),
	})
}
