package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/cosmtrek/air/runner"
)

func main() {
	air()
}

// air --build.cmd "go build -o ./tmp/service_a ./examples/service_a" --build.bin "./tmp/service_a" --build.exclude_dir "dist,tmp,vendor"
func air() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	p, _ := expandPath(".")
	cfg := &runner.Config{}
	cfg.Root = p
	cfg.TmpDir = "tmp"
	cfg.TestDataDir = "testdata"
	cfg.Build.Log = "build-errors.log"

	cfg.Build.Cmd = "go build -o ./tmp/service_a ./examples/service_a"
	cfg.Build.Bin = "./tmp/service_a"

	cfg.Build.ExcludeDir = []string{"assets", "dist", "tmp", "vendor", "testdata"}
	cfg.Build.IncludeExt = []string{"go", "tpl", "tmpl", "html"}
	cfg.Build.ExcludeRegex = []string{"_test.go"}

	cfg.Build.Delay = 1000
	cfg.Build.Rerun = false
	cfg.Build.RerunDelay = 500

	cfg.Build.Bin, _ = filepath.Abs(cfg.Build.Bin)

	r, err := runner.NewEngineWithConfig(cfg, false)
	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		<-sigs
		r.Stop()
	}()

	defer func() {
		if e := recover(); e != nil {
			log.Fatalf("PANIC: %+v", e)
		}
	}()

	r.Run()
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
