package hooks

type Hooks interface {
	PreDeployment(userName string, env string, service string, podVersion string) error
	PostDeployment(userName string, env string, service string, podVersion string) error
}
