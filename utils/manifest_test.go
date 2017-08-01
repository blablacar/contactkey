package utils

import (
	"testing"
)

const testfile = "testdata/manifest.yaml"

func TestLoadDeployfile(t *testing.T) {
	defaults := &DeployManifest{
		Deploy: Deploy{
			Hooks: DeployHooks{
				DeployHookNewRelic: DeployHookNewRelic{
					ApiKey: "abc123_from_config",
				},
			},
		},
	}
	f, err := LoadDeployfile(defaults, testfile)
	if err != nil {
		t.Fatalf("LoadDeployfile failed with err %q", err)
	}

	if f == nil {
		t.Fatal("f is nil")
	}

	if f.ManifestVersion != ManifestVersion {
		t.Errorf("Expected nanifestVersion %q, got %q", ManifestVersion, f.ManifestVersion)
	}

	if f.Deploy.Hooks.DeployHookNewRelic.ApiKey != "abc123_from_config" {
		t.Errorf(
			"Unexpected NewRelic.ApiKey: %q",
			f.Deploy.Hooks.DeployHookNewRelic.ApiKey,
		)
	}

	if f.Deploy.Method != "ggn" {
		t.Errorf("Unexpected deployment.method %q", f.Deploy.Method)
	}
}
