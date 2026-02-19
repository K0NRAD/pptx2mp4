<script lang="ts">
  import FileUpload from './lib/components/FileUpload.svelte';
  import ConfigPanel from './lib/components/ConfigPanel.svelte';
  import ConversionStatus from './lib/components/ConversionStatus.svelte';
  import { jobStore } from './lib/stores/job.store.svelte';
</script>

<main class="container py-4 d-flex flex-column min-vh-100">
  <header class="text-center text-white mb-5">
    <h1 class="display-4 fw-bold" style="text-shadow: 2px 2px 4px rgba(0,0,0,0.2);">
      PPTX to MP4 Converter
    </h1>
    <p class="lead opacity-75">
      Konvertieren Sie Ihre PowerPoint-Präsentationen in hochwertige MP4-Videos
    </p>
  </header>

  <div class="flex-grow-1 d-flex flex-column gap-4">
    {#if !jobStore.jobId || jobStore.status === 'completed' || jobStore.status === 'failed'}
      <section class="card shadow">
        <div class="card-body p-4">
          <FileUpload />
        </div>
      </section>

      <section class="card shadow">
        <div class="card-body p-4">
          <ConfigPanel />
        </div>
      </section>
    {/if}

    <ConversionStatus />
  </div>

  <footer class="text-center text-white py-4 opacity-75">
    <p class="mb-0">
      Entwickelt mit <a href="https://go.dev" target="_blank" rel="noopener" class="text-white text-decoration-underline">Go</a>
      und
      <a href="https://svelte.dev" target="_blank" rel="noopener" class="text-white text-decoration-underline">Svelte</a>
    </p>
  </footer>
</main>
