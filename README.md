## YouTube Zen

A minimalist, installable web app to watch YouTube videos without recommendations, comments, or promotions. Just a pure video-viewing experience.

- No API keys, no Google account, no login — ever
- Not an extension, so you can't accidentally disable it
- Not a native app with feeds, so there's nothing to doomscroll
- Treats Shorts links like normal videos
- Demo at: https://tanishkbansode.github.io/YouTubeZen/
## How it works

Everything runs in your browser:

- **Search** uses public, keyless [Invidious](https://invidious.io) and [Piped](https://github.com/TeamPiped/Piped) instances (with automatic failover)
- **Video details** for pasted URLs come from Invidious or YouTube's official [oEmbed](https://www.youtube.com/oembed) endpoint
- **Playback** uses the privacy-friendly `youtube-nocookie.com` embed

The whole backend from the old desktop version is now ~200 lines of TypeScript in `frontend/src/lib/youtube.ts`.

## Install on Android

1. Open the hosted URL in Chrome on your phone
2. Tap the menu (⋮) → **Add to Home screen** / **Install app**
3. Launch it like any other app — fullscreen, own icon, no browser UI

Works on desktop browsers too (install prompt in Chrome/Edge address bar).

## Development

Prerequisites: Node.js 16+ (npm)

```bash
cd frontend
npm install
npm run dev        # dev server with hot reload
```

To test on your phone during development:

```bash
npm run dev -- --host   # then open http://<your-lan-ip>:5173 on the phone
```

## Build & Host

```bash
cd frontend
npm run build      # static output in dist/
```

**GitHub Pages (automatic):** every push to `main` builds and deploys via the GitHub Action in `.github/workflows/deploy.yml`, live at `https://tanishkbansode.github.io/YouTubeZen/`.

All PWA paths are relative, so `dist/` also works on any other static host (Netlify, Cloudflare Pages) without config.

## Usage

### Searching

Enter a search term, or paste any YouTube URL format:

- `https://www.youtube.com/watch?v=VIDEO_ID`
- `https://youtu.be/VIDEO_ID`
- `https://www.youtube.com/embed/VIDEO_ID`
- `https://www.youtube.com/shorts/VIDEO_ID`

### Keyboard Shortcuts / Gestures

- `Escape` - Close video player or settings modal
- `Enter` - Submit search
- Click the title to reload the app

## Troubleshooting

**"No videos found" / search errors**
Public Invidious/Piped instances occasionally go down or rate-limit. The app tries several automatically; if all fail, try again later or edit the instance lists at the top of `frontend/src/lib/youtube.ts`.

**Video won't play**
Some videos disable embedding; open them directly on YouTube instead.

## Privacy & Security

- No accounts, no tokens, no tracking — nothing is stored except the PWA cache
- Requests go only to the Invidious/Piped instances listed in `youtube.ts` and to YouTube's oEmbed/embed endpoints

## Disclaimer

Not made by or affiliated with Google/YouTube. Just an app for people who want YouTube without getting distracted.

## License

MIT License.
