package hooks

type Hooks interface {
	Init() error
	PreDeployment(userName string, env string, service string, podVersion string) error
	PostDeployment(userName string, env string, service string, podVersion string) error
	StopOnError() bool
}
