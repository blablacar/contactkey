package deployers

import (
	log "github.com/sirupsen/logrus"
)

type Deployer interface {
	ListVersions(env string) (map[string]string, error)
	SetLogLevel(log.Level)
}
