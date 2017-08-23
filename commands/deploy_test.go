package commands

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/remyLemeunier/contactkey/context"
)

func TestExecute(t *testing.T) {
	out := new(bytes.Buffer)
	d := &Deploy{
		Context: &context.Context{
			Deployer:          &DeployerMockGgn{},
			Vcs:               &VCSMock{},
			RepositoryManager: &RepositoryManagerMock{},
		},
		Writer: out,
	}
	d.execute()
	if out.String() == "" {
		t.Errorf("Unexpected stdout : %q", out)
	}

	if !regexp.MustCompile(`locking : 100`).MatchString(out.String()) {
		t.Errorf("Stdout is missing locking step : %q", out)
	}

}
