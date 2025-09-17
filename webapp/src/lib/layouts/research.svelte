<script>
  import { page } from '$app/stores';
  import Icon from '@iconify/svelte';
  import { fade, fly } from 'svelte/transition';

  export let title = '';
  export let author = '';
  export let date = '';
  export let tags = [];
</script>

<article class="research-article" in:fade={{ duration: 600 }}>
  <div class="article-header" in:fly={{ y: 20, duration: 600, delay: 200 }}>
    <a href="/press" class="back-link">
      <Icon icon="mdi:arrow-left" class="back-icon" />
      Back to Press
    </a>

    <h1 class="article-title">{title}</h1>

    <div class="article-meta">
      {#if author}
        <div class="meta-item">
          <Icon icon="mdi:account" class="meta-icon" />
          <span>{author}</span>
        </div>
      {/if}
      {#if date}
        <div class="meta-item">
          <Icon icon="mdi:calendar" class="meta-icon" />
          <span>{new Date(date).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}</span>
        </div>
      {/if}
    </div>

    {#if tags && tags.length > 0}
      <div class="article-tags">
        {#each tags as tag}
          <span class="tag">{tag}</span>
        {/each}
      </div>
    {/if}
  </div>

  <div class="article-content" in:fly={{ y: 30, duration: 600, delay: 400 }}>
    <slot />
  </div>
</article>

<style>
  .research-article {
    max-width: 900px;
    margin: 0 auto;
    padding: 2rem;
  }

  .article-header {
    margin-bottom: 3rem;
  }

  .back-link {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    color: var(--founder);
    text-decoration: none;
    font-weight: 500;
    margin-bottom: 2rem;
    transition: all 0.2s ease;
  }

  .back-link:hover {
    color: var(--pri);
    transform: translateX(-4px);
  }

  :global(.back-icon) {
    font-size: 1.25rem;
  }

  .article-title {
    font-family: "Space Grotesk", sans-serif;
    font-size: 2.5rem;
    font-weight: 700;
    color: var(--pri);
    margin-bottom: 1.5rem;
    line-height: 1.2;
    letter-spacing: -0.02em;
  }

  .article-meta {
    display: flex;
    gap: 2rem;
    margin-bottom: 1.5rem;
    flex-wrap: wrap;
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: var(--text-secondary);
    font-size: 0.925rem;
  }

  :global(.meta-icon) {
    font-size: 1.125rem;
    color: var(--founder);
  }

  .article-tags {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
  }

  .tag {
    background: var(--bg-secondary);
    color: var(--text-secondary);
    padding: 0.375rem 0.75rem;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
  }

  .article-content {
    background: white;
    border-radius: 1.5rem;
    border: 1px solid var(--border);
    padding: 3rem;
    box-shadow: var(--shadow-md);
  }

  /* Article content styling */
  :global(.article-content h2) {
    font-family: "Space Grotesk", sans-serif;
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--pri);
    margin: 2.5rem 0 1rem;
    letter-spacing: -0.01em;
  }

  :global(.article-content h3) {
    font-family: "Space Grotesk", sans-serif;
    font-size: 1.375rem;
    font-weight: 600;
    color: var(--pri);
    margin: 2rem 0 0.875rem;
  }

  :global(.article-content p) {
    line-height: 1.7;
    margin-bottom: 1.25rem;
    color: var(--text-primary);
  }

  :global(.article-content a) {
    color: var(--founder);
    text-decoration: none;
    font-weight: 500;
    border-bottom: 1px solid transparent;
    transition: all 0.2s ease;
  }

  :global(.article-content a:hover) {
    color: var(--pri);
    border-bottom-color: var(--pri);
  }

  :global(.article-content ul),
  :global(.article-content ol) {
    margin: 1.25rem 0;
    padding-left: 2rem;
  }

  :global(.article-content li) {
    margin-bottom: 0.5rem;
    line-height: 1.7;
    color: var(--text-primary);
  }

  :global(.article-content blockquote) {
    border-left: 4px solid var(--founder);
    padding-left: 1.5rem;
    margin: 1.5rem 0;
    font-style: italic;
    color: var(--text-secondary);
  }

  :global(.article-content pre) {
    background: var(--bg-secondary);
    border-radius: 0.75rem;
    padding: 1.5rem;
    overflow-x: auto;
    margin: 1.5rem 0;
  }

  :global(.article-content code) {
    background: var(--bg-secondary);
    padding: 0.125rem 0.375rem;
    border-radius: 0.25rem;
    font-family: "JetBrains Mono", monospace;
    font-size: 0.875rem;
  }

  :global(.article-content pre code) {
    background: none;
    padding: 0;
  }

  :global(.article-content table) {
    width: 100%;
    border-collapse: collapse;
    margin: 1.5rem 0;
  }

  :global(.article-content th) {
    background: var(--bg-secondary);
    padding: 0.75rem;
    text-align: left;
    font-weight: 600;
    border-bottom: 2px solid var(--border);
  }

  :global(.article-content td) {
    padding: 0.75rem;
    border-bottom: 1px solid var(--border);
  }

  :global(.article-content hr) {
    border: none;
    border-top: 1px solid var(--border);
    margin: 2rem 0;
  }

  @media (max-width: 768px) {
    .research-article {
      padding: 1rem;
    }

    .article-title {
      font-size: 2rem;
    }

    .article-content {
      padding: 2rem 1.5rem;
    }

    .article-meta {
      gap: 1rem;
    }
  }

  @media (max-width: 480px) {
    .article-title {
      font-size: 1.75rem;
    }

    .article-content {
      padding: 1.5rem 1rem;
    }

    :global(.article-content h2) {
      font-size: 1.5rem;
    }
  }
</style>