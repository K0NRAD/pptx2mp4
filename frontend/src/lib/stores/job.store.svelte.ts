import type { JobStatus } from '../api/api-client';
import { apiClient } from '../api/api-client';

interface JobState {
  jobId: string | undefined;
  status: JobStatus['status'] | undefined;
  progress: number;
  isPolling: boolean;
  error: string | null;
}

const INITIAL_STATE: JobState = {
  jobId: undefined,
  status: undefined,
  progress: 0,
  isPolling: false,
  error: null,
};

function createJobStore() {
  let state = $state<JobState>({ ...INITIAL_STATE });
  let pollInterval: number | undefined;

  function stopPolling() {
    if (pollInterval) {
      clearInterval(pollInterval);
      pollInterval = undefined;
    }
    state.isPolling = false;
  }

  return {
    get jobId() {
      return state.jobId;
    },
    get status() {
      return state.status;
    },
    get progress() {
      return state.progress;
    },
    get isPolling() {
      return state.isPolling;
    },
    get error() {
      return state.error;
    },

    reset() {
      stopPolling();
      state = { ...INITIAL_STATE };
    },

    setJob(jobId: string, status: string) {
      state.jobId = jobId;
      state.status = status as JobState['status'];
      state.progress = 0;
      state.error = null;
    },

    updateStatus(jobStatus: JobStatus) {
      state.jobId = jobStatus.jobId;
      state.status = jobStatus.status;
      state.progress = jobStatus.progress;
      if (jobStatus.error) {
        state.error = jobStatus.error;
      }
    },

    setError(error: string) {
      state.error = error;
      state.isPolling = false;
      stopPolling();
    },

    async startPolling(jobId: string) {
      stopPolling();
      state.isPolling = true;

      const poll = async () => {
        try {
          const status = await apiClient.getJobStatus(jobId);
          this.updateStatus(status);

          if (status.status === 'completed' || status.status === 'failed') {
            stopPolling();
            if (status.status === 'failed') {
              this.setError(status.error || 'Konvertierung fehlgeschlagen');
            }
          }
        } catch (error) {
          console.error('Fehler beim Abrufen des Job-Status:', error);
          this.setError(
            error instanceof Error ? error.message : 'Unbekannter Fehler',
          );
        }
      };

      await poll();
      pollInterval = window.setInterval(poll, 2000);
    },

    stopPolling,
  };
}

export const jobStore = createJobStore();
