package deployers

import (
	"errors"
	"fmt"
)

var Registry = make(map[string]Deployer)

type Deployer interface {
	ListVersions(env string) (map[string]string, error)
}

func MakeInstance(name string) (Deployer, error) {
	_, ok := Registry[name]
	if !ok {
		return nil, errors.New(
			fmt.Sprintf("Unexpected Deployer type %q", name),
		)
	}
	return Registry[name], nil
}
