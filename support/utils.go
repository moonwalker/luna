package support

import (
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func BoolTostring(b bool, s string) string {
	if b {
		return s
	}
	return ""
}

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

func waitSig() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sigs
}

func makeCmd(command string, chdir string) *exec.Cmd {
	parts := strings.Fields(command)
	name := parts[0]
	arg := parts[1:len(parts)]

	cmd := exec.Command(name, arg...)
	if chdir != "" {
		cmd.Dir = chdir
	}

	return cmd
}
