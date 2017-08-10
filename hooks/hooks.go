package hooks

type Hooks interface {
	PreDeployment(env string, service string, podVersion string) error
	PostDeployment(env string, service string, podVersion string) error
}
