<script lang="ts">
  import FileUpload from './lib/components/FileUpload.svelte';
  import ConfigPanel from './lib/components/ConfigPanel.svelte';
  import ConversionStatus from './lib/components/ConversionStatus.svelte';
</script>

<main class="container py-3 d-flex flex-column min-vh-100">
  <header class="text-center text-white mb-3">
    <h1 class="h2 fw-bold" style="text-shadow: 2px 2px 4px rgba(0,0,0,0.2);">
      PPTX to MP4 Converter
    </h1>
  </header>

  <div class="flex-grow-1 d-flex flex-column gap-3">
    <section class="card shadow">
      <div class="card-body p-3">
        <FileUpload />
      </div>
    </section>

    <section class="card shadow">
      <div class="card-body p-3">
        <ConfigPanel />
      </div>
    </section>

    <ConversionStatus />
  </div>

  <footer class="text-center text-white py-2 opacity-75">
    <p class="mb-0"><b>iTec</b> - learn together grow together</p>
  </footer>
</main>
