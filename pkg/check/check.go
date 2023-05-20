package check

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetKepYaml(url string) (string, error) {
	// Convert the GitHub URL to the raw URL format
	rawURL := convertToRawURL(url)

	// Fetch the raw content of the code file
	content, err := fetchRawContent(rawURL)
	if err != nil {
		return "", fmt.Errorf("Failed to fetch raw content: %s\n", err)
	}

	return content, nil
}

// Convert a GitHub URL to its raw URL format
func convertToRawURL(url string) string {
	// Replace the 'github.com' domain with 'raw.githubusercontent.com'
	rawURL := strings.Replace(url, "github.com", "raw.githubusercontent.com", 1)

	// Remove the '/blob' segment from the URL path
	rawURL = strings.Replace(rawURL, "/blob", "", 1)

	return rawURL
}

// Fetch the raw content of a given URL
func fetchRawContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
