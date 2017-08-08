package commands

import (
	"bytes"
	"strings"
	"testing"

	"github.com/olekukonko/tablewriter"
	"github.com/remyLemeunier/contactkey/context"
	"github.com/remyLemeunier/contactkey/services"
)

type VcsMock struct{}

func (v *VcsMock) RetrieveSha1ForProject(branch string) (string, error) {
	return "abcde", nil
}

func (v *VcsMock) Diff(deployedSha1 string, sha1ToDeploy string) (*services.Changes, error) {
	changes := new(services.Changes)
	commits := services.Commits{}
	commits.Title = "Title"
	commits.DisplayId = "DisplayId"
	commits.AuthorFullName = "AuthorFullName"
	commits.AuthorSlug = "AuthorSlug"
	changes.Commits = append(changes.Commits, commits)

	return changes, nil
}

func TestDiffExecute(t *testing.T) {
	cmd := &Diff{
		Env:     "staging",
		Service: "webhooks",
		Context: &context.Context{
			Deployer: &DeployerMockGgn{},
			Vcs:      &VcsMock{},
		},
	}
	// Catch stdout
	out := new(bytes.Buffer)
	cmd.TableWriter = tablewriter.NewWriter(out)
	cmd.Writer = out

	cmd.execute()
	if out.String() == "" {
		t.Errorf("Unexpected stdout : %q", out)
	}

	// Check if we display at the least the right information
	if !strings.Contains(out.String(), "Diff between \"b0f586a\"(deployed) and \"abcde\"(branch) \n") {
		t.Error("AuthorFullName not found")
	}

	if !strings.Contains(out.String(), "AuthorFullName") {
		t.Error("AuthorFullName not found")
	}

	if !strings.Contains(out.String(), "DisplayId") {
		t.Error("DisplayId not found")
	}

	if !strings.Contains(out.String(), "Title") {
		t.Error("Title not found")
	}
}
