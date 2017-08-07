package deployers

type Deployer interface {
	ListVersions(env string) (map[string]string, error)
}
