package git

import (
	"net/http"
	"net/url"
	"strings"
)

func ValidateGitHubURL(inputURL string) (bool, error) {
	// Parse the input URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return false, err
	}

	// Ensure the URL is a valid GitHub URL
	if !strings.HasPrefix(parsedURL.Host, "github.com") {
		return false, nil
	}

	// Fetch the repository information
	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Check if the repository is public (based on HTTP status code)
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, nil
}
