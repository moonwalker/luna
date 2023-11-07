package pm

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/moonwalker/luna/internal/config"
	"github.com/moonwalker/luna/internal/support"
)

type Runner interface {
	Run()
	Stop()
}

type PM struct {
	wg        sync.WaitGroup
	cfg       *config.Config
	Runnables []*support.Service
}

func NewPM(cfg *config.Config, allServices []*support.Service, serviceNames []string) *PM {
	pm := &PM{
		cfg:       cfg,
		Runnables: make([]*support.Service, 0),
	}

	// collect candidates with possible deps
	candidates := []*support.Service{}
	for _, svc := range allServices {
		if len(serviceNames) > 0 {
			if !support.StringListContains(serviceNames, svc.Name) {
				continue
			}
		}

		// svc's dependencies as candidates
		depSvcs := support.FindServices(svc.Dep...)
		candidates = support.AppendIfNotExists(candidates, depSvcs...)

		// svc as candidate
		candidates = support.AppendIfNotExists(candidates, svc)
	}

	// load runnable candidates
	for _, svc := range candidates {
		if svc.Runnable() {
			pm.Runnables = append(pm.Runnables, svc)
		} else {
			fmt.Printf("no command to execute for %s\n", svc.Name)
		}
	}

	return pm
}

func (pm *PM) Run() {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	defer func() {
		signal.Stop(sigs)
		cancel()
	}()

	for _, svc := range pm.Runnables {
		pm.wg.Add(1)
		go pm.runService(ctx, svc)
	}

	go func() {
		select {
		case <-sigs: // first signal, cancel context
			cancel()
		}
		<-sigs // second signal, hard exit
		os.Exit(2)
	}()

	pm.wg.Wait()
}

func (pm *PM) runService(ctx context.Context, svc *support.Service) {
	defer pm.wg.Done()

	runner, err := pm.runnerFactory(svc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	go func() {
		<-ctx.Done()
		runner.Stop()
		fmt.Print("\r")
		fmt.Println("[STOP]", svc.Name)
	}()

	fmt.Println("[START]", svc.Name)
	runner.Run()
}

func (pm *PM) runnerFactory(svc *support.Service) (Runner, error) {
	if pm.cfg.ExperimentalAirSupport {
		if len(svc.Cmd) > 0 && len(svc.Bin) > 0 {
			return air(svc.Cmd, svc.Bin, svc.Dir)
		}
	}
	return execRunner(svc)
}
