package deployers

type States []State

type State struct {
	Step     string
	Progress int
}

type Deployer interface {
	ListVersions(env string) (map[string]string, error)
	ListVcsVersions(env string) ([]string, error)
	Deploy(env string, podVersion string, c chan State) error
}
