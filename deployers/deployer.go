package deployers

type Deployer interface {
	ListVersions(env string) (map[string]string, error)
	ListVcsVersions(env string) ([]string, error)
	Deploy(env string, podVersion string) error
}
