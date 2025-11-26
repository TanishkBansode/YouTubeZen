package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Start HTTP server for YouTube embed proxy
	go startEmbedServer()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "YouTubeZen",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// startEmbedServer starts a local HTTP server to serve YouTube embeds
func startEmbedServer() {
	http.HandleFunc("/youtube-embed", func(w http.ResponseWriter, r *http.Request) {
		videoID := r.URL.Query().Get("v")
		if videoID == "" {
			http.Error(w, "Missing video ID", http.StatusBadRequest)
			return
		}

		// Set proper headers
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Serve the embed HTML
		html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="referrer" content="strict-origin-when-cross-origin">
    <style>
        * { margin: 0; padding: 0; }
        html, body { height: 100%%; background: #000; }
        iframe { width: 100%%; height: 100%%; border: none; }
    </style>
</head>
<body>
    <iframe 
        src="https://www.youtube-nocookie.com/embed/%s"
        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
        allowfullscreen
        loading="lazy"
    ></iframe>
</body>
</html>`, videoID)

		w.Write([]byte(html))
	})

	fmt.Println("YouTube embed server starting on http://localhost:8888")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		fmt.Printf("Failed to start embed server: %v\n", err)
	}
}
