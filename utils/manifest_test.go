package utils

import (
	"testing"
)

func TestLoadManifest(t *testing.T) {
	configFile, err := ReadFile("./testdata/manifest.yaml")
	if err != nil {
		t.Fatalf("ReadFile failed with err %q", err)
	}
	manifest, err := LoadManifest(configFile)
	if err != nil {
		t.Fatalf("LoadDeployfile failed with err %q", err)
	}

	if manifest == (&Manifest{}) {
		t.Errorf("Unexpected manifest %q", manifest)
	}

	if manifest.Pod != "pod-webhooks" {
		t.Errorf("Pod should be 'pod-webhooks' instead got  %q", manifest.Pod)
	}
}
