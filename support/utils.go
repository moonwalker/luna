package support

import (
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func AppendUnique(slice []string, s string) []string {
	if s == "" {
		return slice
	}
	for _, e := range slice {
		if e == s {
			return slice
		}
	}
	return append(slice, s)
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func MakeCmd(command string, dir string) *exec.Cmd {
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

func waitSig() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sigs
}
