package support

type ServiceKind int

const (
	GoService ServiceKind = iota
)

type Service struct {
	Name string
	Dir  string
	Kind ServiceKind
}

var services = []*Service{}

func RegisterService(name, dir string, kind ServiceKind) {
	services = append(services, &Service{
		Name: name,
		Dir:  dir,
		Kind: kind,
	})
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

func AllServices() []*Service {
	return services
}
