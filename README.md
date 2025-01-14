# YouTube Zen

It is minimalist desktop app to view YouTube videos without the recommendations, comments, promotions;/n Just a pure video-viewing app.

## Why?
- It is not an extension, hence you can't disable it to see something which you weren't going to./n
- It is not a website, so you can't just switch tabs instantly for something new./n
- It treats shorts like videos, so you won't get trapped in doomscrolling./n

## Technologies used
- Golang
- [Wails](https://github.com/wailsapp/wails) 
- Svelte
- Typescript

## How to Use?
Download latest executable for your operating system from the Releases. 
Create a .env file in same directory as the executable and put your Youtube API Key:
```
YOUTUBE_API_KEY=API_KEY
```

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

## This project is not made by Google or Youtube ofc, and not a alternative of YouTube, it is just an app for people who don't want to quit Youtube but are addicted to it(like me).
