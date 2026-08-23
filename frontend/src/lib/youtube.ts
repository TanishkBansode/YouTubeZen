// Keyless YouTube data: URL parsing, search, and video details.
// Uses public Piped/Invidious instances and noembed.com (CORS-enabled
// proxy for YouTube's oEmbed endpoint).

export interface Video {
  id: string;
  title: string;
  channel: string;
  duration: string;
}

// Public mirrors go down or rate-limit regularly; every request tries
// these lists in order until one answers. Edit freely if your favourite
// instance stops working (https://github.com/TeamPiped/Piped/wiki/Instances,
// https://api.invidious.io).
const PIPED_INSTANCES = [
  "https://api.piped.private.coffee",
  "https://pipedapi.reallyaweso.me",
  "https://pipedapi.ducks.party",
];

const INVIDIOUS_INSTANCES = [
  "https://yewtu.be",
  "https://inv.nadeko.net",
  "https://invidious.nerdvpn.de",
];

const REQUEST_TIMEOUT_MS = 10_000;

// Matches YouTube video IDs in all common URL formats:
// watch?v=, /embed/, /v/, /shorts/, youtu.be/
const VIDEO_ID_PATTERN =
  /(?:youtube\.com\/(?:watch\?v=|embed\/|v\/|shorts\/)|youtu\.be\/)([a-zA-Z0-9_-]{11})/;

function isValidURL(input: string): boolean {
  try {
    const parsed = new URL(input);
    return parsed.protocol === "http:" || parsed.protocol === "https:";
  } catch {
    return false;
  }
}

export function extractVideoID(input: string): string | null {
  const match = VIDEO_ID_PATTERN.exec(input);
  return match ? match[1] : null;
}

export function formatDuration(totalSeconds: number): string {
  if (!totalSeconds || totalSeconds < 0) return "";
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;

  if (hours > 0) {
    return `${hours}:${String(minutes).padStart(2, "0")}:${String(seconds).padStart(2, "0")}`;
  }
  return `${minutes}:${String(seconds).padStart(2, "0")}`;
}

async function getJSON<T>(url: string): Promise<T | null> {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);
  try {
    const resp = await fetch(url, { signal: controller.signal });
    if (!resp.ok) return null;
    return (await resp.json()) as T;
  } catch {
    return null;
  } finally {
    clearTimeout(timer);
  }
}

interface InvidiousVideo {
  videoId?: string;
  title?: string;
  author?: string;
  lengthSeconds?: number;
}

interface PipedResponse {
  items?: Array<{
    url?: string;
    title?: string;
    uploaderName?: string;
    duration?: number;
  }>;
}

interface PipedStream {
  title?: string;
  uploader?: string;
  duration?: number;
}

// Same field names as YouTube's oEmbed; served by noembed.com with CORS.
interface OEmbedResponse {
  title?: string;
  author_name?: string;
}

async function searchInvidious(
  instance: string,
  query: string,
): Promise<Video[] | null> {
  const endpoint = `${instance}/api/v1/search?type=video&q=${encodeURIComponent(query)}`;
  const results = await getJSON<InvidiousVideo[]>(endpoint);
  if (!results || results.length === 0) return null;

  const videos: Video[] = [];
  for (const item of results.slice(0, 10)) {
    if (!item.videoId) continue;
    videos.push({
      id: item.videoId,
      title: item.title ?? "",
      channel: item.author ?? "",
      duration: formatDuration(item.lengthSeconds ?? 0),
    });
  }
  return videos.length > 0 ? videos : null;
}

async function searchPiped(
  instance: string,
  query: string,
): Promise<Video[] | null> {
  const endpoint = `${instance}/search?q=${encodeURIComponent(query)}&filter=videos`;
  const response = await getJSON<PipedResponse>(endpoint);
  if (!response?.items || response.items.length === 0) return null;

  const videos: Video[] = [];
  for (const item of response.items) {
    if (videos.length >= 10) break;
    const match = /v=([a-zA-Z0-9_-]{11})/.exec(item.url ?? "");
    if (!match) continue;
    videos.push({
      id: match[1],
      title: item.title ?? "",
      channel: item.uploaderName ?? "",
      duration: formatDuration(item.duration ?? 0),
    });
  }
  return videos.length > 0 ? videos : null;
}

async function searchVideos(query: string): Promise<Video[]> {
  for (const instance of PIPED_INSTANCES) {
    const videos = await searchPiped(instance, query);
    if (videos) return videos;
  }
  for (const instance of INVIDIOUS_INSTANCES) {
    const videos = await searchInvidious(instance, query);
    if (videos) return videos;
  }
  return [];
}

async function videoDetailsFromPiped(id: string): Promise<Video | null> {
  for (const instance of PIPED_INSTANCES) {
    const stream = await getJSON<PipedStream>(`${instance}/streams/${id}`);
    if (stream?.title) {
      return {
        id,
        title: stream.title,
        channel: stream.uploader ?? "",
        duration: formatDuration(stream.duration ?? 0),
      };
    }
  }
  return null;
}

async function videoDetailsFromInvidious(id: string): Promise<Video | null> {
  for (const instance of INVIDIOUS_INSTANCES) {
    const endpoint = `${instance}/api/v1/videos/${id}?fields=videoId,title,author,lengthSeconds`;
    const video = await getJSON<InvidiousVideo>(endpoint);
    if (video?.videoId) {
      return {
        id: video.videoId,
        title: video.title ?? "",
        channel: video.author ?? "",
        duration: formatDuration(video.lengthSeconds ?? 0),
      };
    }
  }
  return null;
}

// Keyless last-resort fallback via noembed.com (CORS-enabled proxy for
// YouTube's official oEmbed endpoint). Title and channel, no duration.
async function videoDetailsOEmbed(id: string): Promise<Video | null> {
  const endpoint =
    `https://noembed.com/embed?url=` +
    encodeURIComponent(`https://www.youtube.com/watch?v=${id}`);
  const meta = await getJSON<OEmbedResponse>(endpoint);
  if (!meta?.title) return null;
  return {
    id,
    title: meta.title,
    channel: meta.author_name ?? "",
    duration: "",
  };
}

async function videoDetails(id: string): Promise<Video | null> {
  return (
    (await videoDetailsFromPiped(id)) ??
    (await videoDetailsFromInvidious(id)) ??
    (await videoDetailsOEmbed(id))
  );
}

// Handles both pasted YouTube URLs and plain text queries. No API keys needed.
export async function searchYouTube(query: string): Promise<Video[]> {
  const trimmed = query.trim();
  if (isValidURL(trimmed)) {
    const id = extractVideoID(trimmed);
    if (!id) return [];
    const details = await videoDetails(id);
    return details ? [details] : [];
  }
  return searchVideos(trimmed);
}
