package builtins

import (
	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
)

type BuiltinsConfig struct {
	RootCommand     *cobra.Command
	TasksGroupID    string
	TasksGroupTitle string
}

var (
	Config      = &BuiltinsConfig{}
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
