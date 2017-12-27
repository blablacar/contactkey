package deployers

import (
	"testing"
)

func TestDeployerRegistry(t *testing.T) {
}

func TestListUniqueVcsVersions(t *testing.T) {
	envs := make(map[string]string)
	envs["staging"] = "staging"
	execCommand = mockggn
	d := DeployerGgn{
		Service:      "webhooks",
		Pod:          "pod-webhooks",
		Environments: envs,
		VcsRegexp:    "-(.+)",
	}

	uniqueVersions, err := ListUniqueVcsVersions(&d, "staging")
	if err != nil {
		t.Fatalf("Error trying to retrieve unique version: %q", err)
	}

	if len(uniqueVersions) != 1 {
		t.Fatalf("uniqueVersions len should be 1 instead got %s", len(uniqueVersions))
	}

	if uniqueVersions[0] != "1" {
		t.Error("uniqueVersion should \"1\" instead got %s", uniqueVersions[0])
	}
}
