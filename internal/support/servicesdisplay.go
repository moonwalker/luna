package support

import (
	"fmt"
	"strings"
)

func ServicesKeys() []string {
	return []string{"Name", "Directory", "Dep", "Watch"}
}

func (s *Service) Fields() []string {
	return []string{
		s.Name,
		s.Dir,
		strings.Join(s.Dep, ","),
		fmt.Sprintf("%v", s.Watch),
	}
}
