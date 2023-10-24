package builtins

import (
	"go.starlark.net/starlark"

	"github.com/moonwalker/luna/internal/support"
)

func go_service(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name, dir string
	err := starlark.UnpackArgs(b.Name(), args, kwargs,
		"name", &name,
		"dir", &dir,
	)
	if err != nil {
		return nil, err
	}

	support.RegisterService(name, dir, support.GoService)

	return starlark.None, nil
}
