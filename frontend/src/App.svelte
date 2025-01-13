<script>
  import { onMount } from 'svelte';
  import { Search } from '../wailsjs/go/main/App.js';

  let searchQuery = '';
  let videos = [];
  let loading = false;
  let hasSearched = false;
  let selectedVideo = null;
  let showModal = false;

  function reloadPage() {
    window.location.reload();
  }

  async function handleSubmit(event) {
    event.preventDefault();
    loading = true;
    hasSearched = true;

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
    if (event.key === 'Escape' && showModal) {
      closeModal();
    }
  }

  onMount(() => {
    document.title = 'YouTube Zen';
  });
</script>

<svelte:window on:keydown={handleKeydown}/>

<div class="app">
  <main class={!hasSearched ? 'centered' : ''}>
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
    {:else if hasSearched}
      <p class="message">No videos found for "{searchQuery}".</p>
    {/if}
  </main>

  {#if showModal}
    <div class="modal-overlay" on:click={closeModal}>
      <div class="modal-content" on:click|stopPropagation>
        <button class="close-button" on:click={closeModal}>×</button>
        <iframe
          title={selectedVideo.title}
          width="100%"
          height="100%"
          src={`https://www.youtube.com/embed/${selectedVideo.id}`}
          frameborder="0"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
          allowfullscreen
        ></iframe>
      </div>
    </div>
  {/if}
</div>

<style>
  @import url('https://fonts.googleapis.com/css2?family=Poppins:wght@400;600;700&display=swap');

  :global(html, body) {
    margin: 0;
    padding: 0;
    height: 100%;
    background-color: #ffffff;
  }

  .app {
    font-family: 'Poppins', sans-serif;
    min-height: 100vh;
    background-color: #ffffff;
    color: #333;
  }

  main {
    max-width: 600px;
    margin: 0 auto;
    padding: 2rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    min-height: 100vh;
  }

  /* New styles for centered initial state */
  .centered {
    margin-top: 10rem;
  }

  /* Modified title styles */
  .title {
    font-size: 2.5rem;
    font-weight: 700;
    margin-bottom: 2rem;
    color: #1a1a1a;
    cursor: pointer;
    transition: color 0.2s ease;
  }

  .title:hover {
    color: #3b82f6;
  }

  .search-container {
    position: relative;
    width: 100%;
    max-width: 500px;
    margin-bottom: 2rem;
  }

  input {
    width: 100%;
    padding: 1rem 3rem 1rem 1.5rem;
    font-size: 1rem;
    border: 2px solid #e0e0e0;
    border-radius: 9999px;
    outline: none;
    transition: border-color 0.3s ease;
    background: #ffffff;
  }

  input:focus {
    border-color: #3b82f6;
  }

  button {
    position: absolute;
    right: 1rem;
    top: 50%;
    transform: translateY(-50%);
    background: none;
    border: none;
    cursor: pointer;
    padding: 0;
  }

  button svg {
    width: 1.5rem;
    height: 1.5rem;
    color: #666;
    transition: color 0.3s ease;
  }

  button:hover svg {
    color: #3b82f6;
  }

  .results {
    width: 100%;
    max-width: 600px;
  }

  .video-item {
    background: #f8f9fa;
    border-radius: 8px;
    padding: 1rem;
    margin-bottom: 1rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    transition: transform 0.2s ease;
  }

  .video-item:hover {
    transform: translateY(-2px);
  }

  .video-info {
    flex: 1;
    margin-right: 1rem;
  }

  .video-info h2 {
    font-size: 1rem;
    font-weight: 600;
    margin: 0 0 0.5rem 0;
    color: #1a1a1a;
  }

  .video-info p {
    font-size: 0.9rem;
    color: #666;
    margin: 0;
  }

  .watch-button {
    position: static;
    color: #3b82f6;
    text-decoration: none;
    font-weight: 600;
    padding: 0.5rem 1rem;
    border-radius: 9999px;
    background: #e6f0ff;
    transition: background-color 0.2s ease;
  }

  .watch-button:hover {
    background: #d1e3ff;
  }

  .message {
    color: #666;
    margin-top: 2rem;
    text-align: center;
  }

  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.75);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .modal-content {
    position: relative;
    width: 90%;
    max-width: 800px;
    height: 0;
    padding-bottom: 56.25%; /* 16:9 Aspect Ratio */
    background: #fff;
    border-radius: 8px;
  }

  .close-button {
    position: absolute;
    top: -40px;
    right: -40px;
    width: 30px;
    height: 30px;
    background: #fff;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    color: #333;
    cursor: pointer;
    z-index: 1001;
  }

  iframe {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    border-radius: 8px;
  }
</style>