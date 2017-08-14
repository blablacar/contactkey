package commands

import (
	"fmt"
	"io"
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
	Writer      io.Writer
}

func (d Diff) execute() {
	// If the branch is null it will use the default one.
	sha1, err := d.Context.Vcs.RetrieveSha1ForProject(branch)
	if err != nil {
		fmt.Fprintf(d.Writer, "Failed to retrieve sha1: %q", err)
		os.Exit(1)
	}
	if sha1 == "" {
		fmt.Fprintf(d.Writer, "No sha1 found for service %q \n", d.Service)
		os.Exit(1)
	}

	versions, err := d.Context.Deployer.ListVcsVersions(d.Env)
	if err != nil {
		fmt.Fprintf(d.Writer, "Failed to list versions with error %q", err)
		os.Exit(1)
	}

	if len(versions) == 0 {
		fmt.Fprintf(d.Writer, "No service (%q) versions found for the Env: %q", d.Service, d.Env)
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
			fmt.Fprintf(d.Writer, "Failed to retrieve sha1: %q", err)
			os.Exit(1)
		}

		fmt.Fprintf(d.Writer, "Diff between %q(deployed) and %q(branch) \n", uniqueVersion, sha1)
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
	d.Writer = os.Stdout
}
