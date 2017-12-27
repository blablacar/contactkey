package hooks

type HookInformation struct {
	UserName      string
	Env           string
	Service       string
	PodVersion    string
	Miscellaneous []string
}

type Hooks interface {
	Init() error
	PreDeployment(hookinformation HookInformation) error
	PostDeployment(hookinformation HookInformation) error
	StopOnError() bool
}
