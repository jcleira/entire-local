// Package git provides read-only access to git repositories via plumbing commands.
package git

import (
	"os/exec"
	"strings"
)

const CheckpointBranch = "entire/checkpoints/v1"

func ExecGit(repoRoot string, args ...string) (string, error) {
	fullArgs := append([]string{"-C", repoRoot}, args...)
	cmd := exec.Command("git", fullArgs...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func RepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func BranchExists(repoRoot, branch string) bool {
	_, err := ExecGit(repoRoot, "rev-parse", "--verify", "refs/heads/"+branch)
	return err == nil
}
