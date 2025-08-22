package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/andrerfcsantos/arbor/lib"
	"github.com/spf13/cobra"
)

var locCmd = &cobra.Command{
	Use:   "loc [repository-path]",
	Short: "Analyze lines of code across commits",
	Long: `Analyze the lines of code for each language across all commits in a Git repository.
If no repository path is provided, the current directory is assumed to be a Git repository.

The command will:
1. Open the specified Git repository
2. Iterate through all commits chronologically
3. Checkout each commit and count lines of code for each language using the scc tool
4. Generate a chart showing the evolution of lines of code over time
5. Restore the repository to its original checkout state`,
	Args: cobra.MaximumNArgs(1),
	RunE: runLocCommand,
}

func init() {
	// Add flags here if needed
}

func runLocCommand(cmd *cobra.Command, args []string) error {
	var repoPath string
	if len(args) > 0 {
		repoPath = args[0]
	} else {
		var err error
		repoPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Open the repository
	repo, err := lib.OpenRepository(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	// Store original state info for user feedback
	originalBranch := repo.GetOriginalBranch()
	fmt.Printf("📁 Analyzing repository: %s\n", repoPath)
	fmt.Printf("📍 Original checkout: %s\n", originalBranch)
	fmt.Println()

	// Get all commits
	commits, err := repo.GetCommits()
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	fmt.Printf("🔍 Found %d commits to analyze\n", len(commits))
	fmt.Println()

	// Track languages found across all commits
	languageSet := make(map[string]bool)

	// Analyze each commit
	for i, commit := range commits {

		// Checkout this commit
		err = repo.CheckoutCommit(commit.Hash)
		if err != nil {
			fmt.Printf("   ❌ Failed to checkout commit: %v\n", err)
			continue
		}

		// Count lines of code for this commit
		locData, err := lib.CountLinesOfCode(repoPath)
		if err != nil {
			fmt.Printf("   ❌ Failed to count lines of code: %v\n", err)
			continue
		}

		// Track languages found
		for lang := range locData {
			languageSet[lang] = true
		}

		// Update commit data with language information
		commits[i].Languages = locData

		// Calculate progress percentage
		progress := float64(i+1) / float64(len(commits)) * 100

		fmt.Printf("📊 Processed commit %d/%d (%.1f%%): %s '%s' by %s\n",
			i+1, len(commits), progress, commit.Hash[:8], strings.TrimSpace(commit.Message), commit.Author)
		fmt.Println()
	}

	// Sort commits by date
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.Before(commits[j].Date)
	})

	fmt.Println("🔄 Restoring original repository state...")

	// Restore the original checkout state
	err = repo.RestoreOriginalState()
	if err != nil {
		return fmt.Errorf("failed to restore original state: %w", err)
	}

	fmt.Printf("✅ Restored to: %s\n", originalBranch)
	fmt.Println()

	// Generate chart
	fmt.Println("📊 Generating chart...")
	err = lib.GenerateLOCChart(commits, languageSet, "loc_analysis.html")
	if err != nil {
		return fmt.Errorf("failed to generate chart: %w", err)
	}

	fmt.Println("🎉 Analysis complete!")
	fmt.Printf("📁 Chart saved as: loc_analysis.html\n")
	fmt.Printf("📊 Analyzed %d commits across %d languages\n", len(commits), len(languageSet))

	return nil
}
