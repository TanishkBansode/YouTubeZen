<script lang="ts">
  import { onMount } from "svelte";
  import { Search } from "../wailsjs/go/main/App.js";
  import "./styles/styles.css";
  import "./styles/theme-light.css";
  import "./styles/theme-dark.css";

  // TypeScript interface for Video data
  interface Video {
    id: string;
    title: string;
    channel: string;
    duration: string;
  }

  let searchQuery = "";
  let videos: Video[] = [];
  let loading = false;
  let HasSearched = false;
  let selectedVideo: Video | null = null;
  let showModal = false;
  let showSettings = false;
  let errorMessage = "";
  let selectedVideoIndex = -1;

  function reloadPage() {
    window.location.reload();
    console.log("hasSearched:", HasSearched);
  }

  function toggleSettings() {
    showSettings = !showSettings;
  }

  async function handleSubmit(event: Event) {
    event.preventDefault();
    loading = true;
    HasSearched = true;
    errorMessage = "";

    // Add the 'visible-overflow' class and remove 'hidden-overflow'
    document.documentElement.classList.add("visible-overflow");
    document.documentElement.classList.remove("hidden-overflow");
    document.body.classList.add("visible-overflow");
    document.body.classList.remove("hidden-overflow");

    try {
      const result = await Search(searchQuery);
      if (result && result.length > 0) {
        videos = result as unknown as Video[];
      } else {
        videos = [];
        errorMessage = `No videos found for "${searchQuery}". Try a different search term.`;
      }
    } catch (error) {
      console.error("Error fetching videos:", error);
      errorMessage =
        "An error occurred while searching. Please check your internet connection and try again.";
      videos = [];
    } finally {
      loading = false;
    }
  }

  function openVideo(video: Video, index: number) {
    selectedVideo = video;
    selectedVideoIndex = index;
    showModal = true;
    // Set focus to close button when modal opens
    setTimeout(() => {
      const closeButton = document.querySelector(
        ".modal-content .close-button",
      ) as HTMLElement;
      if (closeButton) closeButton.focus();
    }, 100);
  }

  function closeModal() {
    showModal = false;
    selectedVideo = null;
    selectedVideoIndex = -1;
    // Return focus to the video item that was opened
    setTimeout(() => {
      const videoButton = document.querySelector(
        `.video-item:nth-child(${selectedVideoIndex + 1}) .watch-button`,
      ) as HTMLElement;
      if (videoButton) videoButton.focus();
    }, 100);
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Escape") {
      if (showSettings) {
        showSettings = false;
      } else if (showModal) {
        closeModal();
      }
    }
  }

  onMount(() => {
    document.title = "YouTube Zen";

    // Set initial styles when the component is mounted
    if (!HasSearched) {
      document.documentElement.classList.add("hidden-overflow");
      document.body.classList.add("hidden-overflow");
    }
  });
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="app">
  <button
    class="settings-button"
    on:click={toggleSettings}
    aria-label="Open settings"
  >
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      aria-hidden="true"
    >
      <circle cx="12" cy="12" r="3"></circle>
      <path
        d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"
      ></path>
    </svg>
  </button>

  <main class={!HasSearched ? "centered" : ""}>
    <h1
      on:click={reloadPage}
      class="title"
      role="button"
      tabindex="0"
      on:keypress={(e) => e.key === "Enter" && reloadPage()}
    >
      YouTube Zen
    </h1>
    <form on:submit={handleSubmit}>
      <div class="search-container">
        <input
          type="search"
          bind:value={searchQuery}
          placeholder="Search for videos..."
          aria-label="Search for YouTube videos"
          required
        />
        <button type="submit" aria-label="Search">
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            aria-hidden="true"
          >
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
        </button>
      </div>
    </form>

    {#if loading}
      <p class="message" role="status" aria-live="polite">Searching...</p>
    {:else if errorMessage}
      <p class="message" role="alert" aria-live="assertive">{errorMessage}</p>
    {:else if videos.length > 0}
      <div class="results" role="list" aria-label="Search results">
        {#each videos as video, index}
          <div class="video-item" role="listitem">
            <div class="video-info">
              <h2>{video.title}</h2>
              <p>{video.channel} • {video.duration}</p>
            </div>
            <button
              class="watch-button"
              on:click={() => openVideo(video, index)}
              aria-label="Watch {video.title}"
            >
              Watch
            </button>
          </div>
        {/each}
      </div>
    {:else if HasSearched}
      <p class="message" role="status">No videos found for "{searchQuery}".</p>
    {/if}
  </main>

  {#if showModal && selectedVideo}
    <div
      class="modal-overlay"
      on:click={closeModal}
      role="dialog"
      aria-modal="true"
      aria-labelledby="video-title"
    >
      <div class="modal-content" on:click|stopPropagation>
        <button
          class="close-button"
          on:click={closeModal}
          aria-label="Close video player">×</button
        >
        <iframe
          id="video-title"
          title={selectedVideo.title}
          src={`http://localhost:8888/youtube-embed?v=${selectedVideo.id}`}
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
          allowfullscreen
          loading="lazy"
        ></iframe>
      </div>
    </div>
  {/if}

  {#if showSettings}
    <div
      class="modal-overlay"
      on:click={toggleSettings}
      role="dialog"
      aria-modal="true"
      aria-labelledby="settings-title"
    >
      <div class="settings-modal" on:click|stopPropagation>
        <button
          class="close-button"
          on:click={toggleSettings}
          aria-label="Close settings">×</button
        >
        <h2 id="settings-title">Settings</h2>
        <div class="settings-content">
          <p>Settings panel content will go here.</p>
        </div>
      </div>
    </div>
  {/if}
</div>
