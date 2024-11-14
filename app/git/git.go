package git

import (
	"regexp"
	"strings"
	"time"

	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/permafrost-dev/git-ninja/app/utils"
)

type BranchInfo struct {
	Name            string
	CheckoutCount   int
	CommitCount     int
	CheckoutHistory []*BranchCheckoutInfo
	CheckedOutLast  time.Time
	Score           float64
}

type BranchCheckoutInfo struct {
	BranchName   string
	RelativeTime string
	Timestamp    time.Time
}

func (b *BranchInfo) UpdateCheckoutCount() {
	logitems, _ := GetRefLogItemsForBranch(b.Name)

	b.CommitCount = 0
	for _, item := range logitems {
		if item.Action == "commit" {
			b.CommitCount++
		}
	}
}

func (b *BranchInfo) UpdateScore() {
	ageInHours := time.Since(b.CheckedOutLast).Hours() * -1
	decay := 1000 - ageInHours // the older the checkout, the less it counts

	b.Score = (float64)(b.CheckoutCount*b.CommitCount) * decay

	if b.CommitCount == 0 {
		b.Score = b.Score / 3.5
	}
}

func (b *BranchInfo) Update() {
	b.UpdateCheckoutCount()
	b.UpdateScore()
}

func GetBranchInfoFromReflogLine(pattern *regexp.Regexp, reflogLine string, minMatchCount int) *BranchCheckoutInfo {
	matches := pattern.FindStringSubmatch(reflogLine)

	if matches == nil || len(matches) < minMatchCount {
		return nil
	}

	return &BranchCheckoutInfo{
		BranchName:   strings.TrimSpace(matches[3]),
		RelativeTime: strings.TrimSpace(matches[4]),
		Timestamp:    utils.ParseTimestampIntoTime(matches[1]),
	}
}

func GetGitReflogLines(prettyFmt string) ([]string, error) {
	output, err := helpers.RunCommandBuffered("git", "reflog", "show", "--pretty=format:'"+prettyFmt+"'", "--date=relative")
	if err != nil {
		return make([]string, 0), err
	}

	return strings.Split(output, "\n"), nil
}

func GetBranchReflogLines(name string, prettyFmt string) ([]string, error) {
	output, err := helpers.RunCommandBuffered("git", "reflog", "--pretty=format:"+prettyFmt, name)
	if err != nil {
		return make([]string, 0), err
	}

	return strings.Split(output, "\n"), nil
}

func SliceContainsBranchCommitData(slice []*BranchCheckoutInfo, info *BranchCheckoutInfo) bool {
	for _, item := range slice {
		if item.BranchName == info.BranchName {
			return true
		}
	}

	return false
}

type RefLogItem struct {
	BranchName  string
	CommitHash  string
	AuthorName  string
	AuthorEmail string
	Timestamp   time.Time
	Action      string
	Message     string
}

func GetRefLogItemsForBranch(branchName string) ([]*RefLogItem, error) {
	lines, err := GetBranchReflogLines(branchName, "%at|%H|%an|%ae|%gs|%gd")
	if err != nil {
		return make([]*RefLogItem, 0), err
	}

	result := make([]*RefLogItem, 0)

	for _, line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) < 6 {
			continue
		}

		timestamp := utils.ParseTimestampIntoTime(parts[0])
		commitHash := parts[1]
		authorName := parts[2]
		authorEmail := parts[3]
		message := parts[4]
		refName := parts[5]

		messagePrefix := strings.Split(message, ":")[0]

		result = append(result, &RefLogItem{
			BranchName:  refName,
			CommitHash:  commitHash,
			AuthorName:  authorName,
			AuthorEmail: authorEmail,
			Timestamp:   timestamp,
			Action:      messagePrefix,
			Message:     message,
		})
	}

	return result, nil
}
