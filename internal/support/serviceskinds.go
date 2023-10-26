package support

import (
	"fmt"
)

type ServiceKind int

const (
	GenericService ServiceKind = iota
	GoService

	// defaults
	buildTmp = "./tmp"

	goServiceRun   = "go run ."
	goServiceBuild = "go build -o %s/%s %s"
)

// enfore specific attributes
func (s *Service) SetKind(kind ServiceKind) error {
	s.Kind = GoService

	dir, err := ExpandPath(s.Dir)
	if err != nil {
		return err
	}

	s.Run = goServiceRun
	s.Cmd = fmt.Sprintf(goServiceBuild, buildTmp, s.Name, dir)
	s.Bin = fmt.Sprintf("%s/%s", buildTmp, s.Name)

	return nil
}
