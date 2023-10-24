package tasks

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"

	"github.com/moonwalker/luna/internal/builtins"
)

const (
	groupID    = "tasks"
	groupTitle = "Tasks:"
)

var (
	thread = &starlark.Thread{
		Name:  "main",
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}
)

func Load(fname string, command *cobra.Command) error {
	builtins.ConfigTasks(command, groupID, groupTitle)

	globals, err := execFile(fname)
	if err != nil {
		return err
	}

	if len(globals.Keys()) == 0 {
		// command.Run = func(cmd *cobra.Command, args []string) {}
		return nil
	}

	fns := make([]*starlark.Function, 0)
	for _, name := range globals.Keys() {
		v := globals[name]
		// fmt.Println(name, v.Type())
		fn, ok := (v).(*starlark.Function)
		if ok {
			// convention for private, skip those
			if !strings.HasPrefix(fn.Name(), "_") {
				fns = append(fns, fn)
			}
		}
	}

	if len(fns) > 0 {
		if !command.ContainsGroup(groupID) {
			command.AddGroup(&cobra.Group{ID: groupID, Title: groupTitle})
		}

		for _, fn := range fns {
			addCommand(command, fn)
		}
	}

	return nil
}

func execFile(name string) (starlark.StringDict, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	globals, err := starlark.ExecFile(thread, name, string(data), builtins.Predeclared)
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			return nil, errors.New(evalErr.Backtrace())
		}
		return nil, err
	}

	return globals, nil
}

func addCommand(command *cobra.Command, fn *starlark.Function) {
	minParams, params := parseParams(fn)
	usage := fmtUsage(fn.Name(), params)

	command.AddCommand(&cobra.Command{
		Use:     usage,
		Args:    cobra.MinimumNArgs(minParams),
		GroupID: groupID,
		Run: func(cmd *cobra.Command, args []string) {
			out, err := starlark.Call(thread, fn, starlarkArgs(args), nil)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			_, none := (out).(starlark.NoneType)
			if out != nil && !none {
				fmt.Println(strings.Trim(out.String(), `"`))
			}
		},
	})
}

func parseParams(fn *starlark.Function) (int, map[string]string) {
	params := make(map[string]string)
	minParams := fn.NumParams()

	for i := 0; i < fn.NumParams(); i++ {
		param, _ := fn.Param(i)
		defval := fn.ParamDefault(i)

		if defval != nil {
			params[param] = defval.String()
			minParams--
		} else {
			params[param] = ""
		}
	}

	return minParams, params
}

func fmtUsage(fnName string, params map[string]string) string {
	usage := fnName

	for name, defval := range params {
		if defval == "" {
			usage += fmt.Sprintf(" <%s>", name)
		} else {
			usage += fmt.Sprintf(" [%s=%s]", name, defval)
		}
	}

	return usage
}

func starlarkArgs(args []string) starlark.Tuple {
	sargs := starlark.Tuple{}
	for _, v := range args {
		sargs = append(sargs, starlark.String(v))
	}
	return sargs
}
