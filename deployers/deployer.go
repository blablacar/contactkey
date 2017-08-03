package deployers

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

var Registry = make(map[string]Deployer)

type Deployer interface {
	ListVersions(env string) (map[string]string, error)
	SetLogLevel(log.Level)
}

func MakeInstance(name string) (Deployer, error) {
	_, ok := Registry[name]
	if !ok {
		return nil, errors.New(
			fmt.Sprintf("Unexpected Deployer type %q", name),
		)
	}
	Registry[name].SetLogLevel(log.DebugLevel)
	return Registry[name], nil
}
