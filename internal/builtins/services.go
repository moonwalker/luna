package builtins

import (
	"go.starlark.net/starlark"

	"github.com/moonwalker/luna/internal/support"
)

func service(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	svc, err := serviceFromArgs(b, args, kwargs)
	if err != nil {
		return nil, err
	}

	err = support.RegisterService(svc)
	if err != nil {
		return nil, err
	}

	return starlark.None, nil
}

func go_service(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	svc, err := serviceFromArgs(b, args, kwargs)
	if err != nil {
		return nil, err
	}

	// enforce go attributes
	err = svc.SetKind(support.GoService)
	if err != nil {
		return nil, err
	}

	err = support.RegisterService(svc)
	if err != nil {
		return nil, err
	}

	return starlark.None, nil
}

func serviceFromArgs(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (*support.Service, error) {
	var name, dir, run, cmd, bin string
	var dep starlark.Value
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

	svc := &support.Service{
		Kind:  support.GenericService,
		Name:  name,
		Dir:   dir,
		Run:   run,
		Cmd:   cmd,
		Bin:   bin,
		Dep:   stringArray(dep),
		Watch: watch,
	}

	return svc, nil
}
