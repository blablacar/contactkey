package commands

import (
	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback to a previous version",
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the service in an environment",
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Diff between what's currently deployed and what's going to be deployed",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the version of the service",
}