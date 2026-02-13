package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gh-flow",
	Short: "CLI tool for streamlined Git workflow across multiple repositories",
	Long: `gh-flow automates the workflow for developing features and hotfixes
across multiple repositories in a system (frontend, backend, microservices, etc.).

It handles:
  - Repository preparation (branch selection, pulling, creating feature branches)
  - Development finalization (committing, pushing, creating PRs)

Commands:
  start   Prepare repositories for development
  finish  Finalize development and create PRs`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(finishCmd)
}
