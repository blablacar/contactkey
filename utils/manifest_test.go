package utils

import (
	"testing"
)

const testfile = "testdata/manifest.yaml"

func TestLoadDeployfile(t *testing.T) {
	f, err := LoadDeployfile(testfile)
	if err != nil {
		t.Fatalf("LoadDeployfile failed with err %q", err)
	}

	if f == nil {
		t.Fatal("f is nil")
	}

	if f.ManifestVersion != ManifestVersion {
		t.Errorf("Expected nanifestVersion %q, got %q", ManifestVersion, f.ManifestVersion)
	}

	if f.Deploy.Method != "ggn" {
		t.Errorf("Unexpected deployment.method %q", f.Deploy.Method)
	}
}
