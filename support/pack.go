package support

import (
	"fmt"
	"os"
)

func PackageServices(cfg Config) {
	for name, svc := range cfg.Services {
		if svc.Pack == "" || !svc.Changed {
			continue
		}

		fmt.Println("[PACK]", name)
		cmd := MakeCmd(svc.Pack, svc.Chdir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Start()
		if err != nil {
			fmt.Println(err)
		}

		err = cmd.Wait()
		if err != nil {
			fmt.Println(err)
		}
	}
}
