package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/permafrost-dev/git-ninja/app/git"
	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/permafrost-dev/git-ninja/app/utils"
	"github.com/spf13/cobra"
)

func getGroupedAndSortedDisplayBranches(branchData map[string]git.BranchInfo, thresholds *FrequentBranchThresholds, limit int) []git.BranchInfo {
	var displayedBranches, veryRecentBranches, otherBranches []git.BranchInfo

	veryRecentBranches, otherBranches = splitIntoRecentAndOldBranches(branchData, thresholds.Recent, thresholds.Older)
	veryRecentBranches, otherBranches = sortRecentAndOtherBranches(veryRecentBranches, otherBranches, thresholds.Older)

	// Combine the very recent and other branches, prioritizing very recent ones
	displayedBranches = append(veryRecentBranches, otherBranches...)

	return displayedBranches[:limit]
}

func splitIntoRecentAndOldBranches(branchData map[string]git.BranchInfo, recentThreshold time.Time, oldThreshold time.Time) ([]git.BranchInfo, []git.BranchInfo) {
	var veryRecentBranches, otherBranches []git.BranchInfo

	for _, bd := range branchData {
		if bd.CheckedOutLast.After(recentThreshold) {
			veryRecentBranches = append(veryRecentBranches, bd) // git.BranchInfo{Name: branch, CheckoutCount: bd.CheckoutCount, CheckedOutLast: bd.CheckedOutLast, CommitCount: bd.CommitCount})
		} else {
			otherBranches = append(otherBranches, bd) // git.BranchInfo{Name: branch, CheckoutCount: bd.CheckoutCount, CheckedOutLast: bd.CheckedOutLast, CommitCount: bd.CommitCount})
		}
	}

	return veryRecentBranches, otherBranches
}

func sortRecentAndOtherBranches(recent []git.BranchInfo, other []git.BranchInfo, oldThreshold time.Time) ([]git.BranchInfo, []git.BranchInfo) {

	// Sort very recent branches by recency (latest first)
	sort.Slice(recent, func(i, j int) bool {
		if recent[i].Score == recent[j].Score {
			return recent[i].CheckedOutLast.After(recent[j].CheckedOutLast)
		}
		return recent[i].Score > recent[j].Score
	})

	// Sort other branches by frequency and then by recency
	sort.Slice(other, func(i, j int) bool {
		if other[i].CheckedOutLast.Before(oldThreshold) {
			if other[i].CheckoutCount == other[j].CheckoutCount {
				return other[i].CheckedOutLast.After(other[j].CheckedOutLast)
			}
			if other[i].CheckedOutLast.Before(other[j].CheckedOutLast) && other[i].CheckoutCount > other[j].CheckoutCount {
				return true
			}
			if other[i].CheckedOutLast.After(other[j].CheckedOutLast) && other[i].CheckoutCount < other[j].CheckoutCount {
				return false
			}

			return other[i].CheckoutCount > other[j].CheckoutCount
		}

		return other[i].CheckoutCount > other[j].CheckoutCount
	})

	return recent, other
}

type FrequentBranchThresholds struct {
	Recent time.Time
	Older  time.Time
}

func processRefLogLines(lines []string, existingBranches map[string]bool) map[string]git.BranchInfo {
	branches := make(map[string]git.BranchInfo)

	for _, line := range lines {
		info := git.GetBranchInfoFromReflogLine(lineRegex, line, 4)
		if info == nil || !utils.MapEntryExists(info.BranchName, existingBranches) {
			continue
		}

		data := git.BranchInfo{Name: info.BranchName, CheckoutCount: 0, CheckedOutLast: time.Time{}}

		if _, ok := branches[info.BranchName]; ok {
			data = branches[info.BranchName]
		}

		data.CheckoutCount++

		if data.CheckedOutLast.IsZero() || info.Timestamp.After(data.CheckedOutLast) {
			data.CheckedOutLast = info.Timestamp
		}

		branches[info.BranchName] = data
	}

	for name, branch := range branches {
		branch.Update()
		branches[name] = branch
	}

	return branches
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "branch:freq",
		Short: "Show recently checked out branch names",
		Run: func(c *cobra.Command, args []string) {
			lines, _ := git.GetGitReflogLines("%at ~ %gs ~ %gd")
			availableBranches, _ := helpers.GetAvailableBranchesMap()

			thresholds := FrequentBranchThresholds{
				Recent: time.Now().AddDate(0, 0, -7),
				Older:  time.Now().AddDate(0, 0, -15),
			}

			branches := processRefLogLines(lines, availableBranches)
			frequent := getGroupedAndSortedDisplayBranches(branches, &thresholds, 15)

			for _, br := range frequent {
				description := fmt.Sprintf("%2d checkouts, %2d commits, %-15s", br.CheckoutCount, br.CommitCount, utils.GetRelativeTime(br.CheckedOutLast))
				fmt.Printf("  \033[33m%28s \033[37;1m %s\033[0m\n", description, br.Name)
			}
		},
	})
}
