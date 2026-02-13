package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"gh-flow/internal/git"
	"gh-flow/internal/prompt"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Prepare repositories for development",
	Long:  `Detects git repositories, checks for uncommitted changes, and sets up branches for development.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStart()
	},
}

func runStart() error {
	// Step 1: Detect repositories
	fmt.Println("🔍 Scanning for git repositories...")
	repos, err := git.DetectRepositories(".")
	if err != nil {
		return fmt.Errorf("failed to detect repositories: %w", err)
	}

	if len(repos) == 0 {
		return fmt.Errorf("no git repositories found in current directory")
	}

	fmt.Printf("✓ Found %d repository(ies)\n\n", len(repos))

	// Step 2: Check for uncommitted changes
	fmt.Println("📋 Checking repository status...")
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

	if len(reposWithChanges) > 0 {
		fmt.Printf("⚠️  Found uncommitted changes in %d repo(s)\n", len(reposWithChanges))

		// Launch TUI for stash/abort decision
		model := prompt.NewStashModel(reposWithChanges)
		p := tea.NewProgram(model)
		m, err := p.Run()
		if err != nil {
			return fmt.Errorf("error in interactive prompt: %w", err)
		}

		result := m.(prompt.StashModel)
		if result.Aborted {
			fmt.Println("❌ Aborted by user")
			os.Exit(0)
		}

		// Stash changes if requested
		if result.ShouldStash {
			fmt.Println("\n📦 Stashing changes...")
			for _, repo := range result.SelectedRepos {
				if err := git.StashChanges(repo); err != nil {
					return fmt.Errorf("failed to stash changes in %s: %w", repo, err)
				}
				fmt.Printf("  ✓ Stashed: %s\n", repo)
			}
		}
	}

	// Step 3: Interactive selection of repos
	fmt.Println("\n📂 Select repositories for development:")
	repoModel := prompt.NewRepoSelectionModel(repos)
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

	// Step 4: Ask about urgency
	fmt.Println("⚡ Select urgency level:")
	urgencyModel := prompt.NewUrgencyModel()
	p = tea.NewProgram(urgencyModel)
	m, err = p.Run()
	if err != nil {
		return fmt.Errorf("error in urgency selection: %w", err)
	}

	urgency := m.(prompt.UrgencyModel).Selected
	baseBranch := "dev"
	if urgency == "urgent" {
		baseBranch = "main"
	}
	fmt.Printf("  ✓ Base branch: %s\n\n", baseBranch)

	// Step 5: Ask for branch type and name
	fmt.Println("🌿 Configure new branch:")
	branchModel := prompt.NewBranchModel()
	p = tea.NewProgram(branchModel)
	m, err = p.Run()
	if err != nil {
		return fmt.Errorf("error in branch configuration: %w", err)
	}

	branchResult := m.(prompt.BranchModel)
	if branchResult.Aborted {
		fmt.Println("❌ Aborted by user")
		os.Exit(0)
	}

	newBranch := fmt.Sprintf("%s/%s", branchResult.BranchType, branchResult.BranchName)
	fmt.Printf("  ✓ New branch: %s\n\n", newBranch)

	// Step 6: Execute git operations
	fmt.Println("🚀 Preparing repositories...")
	for _, repo := range selectedRepos {
		fmt.Printf("\n  📁 %s\n", repo)

		// Checkout base branch
		fmt.Printf("    → Checking out %s...", baseBranch)
		if err := git.CheckoutBranch(repo, baseBranch); err != nil {
			return fmt.Errorf("failed to checkout %s in %s: %w", baseBranch, repo, err)
		}
		fmt.Println(" ✓")

		// Pull latest changes
		fmt.Printf("    → Pulling latest changes...")
		if err := git.Pull(repo); err != nil {
			return fmt.Errorf("failed to pull in %s: %w", repo, err)
		}
		fmt.Println(" ✓")

		// Create new branch
		fmt.Printf("    → Creating branch %s...", newBranch)
		if err := git.CreateBranch(repo, newBranch); err != nil {
			return fmt.Errorf("failed to create branch in %s: %w", repo, err)
		}
		fmt.Println(" ✓")
	}

	fmt.Println("\n✅ All repositories prepared successfully!")
	fmt.Printf("   Base branch: %s\n", baseBranch)
	fmt.Printf("   New branch: %s\n", newBranch)
	fmt.Printf("   Repositories: %d\n", len(selectedRepos))

	return nil
}
