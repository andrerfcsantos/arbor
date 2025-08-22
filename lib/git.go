package lib

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Repository represents a git repository with analysis capabilities
type Repository struct {
	repo           *git.Repository
	path           string
	originalRef    *plumbing.Reference
	originalBranch string
}

// CommitData represents commit information with language analysis
type CommitData struct {
	Hash      string
	Message   string
	Author    string
	Date      time.Time
	Languages map[string]int
}

// OpenRepository opens a git repository at the specified path
func OpenRepository(path string) (*Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository at %s: %w", path, err)
	}

	// Get the current HEAD reference to restore later
	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository head: %w", err)
	}

	// Get the current branch name
	var branchName string
	if head.Name().IsBranch() {
		branchName = head.Name().String()
	} else {
		branchName = head.Hash().String()[:8]
	}

	return &Repository{
		repo:           repo,
		path:           path,
		originalRef:    head,
		originalBranch: branchName,
	}, nil
}

// GetCommits retrieves all commits from the repository
func (r *Repository) GetCommits() ([]CommitData, error) {
	// Get the repository head
	head, err := r.repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository head: %w", err)
	}

	// Get commit iterator
	commitIter, err := r.repo.Log(&git.LogOptions{
		From:  head.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}
	defer commitIter.Close()

	// Collect commit data
	var commits []CommitData

	err = commitIter.ForEach(func(commit *object.Commit) error {
		commitData := CommitData{
			Hash:      commit.Hash.String(),
			Message:   commit.Message,
			Author:    commit.Author.Name,
			Date:      commit.Author.When,
			Languages: make(map[string]int), // Will be populated later
		}

		commits = append(commits, commitData)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to process commits: %w", err)
	}

	return commits, nil
}

// CheckoutCommit checks out a specific commit
func (r *Repository) CheckoutCommit(commitHash string) error {
	worktree, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	hash := plumbing.NewHash(commitHash)
	err = worktree.Checkout(&git.CheckoutOptions{
		Hash:  hash,
		Force: true,
	})
	if err != nil {
		return fmt.Errorf("failed to checkout commit %s: %w", commitHash, err)
	}

	return nil
}

// RestoreOriginalState restores the repository to its original checkout state
func (r *Repository) RestoreOriginalState() error {
	worktree, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Try to checkout the original branch first
	if r.originalRef.Name().IsBranch() {
		err = worktree.Checkout(&git.CheckoutOptions{
			Branch: r.originalRef.Name(),
		})
		if err != nil {
			// If branch checkout fails, try to checkout the specific commit
			err = worktree.Checkout(&git.CheckoutOptions{
				Hash: r.originalRef.Hash(),
			})
			if err != nil {
				return fmt.Errorf("failed to restore original state: %w", err)
			}
		}
	} else {
		// If it was a detached HEAD, checkout the specific commit
		err = worktree.Checkout(&git.CheckoutOptions{
			Hash: r.originalRef.Hash(),
		})
		if err != nil {
			return fmt.Errorf("failed to restore original state: %w", err)
		}
	}

	return nil
}

// GetPath returns the repository path
func (r *Repository) GetPath() string {
	return r.path
}

// GetOriginalBranch returns the original branch name
func (r *Repository) GetOriginalBranch() string {
	return r.originalBranch
}
