package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"go.starlark.net/starlark"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	shsyntax "mvdan.cc/sh/v3/syntax"
)

/* DRAFT API

tags = skaffold_build(raw=True)
kustomize_setimage(tags)


task('terra', args=[fix_target('target'), 'command'], cmd=['echo foo'], transforms=['args.target', fix_target])

run: xxx terra ...


#
skaffold build --file-output=deploy/tags.json

tusk tags deploy/tags.json

task('tags', args=['file'], cmd=[
	'ehco $file',
	'ehco 2',
])

def tags(fname):
	fname = fix_target(fname)
	kustomize_setimage(tags)
	shell_exec('echo $fname')

*/

// https://pkg.go.dev/go.starlark.net/starlark#example-ExecFile
func main() {
	var data = `
# print(greeting + ", world")
# print(repeat("one"))
# print(repeat("mur", 2))
# squares = [x*x for x in range(10)]

# task('terra', dir='tmp', args=['target', 'command'], cmds=['echo foo'])

def fix_target(s):
	return s

task('terra', args=['target', 'command'], cmds=[
	'echo $target $command'
])

task('tags', args='fname', cmds=[
	'echo $fname'
])

def tags(fname, bar=1):
	fname = fix_target(fname)
	# kustomize_setimage(tags)
	# shell_exec('echo $fname')
	print('>>>', fname, bar)
`
	task := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var name string
		var dir, cmds, cmdargs starlark.Value

		err := starlark.UnpackArgs(b.Name(), args, kwargs,
			"name", &name,
			"dir?", &dir,
			"cmds?", &cmds,
			"args?", &cmdargs,
		)
		if err != nil {
			return nil, err
		}

		env := make(map[string]string)
		argsArr := stringArray(cmdargs)

		for i, v := range os.Args[1:] {
			if i < len(argsArr) {
				argname := argsArr[i]
				env[argname] = v
			}
		}

		err = run(strings.Join(stringArray(cmds), "\n"), env)
		if err != nil {
			return nil, err
		}

		return starlark.None, nil
	}

	// repeat(str, n=1) is a Go function called from Starlark.
	// It behaves like the 'string * int' operation.
	repeat := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var s string
		var n int = 1
		if err := starlark.UnpackArgs(b.Name(), args, kwargs, "s", &s, "n?", &n); err != nil {
			return nil, err
		}
		return starlark.String(strings.Repeat(s, n)), nil
	}

	// The Thread defines the behavior of the built-in 'print' function.
	thread := &starlark.Thread{
		Name:  "example",
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}

	// This dictionary defines the pre-declared environment.
	predeclared := starlark.StringDict{
		"greeting": starlark.String("hello"),
		"repeat":   starlark.NewBuiltin("repeat", repeat),

		"task": starlark.NewBuiltin("task", task),
	}

	// Execute a program.
	globals, err := starlark.ExecFile(thread, "apparent/filename.star", data, predeclared)
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			log.Fatal(evalErr.Backtrace())
		}
		log.Fatal(err)
	}

	// Print the global environment.
	fmt.Println("\nGlobals:")
	for _, name := range globals.Keys() {
		v := globals[name]
		fmt.Printf("%s (%s) = %s\n", name, v.Type(), v.String())
	}

	// f, err := syntax.Parse("apparent/filename.star", data, 0)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = resolve.File(f, predeclared.Has, starlark.Universe.Has)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var pos syntax.Position
	// if len(f.Stmts) > 0 {
	// 	pos = syntax.Start(f.Stmts[0])
	// } else {
	// 	pos = syntax.MakePosition(&f.Path, 1, 1)
	// }
	// module := f.Module.(*resolve.Module)

	// fmt.Println("-->", pos, module.Globals)

	// Call Starlark function from Go.
	tags := globals["tags"]

	tagsfn := (tags).(*starlark.Function)

	p, _ := tagsfn.Param(0)
	fmt.Println(tagsfn.NumParams(), p)

	v, err := starlark.Call(thread, tags, starlark.Tuple{starlark.String("a"), starlark.String("b")}, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tags() = %v\n", v)
}

// https://github.com/mvdan/sh/blob/master/interp/example_test.go
// https://github.com/go-task/task/blob/main/internal/execext/exec.go#L35
func run(src string, env map[string]string) error {
	for k, v := range env {
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}

	environ := os.Environ()

	r, err := interp.New(
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
