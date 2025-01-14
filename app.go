package main

import (
	"YouTubeZen/backend"
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	backend.Backend()
}

// Search calls the SearchYoutube function from the backend package
func (a *App) Search(query string) ([]map[string]string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get API key from environment variables
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		log.Fatal("YouTube API key not found in environment")
	}

	results := backend.SearchYouTube(apiKey, query)
	return results, nil
}
