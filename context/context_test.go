package context

import (
	"testing"

	"reflect"

	"github.com/remyLemeunier/contactkey/utils"
	"github.com/spf13/viper"
)

func TestNewContext(t *testing.T) {
	_, err := utils.ReadFile("../utils/testdata/config.yaml")
	if err != nil {
		t.Fatalf("ReadFile failed with err %q", err)
	}

	viper.AddConfigPath("../utils/testdata")
	cfg, err := utils.LoadConfig()
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

	if reflect.TypeOf(ctxt.Binaries).String() != "*services.Nexus" {
		t.Errorf("Type should be *services.Nexus instead got %q", reflect.TypeOf(ctxt.Binaries).String())
	}

	if len(ctxt.Hooks) != 3 {
		t.Fatalf("Unexpected Hooks length: %q", len(ctxt.Hooks))
	}

	if reflect.TypeOf(ctxt.Hooks[0]).String() != "*hooks.Slack" {
		t.Fatalf("Type should be *hooks.Slack instead got %q", reflect.TypeOf(ctxt.Hooks[0]).String())
	}

	if ctxt.Hooks[0].StopOnError() != false {
		t.Error("Exepected StopOnError from Slack to be false")
	}

	if reflect.TypeOf(ctxt.Hooks[1]).String() != "*hooks.NewRelicClient" {
		t.Fatalf("Type should be *hooks.NewRelicClient instead got %q", reflect.TypeOf(ctxt.Hooks[1]).String())
	}
	if reflect.TypeOf(ctxt.Hooks[2]).String() != "*hooks.ExecCommand" {
		t.Fatalf("Type should be *hooks.ExecCommand instead got %q", reflect.TypeOf(ctxt.Hooks[2]).String())
	}

	if ctxt.Hooks[2].StopOnError() != true {
		t.Error("Exepected StopOnError from ExecCommand to be true")
	}

	if reflect.TypeOf(ctxt.LockSystem).String() != "*utils.FileLock" {
		t.Errorf("Type should be *utils.FileLock instead got %q", reflect.TypeOf(ctxt.LockSystem).String())
	}
}
