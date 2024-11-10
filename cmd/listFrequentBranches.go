package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// getRelativeTime takes a time.Time and returns a human-readable relative time string
func getRelativeTime(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration.Minutes() < 1:
		return "just now"
	case duration.Minutes() < 100:
		return fmt.Sprintf("%.0f min ago", duration.Minutes())
	case duration.Hours() < 1:
		return fmt.Sprintf("%.0f min ago", duration.Minutes())
	case duration.Hours() < 36:
		return fmt.Sprintf("%.0f hours ago", duration.Hours())
	case duration.Hours() < 24*7:
		return fmt.Sprintf("%.0f days ago", duration.Hours()/24)
	case duration.Hours() < 24*30:
		return fmt.Sprintf("%.0f weeks ago", duration.Hours()/(24*7))
	default:
		return fmt.Sprintf("%.0f months ago", duration.Hours()/(24*30))
	}
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "branch:freq",
		Short: "Show recently checked out branch names",
		Run: func(c *cobra.Command, args []string) {

			availableBranches, err := getAvailableBranches()
			if err != nil {
				fmt.Println("Error fetching branches:", err)
				return
			}

			cmd := exec.Command("git", "reflog", "show", "--pretty=format:%gs~%ci")
			var out bytes.Buffer
			cmd.Stdout = &out
			err = cmd.Run()
			if err != nil {
				fmt.Println("Error executing git command:", err)
				return
			}

			// Convert command output to string and split by new lines
			output := out.String()
			lines := strings.Split(output, "\n")

			// Prepare regular expressions to match 'checkout:' and extract the branch name and ISO date
			checkoutRegexp := regexp.MustCompile(`^checkout: moving from [^ ]+ to ([^ ]+)~(.*)$`)

			// Map to keep track of branch frequencies and recency
			type branchInfo struct {
				branch string
				count  int
				latest time.Time
			}
			branchData := make(map[string]branchInfo)

			// Time threshold for very recent checkouts (last 2 days)
			recentThreshold := time.Now().AddDate(0, 0, -3)

			// Iterate through each line and apply the regex filters
			for _, line := range lines {
				if matches := checkoutRegexp.FindStringSubmatch(line); len(matches) == 3 {
					branch := matches[1]
					isoDate := strings.TrimSpace(matches[2])

					checkoutDate, err := time.Parse("2006-01-02 15:04:05 -0700", isoDate)
					if err != nil {
						fmt.Println("Error parsing date:", err)
						continue
					}

					if branchExists(branch, availableBranches) {
						info := branchData[branch]
						info.count++

						// Update latest checkout date
						if info.latest.IsZero() || checkoutDate.After(info.latest) {
							info.latest = checkoutDate
						}

						branchData[branch] = info
					}
				}
			}

			// Separate branches into very recent and older groups
			var veryRecentBranches, otherBranches []branchInfo
			for branch, info := range branchData {
				if info.latest.After(recentThreshold) {
					veryRecentBranches = append(veryRecentBranches, branchInfo{branch, info.count, info.latest})
				} else {
					otherBranches = append(otherBranches, branchInfo{branch, info.count, info.latest})
				}
			}

			// Sort very recent branches by recency (latest first)
			sort.Slice(veryRecentBranches, func(i, j int) bool {
				return veryRecentBranches[i].latest.After(veryRecentBranches[j].latest)
			})

			// Sort other branches by frequency and then by recency
			sort.Slice(otherBranches, func(i, j int) bool {
				if otherBranches[i].count == otherBranches[j].count {
					return otherBranches[i].latest.After(otherBranches[j].latest)
				}
				return otherBranches[i].count > otherBranches[j].count
			})

			// Combine the very recent and other branches, prioritizing very recent ones
			displayedBranches := append(veryRecentBranches, otherBranches...)

			// Display the sorted branches, ensuring at least 20 are shown
			count := 0
			for _, branch := range displayedBranches {
				infoStr := fmt.Sprintf("%2d checkouts, %-15s", branch.count, getRelativeTime(branch.latest))
				fmt.Printf("  \033[33m%28s \033[37;1m %s\033[0m\n", infoStr, branch.branch)
				count++
				if count >= 20 {
					break
				}
			}
		},
	})
}
