package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return false, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return false, err
	}

	status, err := worktree.Status()
	if err != nil {
		return false, err
	}

	return !status.IsClean(), nil
}

// StashChanges stashes all changes in a repository
func StashChanges(repoPath string) error {
	// Verify it's a valid repo
	_, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	// Use git command for stash as go-git doesn't support it directly
	return execGitCommand(repoPath, "stash", "push", "-m", "gh-flow auto-stash")
}

// CheckoutBranch switches to a branch
func CheckoutBranch(repoPath, branch string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	return worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	})
}

// Pull pulls latest changes from origin
func Pull(repoPath string) error {
	return execGitCommand(repoPath, "pull")
}

// CreateBranch creates and checks out a new branch
func CreateBranch(repoPath, branchName string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	headRef, err := repo.Head()
	if err != nil {
		return err
	}

	// Create new branch from current HEAD
	ref := plumbing.NewHashReference(plumbing.NewBranchReferenceName(branchName), headRef.Hash())
	if err := repo.Storer.SetReference(ref); err != nil {
		return err
	}

	// Checkout the new branch
	return worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
	})
}

// GetCurrentBranch returns the name of the current branch
func GetCurrentBranch(repoPath string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}

	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	if !head.Name().IsBranch() {
		return "", fmt.Errorf("HEAD is not a branch")
	}

	return head.Name().Short(), nil
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
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	return worktree.AddWithOptions(&git.AddOptions{
		All: true,
	})
}

// Commit creates a commit with the given message
func Commit(repoPath, message string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	_, err = worktree.Commit(message, &git.CommitOptions{})
	return err
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
