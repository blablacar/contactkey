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

	if manifest.NexusManifest.Artifact != "pod-webhooks" {
		t.Errorf("artifact in the NexusManifest not found got %q", manifest.NexusManifest.Artifact)
	}

	if manifest.SlackManifest.Channel != "channel" {
		t.Errorf("channel in the SlackManifest not found got %q", manifest.SlackManifest.Channel)
	}
}
