package builtins

import (
	"fmt"

	"go.starlark.net/starlark"

	"github.com/moonwalker/luna/internal/support"
)

func service(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name, dir, run, cmd, bin string
	var dep []string
	var watch bool

	err := starlark.UnpackArgs(b.Name(), args, kwargs,
		"name", &name,
		"dir", &dir,
		"run?", &run,
		"cmd?", &cmd,
		"bin?", &bin,
		"dep?", &dep,
		"watch?", &watch,
	)
	if err != nil {
		return nil, err
	}

	err = support.RegisterService(&support.Service{
		Kind:  support.GenericService,
		Name:  name,
		Dir:   dir,
		Run:   run,
		Cmd:   cmd,
		Bin:   bin,
		Dep:   dep,
		Watch: watch,
	})
	if err != nil {
		return nil, err
	}

	return starlark.None, nil
}

func go_service(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name, dir string
	var dep []string
	var watch bool

	err := starlark.UnpackArgs(b.Name(), args, kwargs,
		"name", &name,
		"dir", &dir,
		"dep?", &dep,
		"watch?", &watch,
	)
	if err != nil {
		return nil, err
	}

	run := "go run ."
	cmd := fmt.Sprintf("go build -o ./tmp/%s ./%s", name, dir)
	bin := fmt.Sprintf("./tmp/%s", name)

	err = support.RegisterService(&support.Service{
		Kind:  support.GoService,
		Name:  name,
		Dir:   dir,
		Run:   run,
		Cmd:   cmd,
		Bin:   bin,
		Dep:   dep,
		Watch: watch,
	})
	if err != nil {
		return nil, err
	}

	return starlark.None, nil
}
