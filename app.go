package main

import (
	"YouTubeZen/backend"
	"context"
	"net/http"

	"google.golang.org/api/youtube/v3"
)

var client *http.Client
var service *youtube.Service

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
	scope := youtube.YoutubeReadonlyScope
	client = backend.GetClient(scope)
	service = backend.Initialize(client)
}

// Search calls the SearchYoutube function from the backend package
func (a *App) Search(query string) ([]map[string]string, error) {
	results := backend.SearchYouTube(service, query)
	return results, nil
}
