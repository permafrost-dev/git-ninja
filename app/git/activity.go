package git

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/go-git/go-git/v5"
	g "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	object "github.com/go-git/go-git/v5/plumbing/object"
)

// BranchActivity represents a branch and its commit count in the last 48 hours
type BranchActivity struct {
	Name           string
	IsRemote       bool
	CommitCount    int
	LatestCommit   string
	LatestCommitAt time.Time
}

// GetActiveBranches returns a slice of BranchActivity sorted by commit count in descending order
func GetActiveBranches(repoPath string) ([]BranchActivity, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the list of all references (includes branches, tags, etc.)
	refs, err := repo.References()
	if err != nil {
		return nil, fmt.Errorf("failed to list references: %w", err)
	}

	// Define the cutoff time (48 hours ago)
	cutoffTime := time.Now().Add(-48 * time.Hour)

	// Slice to hold branch activities
	var activities []BranchActivity

	// Iterate through each reference
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		var isBranch bool
		var isRemote bool

		// Determine if the reference is a branch
		if ref.Name().IsBranch() {
			isBranch = true
			isRemote = false
		} else if ref.Name().IsRemote() {
			isBranch = true
			isRemote = true
		}

		if !isBranch {
			// Skip non-branch references
			return nil
		}

		branchName := ref.Name().Short()

		// Get the commit iterator for the branch
		commitIter, err := repo.Log(&g.LogOptions{From: ref.Hash()})
		if err != nil {
			// Skip branches that cannot be iterated
			log.Printf("warning: failed to get commits for branch %s: %v", branchName, err)
			return nil
		}

		defer commitIter.Close()

		// Count commits after cutoffTime
		count := 0
		var latestCommitMsg string
		var latestCommitTime time.Time

		err = commitIter.ForEach(func(c *object.Commit) error {
			if c.Committer.When.After(cutoffTime) {
				count++
				// Capture the latest commit details
				if latestCommitTime.IsZero() || c.Committer.When.After(latestCommitTime) {
					latestCommitTime = c.Committer.When
					latestCommitMsg = c.Message
				}
			} else {
				return fmt.Errorf("reached commits older than cutoff")
			}
			return nil
		})

		// Append the activity with latest commit details
		activities = append(activities, BranchActivity{
			Name:           branchName,
			IsRemote:       isRemote,
			CommitCount:    count,
			LatestCommit:   latestCommitMsg,
			LatestCommitAt: latestCommitTime,
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error iterating references: %w", err)
	}

	// Sort the activities by CommitCount in descending order
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].CommitCount > activities[j].CommitCount
	})

	activities = activities[:50]

	return activities, nil
}

func ShowActive() {
	// Example usage
	repoPath, _ := os.Getwd()

	activities, err := GetActiveBranches(repoPath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Branches sorted by activity in the last 48 hours:")
	for _, activity := range activities {
		branchType := "Local"
		if activity.IsRemote {
			branchType = "Remote"
		}
		if activity.CommitCount == 0 {
			continue
		}
		fmt.Printf("Branch: %s (%s), Commits: %d, Latest Commit at %s\n", activity.Name, branchType, activity.CommitCount, activity.LatestCommitAt.Format(time.RFC1123))
		// fmt.Printf("Branch: %s (%s), Commits: %d\n", activity.Name, branchType, activity.CommitCount)
	}
}
