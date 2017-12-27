package commands

import (
	"fmt"
	"os"

	"errors"

	"github.com/olekukonko/tablewriter"
	"github.com/remyLemeunier/contactkey/context"
	"github.com/remyLemeunier/contactkey/deployers"
	log "github.com/sirupsen/logrus"
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

func (d Diff) execute() error {
	// If the branch is null it will use the default one.
	sha1, err := d.Context.Vcs.RetrieveSha1ForProject(branch)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to retrieve sha1: %q \n", err))
	}
	if sha1 == "" {
		return errors.New(fmt.Sprintf("No sha1 found for service %q \n", d.Service))
	}

	uniqueVersions, err := deployers.ListUniqueVcsVersions(d.Context.Deployer, d.Env)
	if err != nil {
		return err
	}

	for _, uniqueVersion := range uniqueVersions {
		changes, err := d.Context.Vcs.Diff(uniqueVersion, sha1)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to retrieve sha1: %q \n", err))
		}

		log.Println(fmt.Sprintf("Diff between %q(deployed) and %q(branch) \n", uniqueVersion, sha1))
		d.TableWriter.SetHeader([]string{"Author", "sha1", "description"})
		for _, change := range changes.Commits {
			d.TableWriter.Append([]string{change.AuthorFullName, change.DisplayId, change.Title})
		}
		d.TableWriter.Render()
	}
	return nil
}

func (d *Diff) fill(context *context.Context, service string, env string) {
	d.Env = env
	d.Service = service
	d.Context = context
	d.TableWriter = tablewriter.NewWriter(os.Stdout)
}
