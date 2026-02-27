// In production (embedded): BASE_URL = '/pptx2mp4/' → API unter /pptx2mp4/api/v1/...
// Im Dev-Modus mit separatem Backend: VITE_API_BASE_URL=http://localhost:8080/pptx2mp4 setzen
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || import.meta.env.BASE_URL.replace(/\/$/, '');

export interface ConversionConfig {
  fps: number;
  resolution: number;
  duration: number;
  transitionDuration: number;
}

export interface JobStatus {
  jobId: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  progress: number;
  error?: string;
}

export interface ConvertResponse {
  jobId: string;
  status: string;
}

export interface ErrorResponse {
  error: string;
  message: string;
}

export class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  async convertFile(
    file: File,
    config: ConversionConfig
  ): Promise<ConvertResponse> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('fps', config.fps.toString());
    formData.append('resolution', config.resolution.toString());
    formData.append('duration', config.duration.toString());
    formData.append('transitionDuration', config.transitionDuration.toString());

    const response = await fetch(`${this.baseUrl}/api/v1/convert`, {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      const error: ErrorResponse = await response.json();
      throw new Error(error.message || 'Konvertierung fehlgeschlagen');
    }

    return response.json();
  }

  async getJobStatus(jobId: string): Promise<JobStatus> {
    const response = await fetch(
      `${this.baseUrl}/api/v1/jobs/${jobId}/status`
    );

    if (!response.ok) {
      const error: ErrorResponse = await response.json();
      throw new Error(error.message || 'Status konnte nicht abgerufen werden');
    }

    return response.json();
  }

  getDownloadUrl(jobId: string): string {
    return `${this.baseUrl}/api/v1/jobs/${jobId}/download`;
  }

  async downloadFile(jobId: string): Promise<Blob> {
    const response = await fetch(this.getDownloadUrl(jobId));

    if (!response.ok) {
      const error: ErrorResponse = await response.json();
      throw new Error(error.message || 'Download fehlgeschlagen');
    }

    return response.blob();
  }

  async checkHealth(): Promise<{ status: string; [key: string]: any }> {
    const response = await fetch(`${this.baseUrl}/api/v1/health`);
    return response.json();
  }
}

export const apiClient = new ApiClient();
