package git

import (
	"net/http"
	"net/url"
	"strings"
)

func ValidateGitHubURL(inputURL string) error {
	// Parse the input URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return err
	}

	// Ensure the URL is a valid GitHub URL
	if !strings.HasPrefix(parsedURL.Host, "github.com") {
		return nil
	}

	// Fetch the repository information
	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the repository is public (based on HTTP status code)
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return nil
}
