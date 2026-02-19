<script lang="ts">
  import { apiClient } from '../api/api-client';
  import { jobStore } from '../stores/job.store.svelte';
  import { configStore } from '../stores/config.store.svelte';

  let isDragging = $state(false);
  let selectedFile: File | null = $state(null);
  let isUploading = $state(false);
  let error: string | null = $state(null);

  function handleDragOver(event: DragEvent) {
    event.preventDefault();
    isDragging = true;
  }

  function handleDragLeave() {
    isDragging = false;
  }

  function handleDrop(event: DragEvent) {
    event.preventDefault();
    isDragging = false;

    const files = event.dataTransfer?.files;
    if (files && files.length > 0) {
      selectFile(files[0]);
    }
  }

  function handleFileInput(event: Event) {
    const target = event.target as HTMLInputElement;
    if (target.files && target.files.length > 0) {
      selectFile(target.files[0]);
    }
  }

  function selectFile(file: File) {
    if (!file.name.endsWith('.pptx')) {
      error = 'Bitte wählen Sie eine PPTX-Datei aus';
      return;
    }

    if (file.size > 100 * 1024 * 1024) {
      error = 'Datei zu groß (max. 100 MB)';
      return;
    }

    selectedFile = file;
    error = null;
  }

  async function handleUpload() {
    if (!selectedFile) {
      return;
    }

    isUploading = true;
    error = null;

    try {
      const response = await apiClient.convertFile(selectedFile, configStore.current);
      jobStore.setJob(response.jobId, response.status);
      jobStore.startPolling(response.jobId);
      selectedFile = null;
    } catch (err) {
      error = err instanceof Error ? err.message : 'Upload fehlgeschlagen';
    } finally {
      isUploading = false;
    }
  }

  function clearFile() {
    selectedFile = null;
    error = null;
  }
</script>

<div class="w-100">
  <div
    class="border border-2 border-dashed rounded p-4 text-center transition-all {isDragging ? 'border-primary bg-primary bg-opacity-10' : 'border-secondary-subtle'}"
    ondragover={handleDragOver}
    ondragleave={handleDragLeave}
    ondrop={handleDrop}
    role="button"
    tabindex="0"
  >
    <input
      type="file"
      accept=".pptx"
      onchange={handleFileInput}
      id="file-input"
      class="d-none"
    />

    {#if selectedFile}
      <div class="d-flex flex-column gap-2">
        <p class="fw-semibold text-dark mb-0">{selectedFile.name}</p>
        <p class="text-secondary mb-0">
          {(selectedFile.size / (1024 * 1024)).toFixed(2)} MB
        </p>
        <div class="d-flex gap-2 justify-content-center mt-2">
          <button
            onclick={handleUpload}
            disabled={isUploading}
            class="btn btn-primary"
          >
            {isUploading ? 'Wird hochgeladen...' : 'Konvertieren'}
          </button>
          <button onclick={clearFile} class="btn btn-secondary">
            Abbrechen
          </button>
        </div>
      </div>
    {:else}
      <label for="file-input" class="d-flex flex-column align-items-center gap-3 mb-0" style="cursor: pointer;">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          width="48"
          height="48"
          class="text-primary"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
          />
        </svg>
        <p class="mb-0 text-dark">
          Ziehen Sie eine PPTX-Datei hierher oder klicken Sie zum Auswählen
        </p>
        <p class="mb-0 text-secondary small">Maximale Dateigröße: 100 MB</p>
      </label>
    {/if}
  </div>

  {#if error}
    <div class="alert alert-danger mt-3 text-center" role="alert">
      {error}
    </div>
  {/if}
</div>
