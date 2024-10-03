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
	svc *support.Service
	cmd *exec.Cmd
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
	if r.cmd.Process != nil {
		err := r.cmd.Process.Signal(os.Interrupt)
		if err != nil {
			err = r.cmd.Process.Kill()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func (r *cmdExecRunner) start() {
	r.cmd = makeCmd(r.svc.Run, r.svc.Dir)
	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr

	err := r.cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	err = r.cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (r *cmdExecRunner) restart() {
	if r.cmd.Process != nil {
		fmt.Println("[RESTART]", r.svc.Name)
		r.Stop()
		r.start()
	}
}

func makeCmd(command string, dir string) *exec.Cmd {
	parts := strings.Fields(command)
	name := parts[0]
	args := parts[1:]

	cmd := exec.Command(name, args...)
	cmd.Env = support.Environ(dir)
	if dir != "" {
		cmd.Dir = dir
	}

	return cmd
}
