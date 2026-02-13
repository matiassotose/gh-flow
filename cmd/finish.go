package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"gh-flow/internal/git"
	"gh-flow/internal/prompt"
)

var finishCmd = &cobra.Command{
	Use:   "finish",
	Short: "Finalize development and create PRs",
	Long:  `Commits changes, pushes branches, and creates pull requests for all selected repositories.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runFinish()
	},
}

func runFinish() error {
	// Step 1: Detect repositories
	fmt.Println("🔍 Scanning for git repositories...")
	repos, err := git.DetectRepositories(".")
	if err != nil {
		return fmt.Errorf("failed to detect repositories: %w", err)
	}

	if len(repos) == 0 {
		return fmt.Errorf("no git repositories found in current directory")
	}

	// Step 2: Check which repos have changes
	fmt.Println("📋 Checking for changes...")
	reposWithChanges := make([]string, 0)
	for _, repo := range repos {
		hasChanges, err := git.HasUncommittedChanges(repo)
		if err != nil {
			return fmt.Errorf("failed to check status of %s: %w", repo, err)
		}
		if hasChanges {
			reposWithChanges = append(reposWithChanges, repo)
		}
	}

	if len(reposWithChanges) == 0 {
		return fmt.Errorf("no repositories with uncommitted changes found")
	}

	fmt.Printf("✓ Found %d repo(s) with changes\n\n", len(reposWithChanges))

	// Step 3: Select repos to finish
	fmt.Println("📂 Select repositories to finalize:")
	repoModel := prompt.NewRepoSelectionModelWithDefaults(reposWithChanges, reposWithChanges)
	p := tea.NewProgram(repoModel)
	m, err := p.Run()
	if err != nil {
		return fmt.Errorf("error in repo selection: %w", err)
	}

	selectedRepos := m.(prompt.RepoSelectionModel).SelectedRepos
	if len(selectedRepos) == 0 {
		fmt.Println("❌ No repositories selected")
		os.Exit(0)
	}

	fmt.Printf("  ✓ Selected %d repo(s)\n\n", len(selectedRepos))

	// Step 4: Get commit message
	fmt.Println("📝 Enter commit message:")
	commitModel := prompt.NewCommitModel()
	p = tea.NewProgram(commitModel)
	m, err = p.Run()
	if err != nil {
		return fmt.Errorf("error getting commit message: %w", err)
	}

	commitResult := m.(prompt.CommitModel)
	if commitResult.Aborted || commitResult.Message == "" {
		fmt.Println("❌ Aborted by user")
		os.Exit(0)
	}

	commitMsg := commitResult.Message
	fmt.Printf("  ✓ Commit message: %s\n\n", commitMsg)

	// Step 5: Determine target branches for each repo
	repoTargets := make(map[string][]string)
	for _, repo := range selectedRepos {
		currentBranch, err := git.GetCurrentBranch(repo)
		if err != nil {
			return fmt.Errorf("failed to get current branch for %s: %w", repo, err)
		}

		// Check if branch exists on remote to determine origin
		originBranch, err := git.GetOriginBranch(repo, currentBranch)
		if err != nil {
			// If we can't determine, ask user
			fmt.Printf("\n⚠️  Could not determine origin branch for %s\n", repo)
			targetModel := prompt.NewTargetBranchModel(repo)
			p = tea.NewProgram(targetModel)
			m, err := p.Run()
			if err != nil {
				return fmt.Errorf("error selecting target branch: %w", err)
			}
			repoTargets[repo] = m.(prompt.TargetBranchModel).SelectedBranches
		} else {
			// Auto-determine based on origin
			if originBranch == "main" || originBranch == "master" {
				repoTargets[repo] = []string{"main", "dev"}
			} else {
				repoTargets[repo] = []string{"dev"}
			}
		}
	}

	// Step 6: Execute git operations for each repo
	fmt.Println("\n🚀 Finalizing repositories...")
	for _, repo := range selectedRepos {
		fmt.Printf("\n  📁 %s\n", repo)

		// Add all changes
		fmt.Printf("    → Adding changes...")
		if err := git.AddAll(repo); err != nil {
			return fmt.Errorf("failed to add changes in %s: %w", repo, err)
		}
		fmt.Println(" ✓")

		// Commit
		fmt.Printf("    → Committing...")
		if err := git.Commit(repo, commitMsg); err != nil {
			return fmt.Errorf("failed to commit in %s: %w", repo, err)
		}
		fmt.Println(" ✓")

		// Push
		currentBranch, _ := git.GetCurrentBranch(repo)
		fmt.Printf("    → Pushing %s...", currentBranch)
		if err := git.Push(repo, currentBranch); err != nil {
			return fmt.Errorf("failed to push in %s: %w", repo, err)
		}
		fmt.Println(" ✓")

		// Create PRs
		targets := repoTargets[repo]
		for _, target := range targets {
			fmt.Printf("    → Creating PR to %s...", target)
			if err := createPR(repo, currentBranch, target, commitMsg); err != nil {
				fmt.Printf(" ⚠️  %v\n", err)
			} else {
				fmt.Println(" ✓")
			}
		}
	}

	// Step 7: Return to original branches
	fmt.Println("\n🔄 Returning to original branches...")
	for _, repo := range selectedRepos {
		// Try to return to dev, fallback to main
		if err := git.CheckoutBranch(repo, "dev"); err != nil {
			if err := git.CheckoutBranch(repo, "main"); err != nil {
				fmt.Printf("  ⚠️  Could not return to base branch in %s\n", repo)
				continue
			}
		}
	}

	fmt.Println("\n✅ Development finalized successfully!")
	fmt.Printf("   Commit message: %s\n", commitMsg)
	fmt.Printf("   Repositories: %d\n", len(selectedRepos))

	return nil
}

func createPR(repo, branch, target, title string) error {
	// Check if PR already exists
	checkCmd := exec.Command("gh", "pr", "list",
		"--repo", getRepoFullName(repo),
		"--head", branch,
		"--base", target,
		"--json", "number",
	)
	output, err := checkCmd.Output()
	if err == nil && len(output) > 2 { // Not empty array
		return fmt.Errorf("PR already exists")
	}

	// Create PR
	body := fmt.Sprintf("This PR merges changes from branch `%s` into `%s`.", branch, target)
	cmd := exec.Command("gh", "pr", "create",
		"--repo", getRepoFullName(repo),
		"--head", branch,
		"--base", target,
		"--title", title,
		"--body", body,
	)

	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gh pr create failed: %s", string(output))
	}

	return nil
}

func getRepoFullName(repoPath string) string {
	// Get the remote URL and extract owner/repo
	cmd := exec.Command("git", "-C", repoPath, "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	url := strings.TrimSpace(string(output))
	// Handle both HTTPS and SSH formats
	url = strings.TrimPrefix(url, "https://github.com/")
	url = strings.TrimPrefix(url, "git@github.com:")
	url = strings.TrimSuffix(url, ".git")

	return url
}
