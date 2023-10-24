package builtins

import (
	"os"

	"github.com/joho/godotenv"
	"go.starlark.net/starlark"
)

func init() {
	// .env (default)
	godotenv.Load()

	// .env.local # local user specific (usually git ignored)
	godotenv.Overload(".env.local")
}

func env(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var key, val string
	err := starlark.UnpackArgs(b.Name(), args, kwargs,
		"key", &key,
		"val", &val,
	)
	if err != nil {
		return nil, err
	}

	_, exists := os.LookupEnv(key)
	if !exists {
		err = os.Setenv(key, val)
		if err != nil {
			return nil, err
		}
	}

	return starlark.None, nil
}
