package pm

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cosmtrek/air/runner"
)

// air --build.cmd "go build -o ./tmp/service_a ./examples/service_a" --build.bin "./tmp/service_a" --build.exclude_dir "dist,tmp,vendor"
func air(cmd, bin, dir string) (*runner.Engine, error) {
	root, err := expandPath(".")
	if err != nil {
		return nil, err
	}

	cfg := &runner.Config{}
	cfg.Root = root
	cfg.TmpDir = "tmp"
	cfg.TestDataDir = "testdata"
	cfg.Build.Log = "build-errors.log"

	cfg.Build.ExcludeDir = []string{"assets", "dist", "tmp", "vendor", "testdata"}
	cfg.Build.IncludeExt = []string{"go", "tpl", "tmpl", "html"}
	cfg.Build.ExcludeRegex = []string{"_test.go"}

	cfg.Build.Delay = 1000
	cfg.Build.Rerun = false
	cfg.Build.RerunDelay = 500

	cfg.Build.Cmd = cmd
	cfg.Build.Bin = bin
	cfg.Build.Bin, _ = filepath.Abs(cfg.Build.Bin)
	cfg.Build.IncludeDir = []string{"pkg", "internal", dir}

	return runner.NewEngineWithConfig(cfg, false)
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home := os.Getenv("HOME")
		return home + path[1:], nil
	}

	var err error
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if path == "." {
		return wd, nil
	}

	if strings.HasPrefix(path, "./") {
		return wd + path[1:], nil
	}

	return path, nil
}
