package hooks

type hooks interface {
	preDeployment()
	postDeployment()
}

func PreDeployment(h hooks) {
	h.preDeployment()
}

func PostDeployment(h hooks) {
	h.postDeployment()
}
