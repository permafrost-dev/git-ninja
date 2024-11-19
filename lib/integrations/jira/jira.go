package jira

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// JiraSearchResponse represents the structure of Jira's search API response
type JiraSearchResponse struct {
	Issues []struct {
		ID    string `json:"id"`
		Key   string `json:"key"`
		Title string `json:"fields"` // Assuming 'fields' contains 'summary' or similar for title
	} `json:"issues"`
}

// IssueCache represents the structure of the cache file
type IssueCache struct {
	Timestamp time.Time `json:"timestamp"`
	JiraHash  string    `json:"jira_hash"`
	IssueIDs  []string  `json:"issue_ids"`
}

var (
	// Mutex to ensure thread-safe access to the cache file
	cacheMutex sync.Mutex
)

func GetJiraCacheFileName() string {
	cwd, err := os.UserHomeDir()

	if err != nil {
		cwd = os.Getenv("HOME")
	}

	cacheFileName := cwd + "/.gitninja.jira-cache.json"

	return cacheFileName
}

// GetCurrentUserActiveIssueIDs queries Jira for the current user's active issues and returns their IDs.
// It caches the response in cache.json and only fetches new data if the cache is older than 5 minutes.
func GetCurrentUserActiveIssueIDs(jiraBaseURL, email, apiToken string) ([]string, error) {

	const cacheDuration = 5 * time.Minute
	cacheFileName := GetJiraCacheFileName()

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Attempt to read from cache
	cachedData, err := readCache(cacheFileName)
	if err == nil {
		if getCurrentJiraHash() != cachedData.JiraHash {
			// If the Jira hash has changed, skip returning the cached data
		} else if time.Since(cachedData.Timestamp) < cacheDuration {
			return cachedData.IssueIDs, nil
		}
	}

	// If cache is invalid or reading failed, fetch new data
	issueIDs, err := fetchJiraIssues(jiraBaseURL, email, apiToken)
	if err != nil {
		// If fetching fails but cache exists, return cached data
		if cachedData != nil {
			// fmt.Printf("Warning: Failed to fetch new data: %v. Returning cached data.\n", err)
			return cachedData.IssueIDs, nil
		}
		return nil, err
	}

	// Update cache with new data
	err = writeCache(cacheFileName, issueIDs)
	if err != nil {
		// fmt.Printf("Warning: Failed to write cache: %v\n", err)
		// Proceed without failing, as we have the issue IDs
	}

	return issueIDs, nil
}

// fetchJiraIssues performs the HTTP request to Jira and retrieves the active issue IDs.
func fetchJiraIssues(jiraBaseURL, email, apiToken string) ([]string, error) {
	// Define the JQL query
	// This JQL fetches issues assigned to the current user that are not resolved or closed.
	jql := `assignee = currentUser()  AND sprint in openSprints() AND
	statusCategory IN ("To Do", "In Progress")
	ORDER BY updated DESC`

	// Prepare the request URL
	// Jira's search API endpoint
	searchURL, err := url.Parse(fmt.Sprintf("%s/rest/api/3/search", jiraBaseURL))
	if err != nil {
		return nil, fmt.Errorf("invalid Jira base URL: %v", err)
	}

	// Set query parameters
	query := searchURL.Query()
	query.Set("jql", jql)
	query.Set("fields", "id,key")  // We need the issue ID and Key
	query.Set("maxResults", "100") // Adjust as needed
	searchURL.RawQuery = query.Encode()

	// Create a new HTTP request with context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set authentication headers using Basic Auth (email and API token)
	auth := email + ":" + apiToken
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+encodedAuth)
	req.Header.Set("Accept", "application/json")

	// Initialize HTTP client
	client := &http.Client{}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Jira API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON response
	var searchResult JiraSearchResponse
	if err := json.Unmarshal(body, &searchResult); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Extract issue IDs
	var issueIDs []string
	for _, issue := range searchResult.Issues {
		issueIDs = append(issueIDs, issue.Key) // Using Key instead of ID for readability
	}

	return issueIDs, nil
}

func getCurrentJiraHash() string {
	jiraVars := []string{
		os.Getenv("JIRA_SUBDOMAIN"),
		os.Getenv("JIRA_EMAIL_ADDRESS"),
		os.Getenv("JIRA_API_TOKEN"),
	}

	return fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(jiraVars, ""))))
}

// readCache reads the cache file and returns the cached data.
// Returns an error if the file doesn't exist or is invalid.
func readCache(cacheFileName string) (*IssueCache, error) {
	filePath, err := filepath.Abs(cacheFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to determine absolute path for cache file: %v", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err // File doesn't exist or cannot be opened
	}
	defer file.Close()

	var cache IssueCache
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cache); err != nil {
		return nil, fmt.Errorf("failed to decode cache file: %v", err)
	}

	return &cache, nil
}

// writeCache writes the issue IDs and current timestamp to the cache file.
func writeCache(cacheFileName string, issueIDs []string) error {
	filePath, err := filepath.Abs(cacheFileName)
	if err != nil {
		return fmt.Errorf("failed to determine absolute path for cache file: %v", err)
	}

	cache := IssueCache{
		Timestamp: time.Now(),
		JiraHash:  getCurrentJiraHash(),
		IssueIDs:  issueIDs,
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create cache file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // For readability
	if err := encoder.Encode(&cache); err != nil {
		return fmt.Errorf("failed to encode cache data: %v", err)
	}

	return nil
}

// GetJiraTicketIDs is an example usage of GetCurrentUserActiveIssueIDs.
// It fetches and prints the active Jira issue IDs assigned to the current user.
func GetJiraTicketIDs(subdomain string, email string) []string {
	jiraBaseURL := "https://" + subdomain + ".atlassian.net"
	apiToken := os.Getenv("JIRA_API_TOKEN")

	if apiToken == "" {
		fmt.Println("Error: JIRA_API_TOKEN environment variable is not set.")
		return []string{}
	}

	issueIDs, err := GetCurrentUserActiveIssueIDs(jiraBaseURL, email, apiToken)
	if err != nil {
		// fmt.Printf("Error fetching issues: %v\n", err)
		return []string{}
	}

	if len(issueIDs) == 0 {
		// fmt.Println("No active issues assigned to the current user.")
		return []string{}
	}

	return issueIDs
}

func HashJiraIssueKey(issueKey string) (int64, error) {
	if issueKey == "" {
		return 0, errors.New("issue key is empty")
	}

	// Find the last hyphen to handle keys with multiple hyphens in the prefix
	lastHyphen := strings.LastIndex(issueKey, "-")
	if lastHyphen == -1 || lastHyphen == len(issueKey)-1 {
		return 0, fmt.Errorf("invalid issue key format: %s", issueKey)
	}

	// Extract the numeric part after the last hyphen
	numPart := issueKey[lastHyphen+1:]
	num, err := strconv.ParseInt(numPart, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric part in issue key '%s': %v", issueKey, err)
	}

	return num, nil
}
