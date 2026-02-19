<script lang="ts">
  import { jobStore } from '../stores/job.store.svelte';
  import DownloadButton from './DownloadButton.svelte';

  const statusTextMap: Record<string, string> = {
    pending: 'Warten auf Start...',
    processing: 'Wird konvertiert...',
    completed: 'Konvertierung abgeschlossen!',
    failed: 'Konvertierung fehlgeschlagen',
  };

  const statusVariantMap: Record<string, string> = {
    pending: 'warning',
    processing: 'primary',
    completed: 'success',
    failed: 'danger',
  };

  let statusText = $derived(statusTextMap[jobStore.status || ''] || 'Unbekannt');
  let statusVariant = $derived(statusVariantMap[jobStore.status || ''] || 'secondary');
</script>

{#if jobStore.jobId}
  <div class="card shadow mt-4">
    <div class="card-body p-4">
      <div class="d-flex justify-content-between align-items-center mb-3">
        <h3 class="h5 fw-semibold text-dark mb-0">{statusText}</h3>
        <span class="badge bg-{statusVariant} text-uppercase">
          {jobStore.status}
        </span>
      </div>

      {#if jobStore.status === 'processing' || jobStore.status === 'pending'}
        <div class="mb-3">
          <div class="progress mb-2" style="height: 1rem;">
            <div
              class="progress-bar bg-{statusVariant}"
              role="progressbar"
              style="width: {jobStore.progress}%;"
              aria-valuenow={jobStore.progress}
              aria-valuemin={0}
              aria-valuemax={100}
            ></div>
          </div>
          <p class="text-center text-secondary small mb-0">{jobStore.progress}%</p>
        </div>
      {/if}

      {#if jobStore.status === 'completed'}
        <div class="alert alert-success d-flex align-items-center gap-2 mb-3" role="alert">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            width="32"
            height="32"
            class="flex-shrink-0"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <p class="mb-0 fw-medium">Ihre Datei ist bereit zum Download!</p>
        </div>
        <DownloadButton jobId={jobStore.jobId} />
      {/if}

      {#if jobStore.status === 'failed' && jobStore.error}
        <div class="alert alert-danger d-flex align-items-center gap-2 mb-3" role="alert">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            width="32"
            height="32"
            class="flex-shrink-0"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <p class="mb-0 fw-medium">{jobStore.error}</p>
        </div>
      {/if}

      <div class="border-top pt-3 mt-3">
        <p class="text-secondary small mb-0 font-monospace">Job-ID: {jobStore.jobId}</p>
      </div>

      {#if jobStore.status === 'completed' || jobStore.status === 'failed'}
        <button
          onclick={() => jobStore.reset()}
          class="btn btn-primary w-100 mt-3"
        >
          Neue Konvertierung starten
        </button>
      {/if}
    </div>
  </div>
{/if}
