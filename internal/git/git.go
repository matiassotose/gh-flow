package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DetectRepositories finds all git repositories in subdirectories
func DetectRepositories(root string) ([]string, error) {
	repos := make([]string, 0)

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		repoPath := filepath.Join(root, entry.Name())
		gitPath := filepath.Join(repoPath, ".git")

		if _, err := os.Stat(gitPath); err == nil {
			// It's a git repository
			repos = append(repos, entry.Name())
		}
	}

	return repos, nil
}

// HasUncommittedChanges checks if a repository has uncommitted changes
func HasUncommittedChanges(repoPath string) (bool, error) {
	output, err := execGitCommandOutput(repoPath, "status", "--porcelain")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(output) != "", nil
}

// StashChanges stashes all changes in a repository
func StashChanges(repoPath string) error {
	return execGitCommand(repoPath, "stash", "push", "-m", "gh-flow auto-stash")
}

// CheckoutBranch switches to a branch
func CheckoutBranch(repoPath, branch string) error {
	return execGitCommand(repoPath, "checkout", branch)
}

// Pull pulls latest changes from origin
func Pull(repoPath string) error {
	return execGitCommand(repoPath, "pull")
}

// CreateBranch creates and checks out a new branch
func CreateBranch(repoPath, branchName string) error {
	return execGitCommand(repoPath, "checkout", "-b", branchName)
}

// GetCurrentBranch returns the name of the current branch
func GetCurrentBranch(repoPath string) (string, error) {
	output, err := execGitCommandOutput(repoPath, "branch", "--show-current")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// GetOriginBranch tries to determine from which branch the current branch was created
func GetOriginBranch(repoPath, currentBranch string) (string, error) {
	// Try to get the upstream branch
	output, err := execGitCommandOutput(repoPath, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err != nil {
		return "", err
	}

	upstream := strings.TrimSpace(output)
	parts := strings.Split(upstream, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-1], nil
	}

	return "", fmt.Errorf("could not determine origin branch")
}

// AddAll stages all changes
func AddAll(repoPath string) error {
	return execGitCommand(repoPath, "add", "-A")
}

// Commit creates a commit with the given message
func Commit(repoPath, message string) error {
	return execGitCommand(repoPath, "commit", "-m", message)
}

// Push pushes the current branch to origin
func Push(repoPath, branch string) error {
	return execGitCommand(repoPath, "push", "-u", "origin", branch)
}

// Helper functions
func execGitCommand(repoPath string, args ...string) error {
	cmd := execGit(repoPath, args...)
	return cmd.Run()
}

func execGitCommandOutput(repoPath string, args ...string) (string, error) {
	cmd := execGit(repoPath, args...)
	output, err := cmd.Output()
	return string(output), err
}

func execGit(repoPath string, args ...string) Cmd {
	fullArgs := append([]string{"-C", repoPath}, args...)
	return Command("git", fullArgs...)
}

// Cmd interface for mocking in tests
type Cmd interface {
	Run() error
	Output() ([]byte, error)
}

// Command creates a new command
type realCmd struct {
	name string
	args []string
}

func Command(name string, args ...string) Cmd {
	return &realCmd{name: name, args: args}
}

func (c *realCmd) Run() error {
	cmd := execCommand(c.name, c.args...)
	return cmd.Run()
}

func (c *realCmd) Output() ([]byte, error) {
	cmd := execCommand(c.name, c.args...)
	return cmd.Output()
}

var execCommand = exec.Command
