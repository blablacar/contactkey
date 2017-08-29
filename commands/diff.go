package commands

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/remyLemeunier/contactkey/context"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Diff between what's currently deployed and what's going to be deployed",
}

func init() {
	diffCmd.PersistentFlags().StringVar(&branch, "branch", "", "Change the branch from the default one.")
}

type Diff struct {
	Env         string
	Service     string
	Context     *context.Context
	TableWriter *tablewriter.Table
}

func (d Diff) execute() {
	// If the branch is null it will use the default one.
	sha1, err := d.Context.Vcs.RetrieveSha1ForProject(branch)
	if err != nil {
		d.Context.Log.Error(fmt.Sprintf("Failed to retrieve sha1: %q \n", err))
		return
	}
	if sha1 == "" {
		d.Context.Log.Error(fmt.Sprintf("No sha1 found for service %q \n", d.Service))
		return
	}

	versions, err := d.Context.Deployer.ListVcsVersions(d.Env)
	if err != nil {
		d.Context.Log.Error(fmt.Sprintf("Failed to list versions with error %q \n", err))
		return
	}

	if len(versions) == 0 {
		d.Context.Log.Error(fmt.Sprintf("No service (%q) versions found for the Env: %q \n", d.Service, d.Env))
		return
	}

	// Retrieve only unique versions
	encountered := map[string]bool{}
	for v := range versions {
		encountered[versions[v]] = true
	}
	uniqueVersions := []string{}
	for key := range encountered {
		uniqueVersions = append(uniqueVersions, key)
	}

	for _, uniqueVersion := range uniqueVersions {
		changes, err := d.Context.Vcs.Diff(uniqueVersion, sha1)
		if err != nil {
			d.Context.Log.Error(fmt.Sprintf("Failed to retrieve sha1: %q \n", err))
			return
		}

		d.Context.Log.Println(fmt.Sprintf("Diff between %q(deployed) and %q(branch) \n", uniqueVersion, sha1))
		d.TableWriter.SetHeader([]string{"Author", "sha1", "description"})
		for _, change := range changes.Commits {
			d.TableWriter.Append([]string{change.AuthorFullName, change.DisplayId, change.Title})
		}
		d.TableWriter.Render()
	}
}

func (d *Diff) fill(context *context.Context, service string, env string) {
	d.Env = env
	d.Service = service
	d.Context = context
	d.TableWriter = tablewriter.NewWriter(os.Stdout)
}
