package pm

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type cmdRunner struct {
	run  string
	dir  string
	proc *os.Process
}

func newCmdRunner(run, dir string) (*cmdRunner, error) {
	return &cmdRunner{run, dir, nil}, nil
}

func (r *cmdRunner) Run() {
	cmd := makeExecCmd(r.run, r.dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	r.proc = cmd.Process
}

func (r *cmdRunner) Stop() {
	if r.proc != nil {
		r.proc.Kill()
	}
}

func makeExecCmd(command string, dir string) *exec.Cmd {
	parts := strings.Fields(command)
	name := parts[0]
	args := parts[1:]

	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()
	if dir != "" {
		cmd.Dir = dir
	}

	return cmd
}
