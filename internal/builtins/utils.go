package builtins

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"go.starlark.net/starlark"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	shsyntax "mvdan.cc/sh/v3/syntax"

	"github.com/moonwalker/luna/internal/support"
)

// https://github.com/mvdan/sh/blob/master/interp/example_test.go
// https://github.com/go-task/task/blob/main/internal/execext/exec.go#L35
func run(src string, dir string, env []string) error {
	environ := support.Environ("", env...)

	params := []string{"-e"}
	params = append(params, support.TaskParams...)

	r, err := interp.New(
		interp.Params(params...),
		interp.Dir(dir),
		interp.Env(expand.ListEnviron(environ...)),
		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
	)
	if err != nil {
		return err
	}

	p, err := shsyntax.NewParser().Parse(strings.NewReader(src), "")
	if err != nil {
		return err
	}

	return r.Run(context.Background(), p)
}

func stringArray(in starlark.Value) (res []string) {
	if in == nil {
		return
	}
	switch v := in.(type) {
	case starlark.String:
		s := in.(starlark.String)
		res = append(res, string(s))
	case *starlark.List:
		l := in.(*starlark.List)
		for i := 0; i < l.Len(); i++ {
			s, ok := l.Index(i).(starlark.String)
			if ok {
				res = append(res, string(s))
			}
		}
	default:
		fmt.Printf("unknown type %T\n", v)
	}
	return
}

func stringAndBool(in starlark.Value) (res string) {
	if in == nil {
		return
	}
	switch v := in.(type) {
	case starlark.String:
		s := in.(starlark.String)
		res = string(s)
	case starlark.Bool:
		b := in.(starlark.Bool)
		res = strconv.FormatBool(bool(b))
	default:
		fmt.Printf("unknown type %T\n", v)
	}
	return
}
