package builtins

import (
	"go.starlark.net/starlark"
)

func sh(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var cmd, dir string
	var senv starlark.Value
	err := starlark.UnpackArgs(b.Name(), args, kwargs,
		"cmd", &cmd,
		"dir?", &dir,
		"env?", &senv,
	)
	if err != nil {
		return nil, err
	}

	env := stringArray(senv)
	err = run(cmd, dir, env)
	if err != nil {
		return nil, err
	}

	return starlark.None, nil
}
