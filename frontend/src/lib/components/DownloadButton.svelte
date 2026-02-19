<script lang="ts">
  import { apiClient } from '../api/api-client';

  let { jobId }: { jobId: string } = $props();

  let isDownloading = $state(false);

  async function handleDownload() {
    if (isDownloading) return;

    isDownloading = true;

    try {
      const blob = await apiClient.downloadFile(jobId);
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = 'output.mp4';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Download fehlgeschlagen:', error);
      alert(
        'Download fehlgeschlagen: ' +
          (error instanceof Error ? error.message : 'Unbekannter Fehler'),
      );
    } finally {
      isDownloading = false;
    }
  }
</script>

<button
  onclick={handleDownload}
  disabled={isDownloading}
  class="btn btn-success w-100 d-flex align-items-center justify-content-center gap-2 fw-semibold"
>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
    width="24"
    height="24"
  >
    <path
      stroke-linecap="round"
      stroke-linejoin="round"
      stroke-width="2"
      d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
    />
  </svg>
  {isDownloading ? 'Wird heruntergeladen...' : 'MP4 herunterladen'}
</button>
