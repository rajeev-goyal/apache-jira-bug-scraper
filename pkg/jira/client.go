// In file: pkg/jira/client.go
package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Define the structure of the JSON response we expect from JIRA
// We only need the 'issues' array and the 'key' from each issue
type SearchResult struct {
	Issues []Issue `json:"issues"`
}

type Issue struct {
	Key string `json:"key"`
}

// Client holds the HTTP client and JIRA URL
type Client struct {
	BaseURL    *url.URL
	HttpClient *http.Client
}

// NewClient creates a new JIRA client
func NewClient() *Client {
	baseURL, err := url.Parse("https://issues.apache.org/jira/")
	if err != nil {
		// This panic is okay because the URL is a static constant
		log.Fatalf("Failed to parse base URL: %v", err)
	}

	return &Client{
		BaseURL:    baseURL,
		HttpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// SearchBugs performs the JQL query and returns a list of bug IDs (e.g., "KAFKA-1234")
// We will add authentication here later if needed.
func (c *Client) SearchBugs(projectName string, maxResults int) ([]string, error) {
	// 1. Construct the JQL query
	// project = KAFKA AND issuetype = Bug AND status in (Resolved, Fixed) ORDER BY updated DESC
	jql := fmt.Sprintf(
		"project = %s AND issuetype = Bug AND status in (Resolved, Fixed) ORDER BY updated DESC",
		projectName,
	)

	// 2. Build the full request URL
	relURL, _ := url.Parse("rest/api/2/search")
	reqURL := c.BaseURL.ResolveReference(relURL)

	// Add JQL query parameters
	params := url.Values{}
	params.Add("jql", jql)
	params.Add("fields", "key") // We only need the issue key
	params.Add("maxResults", fmt.Sprintf("%d", maxResults))
	reqURL.RawQuery = params.Encode()

	fmt.Printf("Querying JIRA: %s\n", reqURL.String())

	// 3. Create the HTTP request
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 4. Send the request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 5. Check for errors
	if resp.StatusCode != http.StatusOK {
		// Read the body for more error info
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("JIRA request failed with status %s: %s", resp.Status, string(bodyBytes))
	}

	// 6. Parse the JSON response
	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode JIRA response: %w", err)
	}

	// 7. Extract just the keys into a string slice
	var bugKeys []string
	for _, issue := range result.Issues {
		bugKeys = append(bugKeys, issue.Key)
	}

	return bugKeys, nil
}
