## YouTube Zen

It is a minimalist desktop app to view YouTube videos without the recommendations, comments, or promotions. Just a pure video-viewing experience.

[![Watch the video](https://img.youtube.com/vi/2xioDUufMu0/0.jpg)](https://www.youtube.com/watch?v=2xioDUufMu0)
## Why?
- It is not an extension, so you can't disable it and accidentally see things you weren't intending to.
- It is not a website, so you can't just switch tabs instantly for something else.
- It treats shorts like videos, so you won’t get trapped in doomscrolling.

## Technologies Used
- Golang
- [Wails](https://github.com/wailsapp/wails)
- Svelte
- Typescript

## How to Use?
Download the latest executable for your operating system from the [Releases](https://github.com/your-repo/releases) section. 
Add `client_secret.json` file in the same directory as the executable, get it from the google cloud console, and then after running, it'll get you to the login page.


## Live Development

To run in live development mode, run `wails dev` in the project directory. This will launch a Vite development server, providing very fast hot reloads for your frontend changes. If you want to develop in a browser and access your Go methods, there is also a dev server available at [http://localhost:34115](http://localhost:34115). Connect to this in your browser, and you can call your Go code from the devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

## Disclaimer
This project is not made by Google or YouTube, and it is not an alternative to YouTube. It is simply an app for people who want to keep using YouTube without getting distracted or addicted, like I did.
