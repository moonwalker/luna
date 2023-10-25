package support

import (
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type ServiceKind int

const (
	GoService ServiceKind = iota
)

type Service struct {
	Kind    ServiceKind
	Name    string
	Dir     string
	Run     string
	Dep     []string
	Watch   bool
	Changed bool
	Cmd     *exec.Cmd
}

// store services internally
var services = map[string]*Service{}

func Services() map[string]*Service {
	return services
}

// load services from yaml
func LoadYaml(f string) error {
	in, err := os.ReadFile(f)
	if err != nil {
		return err
	}

	out := struct {
		Services map[string]*Service
	}{}

	err = yaml.Unmarshal(in, &out)
	if err != nil {
		return err
	}

	// map services
	for name, s := range out.Services {
		s.Name = name
		services[name] = s
	}

	return nil
}

// register a service from lunafile
func RegisterService(name, dir string, kind ServiceKind) {
	services[name] = &Service{
		Name: name,
		Dir:  dir,
		Kind: kind,
	}
}

func FindServices(names ...string) []*Service {
	res := make([]*Service, 0)
	for _, n := range names {
		for _, s := range services {
			if s.Name == n {
				res = append(res, s)
			}
		}
	}
	return res
}
