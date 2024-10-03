package builtins

import (
	"os"

	"go.starlark.net/starlark"
)

func env(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var key string
	var val starlark.Value
	err := starlark.UnpackArgs(b.Name(), args, kwargs,
		"key", &key,
		"val", &val,
	)
	if err != nil {
		return nil, err
	}

	_, exists := os.LookupEnv(key)
	if !exists {
		err = os.Setenv(key, stringAndBool(val))
		if err != nil {
			return nil, err
		}
	}

	return starlark.None, nil
}
