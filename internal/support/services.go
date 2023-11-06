package support

import (
	"errors"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

type Service struct {
	Kind  ServiceKind
	Name  string
	Dir   string
	Run   string // simple run
	Cmd   string // build
	Bin   string // built bin to run
	Dep   []string
	Watch bool
}

var (
	// store services internally
	services            = []*Service{}
	invalidServiceError = errors.New("invalid service definition")
)

func (s *Service) Runnable() bool {
	if len(s.Dir) == 0 {
		return false
	}
	return len(s.Run) > 0 || (len(s.Cmd) > 0 && len(s.Bin) > 0)
}

func Services() []*Service {
	return services
}

func ServicesSorted() []*Service {
	res := make([]*Service, 0, len(services))

	for _, svc := range services {
		res = append(res, svc)
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	return res
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
		services = append(services, s)
	}

	return nil
}

// register a service from lunafile
func RegisterService(s *Service) error {
	if s == nil || len(s.Name) == 0 {
		return invalidServiceError
	}
	services = append(services, s)
	return nil
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
