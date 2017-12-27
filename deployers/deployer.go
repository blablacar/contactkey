package deployers

import (
	"errors"
	"fmt"
)

type State struct {
	Step     string
	Progress int
}

type States []State

type Instance struct {
	Name    string
	State   string
	Version string
}

type Deployer interface {
	ListInstances(env string) ([]Instance, error)
	ListVcsVersions(env string) ([]string, error)
	Deploy(env string, podVersion string, c chan State) error
}

func ListUniqueVcsVersions(d Deployer, env string) ([]string, error) {
	uniqueVersions := []string{}
	versions, err := d.ListVcsVersions(env)
	if err != nil {
		return uniqueVersions, errors.New(fmt.Sprintf("Failed to list versions with error %q \n", err))
	}

	if len(versions) == 0 {
		return uniqueVersions, errors.New(fmt.Sprintf("No versions found for the Env: %q \n", env))
	}

	// Retrieve only unique versions
	encountered := map[string]bool{}
	for v := range versions {
		encountered[versions[v]] = true
	}

	for key := range encountered {
		uniqueVersions = append(uniqueVersions, key)
	}

	return uniqueVersions, nil
}
