package support

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type PM struct {
	config   Config
	services map[string]*service
}

func NewPM(cfg Config, serviceNames []string) *PM {
	pm := &PM{
		config:   cfg,
		services: make(map[string]*service, 0),
	}
	for name, svc := range pm.config.Services {
		if len(serviceNames) > 0 && !StringInSlice(name, serviceNames) {
			continue
		} else {
			for _, dep := range svc.Dep {
				pm.services[dep] = pm.config.Services[dep]
			}
		}
		pm.services[name] = svc
	}
	return pm
}

func (pm *PM) Run() {
	for name, svc := range pm.services {
		svc.name = name
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

func (pm *PM) spawn(svc *service) {
	fmt.Println("[START]", svc.name)

	svc.cmd = MakeCmd(svc.Run, svc.Dir)
	svc.cmd.Stdout = os.Stdout
	svc.cmd.Stderr = os.Stderr

	err := svc.cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
}

func (pm *PM) kill(svc *service) {
	if svc.cmd.Process == nil {
		return
	}

	fmt.Println("[KILL]", svc.name)
	err := svc.cmd.Process.Kill()
	if err != nil {
		fmt.Println(err)
	}
}

func (pm *PM) watch(svc *service) *Batcher {
	watcher, err := NewWatcher(1 * time.Second)
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
						fmt.Println("[CHANGE]", svc.name)
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
