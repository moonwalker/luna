package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type service struct {
	Dir string
	Cmd string
	Del string

	exe exec.Cmd
}

type config struct {
	Services map[string]service
}

var procs []service

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		var cfg config

		err := viper.Unmarshal(&cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		procs = make([]service, 0)
		for name, svc := range cfg.Services {
			spawn(name, svc)
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-sigs

		for _, cmd := range procs {
			kill(cmd)
		}
	},
}

func init() {
	servicesCmd.AddCommand(runCmd)
}

func spawn(name string, svc service) {
	cmd := exec.Command("/bin/sh", "-c", svc.Cmd)
	cmd.Dir = svc.Dir
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}

	svc.exe = *cmd
	procs = append(procs, svc)
}

func kill(svc service) {
	err := svc.exe.Process.Kill()
	if err != nil {
		fmt.Println(err)
	}

	if svc.Del != "" {
		fp, _ := filepath.Abs(svc.Dir + "/" + svc.Del)
		os.Remove(fp)
	}
}
