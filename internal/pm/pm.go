package pm

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/moonwalker/luna/internal/support"
	"github.com/moonwalker/luna/internal/watcher"
)

type Runner interface {
	Run()
	Stop()
}

type PM struct {
	// loaded services
	services map[string]*support.Service
	// runners
	runners map[string]Runner
}

func NewPM(allServices map[string]*support.Service, serviceNames []string) *PM {
	pm := &PM{
		services: make(map[string]*support.Service, 0),
		runners:  make(map[string]Runner, 0),
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
		if len(svc.Run) > 0 || (len(svc.Cmd) > 0 && len(svc.Bin) > 0) {
			pm.spawn(svc)
			if svc.Watch {
				pm.watch(svc)
			}
		} else {
			fmt.Printf("no command specified for %s\n", svc.Name)
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

func (pm *PM) runnerFactory(svc *support.Service) (Runner, error) {
	// if svc.Kind == support.GoService {
	// 	return air(svc.Cmd, svc.Bin, svc.Dir)
	// }
	return newCmdRunner(svc.Run, svc.Dir)
}

func (pm *PM) spawn(svc *support.Service) {
	fmt.Println("[START]", svc.Name)

	runner, err := pm.runnerFactory(svc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	pm.runners[svc.Name] = runner
	pm.runners[svc.Name].Run()
}

func (pm *PM) kill(svc *support.Service) {
	if pm.runners[svc.Name] != nil {
		fmt.Println("[KILL]", svc.Name)
		pm.runners[svc.Name].Stop()
	}
}

func (pm *PM) watch(svc *support.Service) *watcher.Batcher {
	watcher, err := watcher.NewWatcher(1 * time.Second)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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

func waitSig() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-sigs
}
