package deployers

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/remyLemeunier/contactkey/utils"
	"github.com/remyLemeunier/k8s-deploy/releases"
)

func TestNewDeployerK8s(t *testing.T) {
	cfg := utils.DeployerK8sConfig{
		Environments: map[string]utils.K8sEnvironment{
			"testEnv": {
				Cluster: "testCluster",
			},
		},
	}
	man := utils.DeployerK8sManifest{
		Release: "test",
	}

	d, err := NewDeployerK8s(cfg, man)
	if err != nil {
		t.Fatalf("Unexpected error: %q", err)
	}

	if d.Environments["testEnv"] == (utils.K8sEnvironment{}) {
		t.Errorf("Unexpected Environments: %q", d.Environments)
	}
}

func TestGetContext(t *testing.T) {
	//TODO
}

func TestGetReleasePath(t *testing.T) {
	d := &DeployerK8s{
		ReleaseName: "sleepy",
		workPath:    "/workpath",
	}
	rp := d.getReleasePath("cluster", "namespace")
	if rp != "/workpath/deployments/cluster/namespace/sleepy/release.yaml" {
		t.Errorf("Unexpected releasePath: %q", rp)
	}
}

func (d *DeployerK8s) MockGetRelease() {}

func TestDeploy(t *testing.T) {
	d := &DeployerK8s{
		Release:   &releases.FakeRelease{},
		Namespace: "payment",
		Environments: map[string]utils.K8sEnvironment{
			"test": {
				Cluster: "cluster",
			},
		},
	}

	c := make(chan State, 1)
	go func() {
		for {
			state := <-c
			spew.Dump(state)
		}
	}()
	err := d.Deploy("test", "1.0", c)
	if err != nil {
		t.Fatalf("Unexpected error : %q", err)
	}

}
