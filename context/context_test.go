package context

import (
	"testing"

	"reflect"

	"github.com/remyLemeunier/contactkey/utils"
)

func TestNewContext(t *testing.T) {
	configFile, err := utils.ReadFile("../utils/testdata/config.yaml")
	if err != nil {
		t.Fatalf("ReadFile failed with err %q", err)
	}

	cfg, err := utils.LoadConfig(configFile)
	if err != nil {
		t.Fatalf("LoadConfig failed with err %q", err)
	}

	manifestFile, err := utils.ReadFile("../utils/testdata/manifest.yaml")
	if err != nil {
		t.Fatalf("ReadFile failed with err %q", err)
	}
	manifest, err := utils.LoadManifest(manifestFile)
	if err != nil {
		t.Fatalf("LoadConfig failed with err %q", err)
	}

	ctxt, err := NewContext(cfg, manifest)
	if err != nil {
		t.Fatalf("NewContext failed with err %q", err)
	}

	if reflect.TypeOf(ctxt.Deployer).String() != "*deployers.DeployerGgn" {
		t.Errorf("Type should be *deployers.DeployerGgn instead got %q", reflect.TypeOf(ctxt.Deployer).String())
	}

	if reflect.TypeOf(ctxt.Vcs).String() != "*services.Stash" {
		t.Errorf("Type should be *services.Stash instead got %q", reflect.TypeOf(ctxt.Vcs).String())
	}

	if reflect.TypeOf(ctxt.RepositoryManager).String() != "*services.Nexus" {
		t.Errorf("Type should be *services.Nexus instead got %q", reflect.TypeOf(ctxt.RepositoryManager).String())
	}

	if len(ctxt.Hooks) != 1 {
		t.Fatalf("Hooks length shound be 1 instead got %q", len(ctxt.Hooks))
	}

	if reflect.TypeOf(ctxt.Hooks[0]).String() != "*hooks.Slack" {
		t.Errorf("Type should be *hooks.Slack instead got %q", reflect.TypeOf(ctxt.Hooks[1]).String())
	}
}
