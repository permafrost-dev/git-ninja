package utils

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func RunCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	result, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(result)), nil
}

func GetRelativeTime(t time.Time) string {
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
