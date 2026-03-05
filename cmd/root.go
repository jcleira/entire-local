// Package cmd implements the CLI entry point for entire-local.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jcleira/entire-local/pkg/checkpoint"
	"github.com/jcleira/entire-local/pkg/git"
	"github.com/jcleira/entire-local/pkg/ui/dashboard"
)

var version = "dev"

func SetVersion(v string) {
	version = v
}

var rootCmd = &cobra.Command{
	Use:     "entire-local",
	Short:   "Local TUI viewer for entire.io checkpoints",
	Version: version,
	RunE:    runDashboard,
}

func Execute() {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runDashboard(_ *cobra.Command, _ []string) error {
	repoRoot, err := git.RepoRoot()
	if err != nil {
		return fmt.Errorf("not inside a git repository")
	}

	if !git.BranchExists(repoRoot, git.CheckpointBranch) {
		return fmt.Errorf("no entire.io checkpoint branch found\nEnable entire.io in this repository first")
	}

	reader := git.NewReader(repoRoot)
	loader := checkpoint.NewLoader(reader)

	return dashboard.RunDashboard(loader)
}
