package builtins

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

//

func init() {
	starlark.Universe["http"] = httpModule
}

var httpModule = &starlarkstruct.Module{
	Name: "json",
	Members: starlark.StringDict{
		"get": starlark.NewBuiltin("http.get", http_get),
	},
}

func http_get(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var url string
	err := starlark.UnpackArgs(b.Name(), args, kwargs,
		"url", &url,
	)
	if err != nil {
		return nil, err
	}

	fmt.Println(url)

	return starlark.None, nil
}
