package support

import (
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func makeCmd(command string, chdir string) *exec.Cmd {
	parts := strings.Fields(command)
	name := parts[0]
	arg := parts[1:len(parts)]

	cmd := exec.Command(name, arg...)
	cmd.Dir = chdir

	return cmd
}

func waitSig() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sigs
}
