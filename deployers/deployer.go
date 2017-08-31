package deployers

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
