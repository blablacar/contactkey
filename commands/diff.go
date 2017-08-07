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

var branch = ""

func init() {
	diffCmd.PersistentFlags().StringVar(&branch, "branch", "", "Change the branch from the default one.")
}

type Diff struct {
	Env     string
	Service string
	Context *context.Context
}

func (d Diff) execute() {
	// If the branch is null it will use the default one.
	sha1, err := d.Context.Vcs.RetrieveSha1ForProject(branch)
	if err != nil {
		fmt.Printf("Failed to retrieve sha1: %q", err)
		os.Exit(1)
	}

	versions, err := d.Context.Deployer.ListVcsVersions(d.Env)
	if err != nil {
		fmt.Printf("Failed to list versions with error %q", err)
		os.Exit(1)
	}

	if len(versions) == 0 {
		fmt.Printf("No service (%q) versions found for the Env: %q", d.Service, d.Env)
		os.Exit(1)
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
			fmt.Printf("Failed to retrieve sha1: %q", err)
			os.Exit(1)
		}

		fmt.Printf("Diff between %q(deployed) and %q(branch) \n", uniqueVersion, sha1)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Author", "sha1", "description"})
		for _, change := range changes.Commits {
			table.Append([]string{change.AuthorFullName, change.DisplayId, change.Title})
		}
		table.Render()
	}
}

func (d *Diff) fill(context *context.Context, service string, env string) {
	d.Env = env
	d.Service = service
	d.Context = context

}
