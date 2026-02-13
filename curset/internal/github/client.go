package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	collectionURL = "https://raw.githubusercontent.com/bilgehannal/cursor-config/refs/heads/main/data/collection.json"
	contentsAPI   = "https://api.github.com/repos/bilgehannal/cursor-config/contents/data/.cursor"
	rawBaseURL    = "https://raw.githubusercontent.com/bilgehannal/cursor-config/refs/heads/main/data/.cursor"
)

// ContentEntry represents a single entry from the GitHub Contents API.
type ContentEntry struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"` // "file" or "dir"
	DownloadURL string `json:"download_url"`
}

// ContentsResult holds the result of a GitHub Contents API call.
type ContentsResult struct {
	Entries []ContentEntry
	IsDir   bool // true if the API returned an array (directory listing)
}

// Client is an HTTP client for fetching data from GitHub.
type Client struct {
	httpClient *http.Client
	cache      map[string]*ContentsResult // cache for ListContents results
}

// NewClient creates a new GitHub client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		cache:      make(map[string]*ContentsResult),
	}
}

// FetchCollectionJSON fetches the collection.json from the raw GitHub URL and returns the bytes.
func (c *Client) FetchCollectionJSON() ([]byte, error) {
	resp, err := c.httpClient.Get(collectionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch collection.json: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch collection.json: HTTP %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// ListContents lists the contents of a path under data/.cursor/ using the GitHub Contents API.
// For example, path "rules/common" lists files in data/.cursor/rules/common/.
// Returns ContentsResult which indicates whether the path is a directory or a file.
// Results are cached to avoid redundant API calls.
func (c *Client) ListContents(path string) (*ContentsResult, error) {
	// Check cache first.
	if cached, ok := c.cache[path]; ok {
		return cached, nil
	}

	url := fmt.Sprintf("%s/%s", contentsAPI, path)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to list contents at %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("path not found: %s", path)
	}

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("GitHub API rate limit exceeded. Try again later or use a GitHub token")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list contents at %s: HTTP %d", path, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// The API returns an array for directories, or a single object for files.
	// Try array first.
	var entries []ContentEntry
	if err := json.Unmarshal(body, &entries); err == nil {
		result := &ContentsResult{Entries: entries, IsDir: true}
		c.cache[path] = result
		return result, nil
	}

	// Try single object (it's a file, not a directory).
	var single ContentEntry
	if err := json.Unmarshal(body, &single); err != nil {
		return nil, fmt.Errorf("failed to parse contents response: %w", err)
	}

	result := &ContentsResult{Entries: []ContentEntry{single}, IsDir: false}
	c.cache[path] = result
	return result, nil
}

// DownloadFile downloads a raw file from the repository.
// The filePath is relative to data/.cursor/, e.g. "rules/common/clean-code.mdc".
func (c *Client) DownloadFile(filePath string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", rawBaseURL, filePath)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %w", filePath, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download %s: HTTP %d", filePath, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
