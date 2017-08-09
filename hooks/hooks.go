package hooks

type Hooks interface {
	PreDeployment() error
	PostDeployment() error
}
