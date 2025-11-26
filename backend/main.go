package backend

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func Initialize(client *http.Client) *youtube.Service {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Println("Error initializing YouTube service:", err)
		return nil
	}
	return service
}

func GetSubscriptions(service *youtube.Service) []map[string]string {
	subscriptions := make([]map[string]string, 0)

	subscriptionsCall := service.Subscriptions.List([]string{"snippet"}).Mine(true).MaxResults(50)
	subscriptionsResponse, err := subscriptionsCall.Do()
	if err != nil {
		fmt.Println("Error fetching subscriptions:", err)
		return nil
	}

	for _, item := range subscriptionsResponse.Items {
		subscription := map[string]string{
			"id":      item.Id,
			"title":   item.Snippet.Title,
			"channel": item.Snippet.ChannelId,
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions
}

// fetchVideoDetails is a helper function to fetch video details by IDs and avoid code duplication
func fetchVideoDetails(service *youtube.Service, videoIDs []string) []map[string]string {
	detailsCall := service.Videos.List([]string{"snippet", "contentDetails"}).Id(strings.Join(videoIDs, ","))
	detailsResponse, err := detailsCall.Do()
	if err != nil {
		fmt.Println("Error fetching video details:", err)
		return nil
	}

	videos := make([]map[string]string, 0, len(detailsResponse.Items))
	for _, item := range detailsResponse.Items {
		video := map[string]string{
			"id":       item.Id,
			"title":    item.Snippet.Title,
			"channel":  item.Snippet.ChannelTitle,
			"duration": formatDuration(item.ContentDetails.Duration),
		}
		videos = append(videos, video)
	}

	return videos
}

// Search YouTube using the API key and return video details
func SearchYouTube(service *youtube.Service, query string) []map[string]string {

	// Handle query, check if it's a URL or a search query
	query, errorcode := HandleQuery(query)
	switch errorcode {
	// If int is 2, the query is not a valid URL.
	case 2:
		// Search for the top 10 videos based on the query
		searchCall := service.Search.List([]string{"id", "snippet"}).Q(query).MaxResults(10).Type("video")
		searchResponse, err := searchCall.Do()
		if err != nil {
			fmt.Println("Error searching YouTube:", err)
			return nil
		}

		// Collect video IDs for content details request
		var videoIDs []string
		for _, item := range searchResponse.Items {
			videoIDs = append(videoIDs, item.Id.VideoId)
		}

		// Fetch additional details (like duration) using the video IDs
		videos := fetchVideoDetails(service, videoIDs)
		if len(videos) == 0 {
			return nil
		}

		return videos
	// If int is 0, the query is a valid URL and the video ID was extracted.
	case 0:
		videoIDs := []string{query}
		videos := fetchVideoDetails(service, videoIDs)
		if len(videos) == 0 {
			return nil
		}

		return videos
	case 1: // Invalid URL
		return nil
	}
	return nil
}

// formatDuration converts ISO 8601 duration format (from YouTube API) to H:MM:SS or MM:SS format
// Examples: PT1H2M3S -> 1:02:03, PT15M30S -> 15:30, PT45S -> 0:45
func formatDuration(duration string) string {
	// YouTube uses ISO 8601 duration format: PT#H#M#S
	// Remove the PT prefix
	duration = strings.TrimPrefix(duration, "PT")

	var hours, minutes, seconds int

	// Use regex to extract hours, minutes, and seconds
	hoursRegex := regexp.MustCompile(`(\d+)H`)
	minutesRegex := regexp.MustCompile(`(\d+)M`)
	secondsRegex := regexp.MustCompile(`(\d+)S`)

	if matches := hoursRegex.FindStringSubmatch(duration); len(matches) > 1 {
		hours, _ = strconv.Atoi(matches[1])
	}
	if matches := minutesRegex.FindStringSubmatch(duration); len(matches) > 1 {
		minutes, _ = strconv.Atoi(matches[1])
	}
	if matches := secondsRegex.FindStringSubmatch(duration); len(matches) > 1 {
		seconds, _ = strconv.Atoi(matches[1])
	}

	// Format output based on presence of hours
	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
