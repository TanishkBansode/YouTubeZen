<script>
  import { onMount } from 'svelte';
  import { Search } from '../wailsjs/go/main/App.js';
  import './styles/styles.css';
  import './styles/theme-light.css';
  import './styles/theme-dark.css';

  let searchQuery = '';
  let videos = [];
  let loading = false;
  let HasSearched = false;
  let selectedVideo = null;
  let showModal = false;
  let showSettings = false;  // New state for settings

  function reloadPage() {
    window.location.reload();
    console.log('hasSearched:', HasSearched);
  }

  function toggleSettings() {  // New function for settings
    showSettings = !showSettings;
  }

  async function handleSubmit(event) {
    event.preventDefault();
    loading = true;
    HasSearched = true;

    // Add the 'visible-overflow' class and remove 'hidden-overflow'
    document.documentElement.classList.add('visible-overflow');
    document.documentElement.classList.remove('hidden-overflow');
    document.body.classList.add('visible-overflow');
    document.body.classList.remove('hidden-overflow');

    try {
      const result = await Search(searchQuery);
      videos = result;
    } catch (error) {
      console.error('Error fetching videos:', error);
    } finally {
      loading = false;
    }
  }

  function openVideo(video) {
    selectedVideo = video;
    showModal = true;
  }

  function closeModal() {
    showModal = false;
    selectedVideo = null;
  }

  function handleKeydown(event) {
    if (event.key === 'Escape') {
      if (showSettings) {
        showSettings = false;
      } else if (showModal) {
        closeModal();
      }
    }
  }

  onMount(() => {
    document.title = 'YouTube Zen';

    // Set initial styles when the component is mounted
    if (!HasSearched) {
      document.documentElement.classList.add('hidden-overflow');
      document.body.classList.add('hidden-overflow');
    }
  });
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="app">
  <button class="settings-button" on:click={toggleSettings} aria-label="Settings">
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="12" cy="12" r="3"></circle>
      <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
    </svg>
  </button>

  <main class={!HasSearched ? 'centered' : ''}>
    <h1 on:click={reloadPage} class="title">YouTube Zen</h1>
    <form on:submit={handleSubmit}>
      <div class="search-container">
        <input
          type="search"
          bind:value={searchQuery}
          placeholder="Search for videos..."
          required
        />
        <button type="submit" aria-label="Search">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
        </button>
      </div>
    </form>

    {#if loading}
      <p class="message">Searching...</p>
    {:else if videos.length > 0}
      <div class="results">
        {#each videos as video}
          <div class="video-item">
            <div class="video-info">
              <h2>{video.title}</h2>
              <p>{video.channel} • {video.duration}</p>
            </div>
            <button class="watch-button" on:click={() => openVideo(video)}>
              Watch
            </button>
          </div>
        {/each}
      </div>
    {:else if HasSearched}
      <p class="message">No videos found for "{searchQuery}".</p>
    {/if}
  </main>

  {#if showModal}
    <div class="modal-overlay" on:click={closeModal}>
      <div class="modal-content" on:click|stopPropagation>
        <button class="close-button" on:click={closeModal}>×</button>
        <iframe
          title={selectedVideo.title}
          src={`https://www.youtube.com/embed/${selectedVideo.id}`}
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
          allowfullscreen
        ></iframe>
      </div>
    </div>
  {/if}

  {#if showSettings}
    <div class="modal-overlay" on:click={toggleSettings}>
      <div class="settings-modal" on:click|stopPropagation>
        <button class="close-button" on:click={toggleSettings}>×</button>
        <h2>Settings</h2>
        <div class="settings-content">
          <p>Settings panel content will go here.</p>
        </div>
      </div>
    </div>
  {/if}
</div>
