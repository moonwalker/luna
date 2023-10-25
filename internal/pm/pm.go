package pm

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/moonwalker/luna/internal/support"
	"github.com/moonwalker/luna/internal/watcher"
)

type PM struct {
	// loaded services
	services map[string]*support.Service
}

func NewPM(allServices map[string]*support.Service, serviceNames []string) *PM {
	pm := &PM{
		services: make(map[string]*support.Service, 0),
	}

	// load services
	for name, svc := range allServices {
		if len(serviceNames) > 0 && !stringInSlice(name, serviceNames) {
			continue
		} else {
			for _, dep := range svc.Dep {
				pm.services[dep] = allServices[dep]
			}
		}
		pm.services[name] = svc
	}

	return pm
}

func (pm *PM) Run() {
	for _, svc := range pm.services {
		if len(svc.Run) > 0 {
			pm.spawn(svc)
			if svc.Watch {
				pm.watch(svc)
			}
		}
	}
	waitSig()
	pm.Stop()
}

func (pm *PM) Stop() {
	for _, svc := range pm.services {
		pm.kill(svc)
	}
}

func (pm *PM) spawn(svc *support.Service) {
	fmt.Println("[START]", svc.Name)

	svc.Cmd = makeCmd(svc.Run, svc.Dir)
	svc.Cmd.Stdout = os.Stdout
	svc.Cmd.Stderr = os.Stderr

	err := svc.Cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
}

func (pm *PM) kill(svc *support.Service) {
	if svc.Cmd.Process == nil {
		return
	}

	fmt.Println("[KILL]", svc.Name)
	err := svc.Cmd.Process.Kill()
	if err != nil {
		fmt.Println(err)
	}
}

func (pm *PM) watch(svc *support.Service) *watcher.Batcher {
	watcher, err := watcher.NewWatcher(1 * time.Second)
	if err != nil {
		fmt.Println(err)
	}

	watcher.Add(svc.Dir)

	go func() {
		for {
			select {
			case evs := <-watcher.Events:
				// fmt.Println("Received System Events:", evs)
				for _, ev := range evs {
					// sometimes during rm -rf operations a '"": REMOVE' is triggered, just ignore these
					if ev.Name == "" {
						continue
					}
					// events to watch
					importantEvent := (ev.Op == fsnotify.Create || ev.Op == fsnotify.Write || ev.Op == fsnotify.Rename || ev.Op == fsnotify.Remove)
					if importantEvent {
						fmt.Println("[CHANGE]", svc.Name)
						pm.kill(svc)
						pm.spawn(svc)
						break
					}
				}
			}
		}
	}()

	return watcher
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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

func waitSig() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sigs
}
