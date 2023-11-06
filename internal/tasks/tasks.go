package tasks

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"

	"github.com/moonwalker/luna/internal/builtins"
	"github.com/moonwalker/luna/internal/support"
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

type fnParam struct {
	name  string
	value string
}

func Load(fname string, command *cobra.Command) error {
	builtins.ConfigTasks(command, groupID, groupTitle)

	globals, err := execFile(fname)
	if err != nil {
		return err
	}

	if len(globals.Keys()) == 0 {
		return nil
	}

	fns := make([]*starlark.Function, 0)
	for _, name := range globals.Keys() {
		v := globals[name]
		fn, ok := (v).(*starlark.Function)
		if ok {
			// convention for private, skip functions prefixed with underscore
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
	use := fmtParams(fn.Name(), params)

	command.AddCommand(&cobra.Command{
		Use:                use,
		Args:               cobra.MinimumNArgs(minParams),
		GroupID:            groupID,
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			// set params as positional params for later usage in shell
			support.TaskParams = args

			// set params as env vars for later usage in shell
			paramsWithArgsToEnv(params, args)

			// call the function
			out, err := starlark.Call(thread, fn, starlarkArgs(args), nil)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			// report output
			_, none := (out).(starlark.NoneType)
			if out != nil && !none {
				fmt.Println(strings.Trim(out.String(), `"`))
			}
		},
	})
}

func parseParams(fn *starlark.Function) (int, []*fnParam) {
	var params []*fnParam
	minParams := fn.NumParams()

	if fn.HasVarargs() {
		minParams--
	}
	if fn.HasKwargs() {
		minParams--
	}

	p := minParams
	for i := 0; i < p; i++ {
		name, _ := fn.Param(i)
		value := fn.ParamDefault(i)

		// has default value, less required params
		if value != nil {
			params = append(params, &fnParam{name, value.String()})
			minParams--
		} else {
			params = append(params, &fnParam{name, ""})
		}
	}

	return minParams, params
}

func fmtParams(s string, params []*fnParam) string {
	for _, p := range params {
		if p.value == "" {
			s += fmt.Sprintf(" <%s>", p.name)
		} else {
			s += fmt.Sprintf(" [%s=%s]", p.name, p.value)
		}
	}
	return s
}

func paramsWithArgsToEnv(params []*fnParam, args []string) {
	for i, p := range params {
		v := args[i]
		os.Setenv(p.name, v)
	}
}

func starlarkArgs(args []string) starlark.Tuple {
	sargs := starlark.Tuple{}
	for _, v := range args {
		sargs = append(sargs, starlark.String(v))
	}
	return sargs
}
