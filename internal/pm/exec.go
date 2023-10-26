package pm

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/moonwalker/luna/internal/support"
	"github.com/moonwalker/luna/internal/watcher"
)

type cmdExecRunner struct {
	svc  *support.Service
	proc *os.Process
}

func execRunner(svc *support.Service) (*cmdExecRunner, error) {
	return &cmdExecRunner{svc, nil}, nil
}

func (r *cmdExecRunner) Run() {
	if r.svc.Watch {
		watcher.Watch(r.svc.Dir, func() {
			r.restart()
		})
	}
	r.start()
}

func (r *cmdExecRunner) Stop() {
	if r.proc != nil {
		r.proc.Kill()
	}
}

func (r *cmdExecRunner) start() {
	cmd := makeCmd(r.svc.Run, r.svc.Dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// keep underlying process for later
	r.proc = cmd.Process

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (r *cmdExecRunner) restart() {
	if r.proc != nil {
		fmt.Println("[RESTART]", r.svc.Name)
		r.proc.Kill()
		r.start()
	}
}

func makeCmd(command string, dir string) *exec.Cmd {
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
