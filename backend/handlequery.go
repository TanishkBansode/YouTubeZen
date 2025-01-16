package backend

import (
	"fmt"
	"net/url"
	"regexp"
)

// isValidURL checks if the given string is a valid URL.
func isValidURL(input string) bool {
	parsedURL, err := url.ParseRequestURI(input)
	return err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

// extractVideoID extracts the video ID from a valid YouTube URL.
func extractVideoID(youtubeURL string) (string, error) {
	// Regex pattern to match YouTube video IDs
	pattern := `(https?://)?(www\.)?(youtube\.com/(watch\?v=|embed/|v/|shorts/)|youtu\.be/)([a-zA-Z0-9_-]{11})`
	re := regexp.MustCompile(pattern)

	// Match the regex with the given URL
	match := re.FindStringSubmatch(youtubeURL)
	if len(match) > 5 {
		return match[5], nil // Return the captured video ID
	}

	return "", fmt.Errorf("no valid YouTube video ID found in URL: %s", youtubeURL)
}

// HandleQuery checks if the query is a valid URL or a search query.
// If int is 0, the query is a valid URL and the video ID was extracted.
// If int is 1, the query is a URL but the video ID could not be extracted.
// If int is 2, the query is not a valid URL.
func HandleQuery(query string) (string, int) {
	if isValidURL(query) {
		videoID, err := extractVideoID(query)
		if err != nil {
			return query, 1
		} else {
			return videoID, 0
		}
	} else {
		return query, 2
	}
}
