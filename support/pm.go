package support

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

type service struct {
	Chdir string
	Build string
	Start string
	Clean string
	Watch bool

	name string
	cmd  *exec.Cmd
}

type Config struct {
	Services map[string]*service
}

type PM struct {
	config Config
}

func NewPM(cfg Config) *PM {
	return &PM{
		config: cfg,
	}
}

func (pm *PM) Run() {
	pm.start()

	// wait signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sigs

	pm.stop()
}

func (pm *PM) start() {
	for name, svc := range pm.config.Services {
		svc.name = name

		if svc.Build != "" {
			pm.build(svc)
		}

		if svc.Start != "" {
			pm.spawn(svc)
		}

		if svc.Watch {
			w := pm.watch(svc)
			defer w.Close()
		}
	}
}

func (pm *PM) stop() {
	for _, svc := range pm.config.Services {
		pm.kill(svc, true)
	}
}

func (pm *PM) build(svc *service) {
	parts := strings.Fields(svc.Build)
	name := parts[0]
	arg := parts[1:len(parts)]

	cmd := exec.Command(name, arg...)
	cmd.Dir = svc.Chdir

	fmt.Println("[BUILD]", svc.name)
	out, err := cmd.CombinedOutput()
	res := strings.TrimSpace(string(out))
	if err != nil {
		fmt.Println(res)
	}
}

func (pm *PM) spawn(svc *service) {
	parts := strings.Fields(svc.Start)
	name := parts[0]
	arg := parts[1:len(parts)]

	svc.cmd = exec.Command(name, arg...)
	svc.cmd.Dir = svc.Chdir
	svc.cmd.Stdout = os.Stdout
	svc.cmd.Stdin = os.Stdin
	svc.cmd.Stderr = os.Stderr

	fmt.Println("[START]", svc.name)
	err := svc.cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
}

func (pm *PM) kill(svc *service, del bool) {
	if svc.cmd.Process == nil {
		return
	}

	fmt.Println("[KILL]", svc.name)
	err := svc.cmd.Process.Kill()
	if err != nil {
		fmt.Println(err)
	}

	if del && svc.Clean != "" {
		fmt.Println("[CLEAN]", svc.Clean)

		parts := strings.Fields(svc.Clean)
		name := parts[0]
		arg := parts[1:len(parts)]

		cmd := exec.Command(name, arg...)
		cmd.Dir = svc.Chdir

		out, err := cmd.CombinedOutput()
		res := strings.TrimSpace(string(out))
		if err != nil {
			fmt.Println(res)
		}
	}
}

func (pm *PM) watch(svc *service) *Batcher {
	watcher, err := NewWatcher(1 * time.Second)
	if err != nil {
		fmt.Println(err)
	}

	watcher.Add(svc.Chdir)

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
					// if change file is actually build file, than skip it
					if strings.HasSuffix(ev.Name, filepath.Base(svc.Start)) {
						continue
					}
					// events to watch
					importantEvent := (ev.Op == fsnotify.Create || ev.Op == fsnotify.Write || ev.Op == fsnotify.Rename || ev.Op == fsnotify.Remove)
					if importantEvent {
						fmt.Println("[CHANGE]", svc.name)
						pm.kill(svc, false)
						pm.build(svc)
						pm.spawn(svc)
						break
					}
				}
			}
		}
	}()

	return watcher
}
