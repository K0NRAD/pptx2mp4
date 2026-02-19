import type { ConversionConfig } from '../api/api-client';

const DEFAULT_CONFIG: ConversionConfig = {
  fps: 24,
  resolution: 1080,
  duration: 5,
  transitionDuration: 1.0,
};

const STORAGE_KEY = 'pptx2mp4_config';

function loadConfigFromStorage(): ConversionConfig {
  if (typeof window === 'undefined') {
    return DEFAULT_CONFIG;
  }

  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored) {
      return JSON.parse(stored);
    }
  } catch (error) {
    console.error('Fehler beim Laden der Konfiguration:', error);
  }

  return DEFAULT_CONFIG;
}

function saveConfigToStorage(config: ConversionConfig) {
  if (typeof window === 'undefined') {
    return;
  }

  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(config));
  } catch (error) {
    console.error('Fehler beim Speichern der Konfiguration:', error);
  }
}

function createConfigStore() {
  let config = $state<ConversionConfig>(loadConfigFromStorage());

  return {
    get fps() {
      return config.fps;
    },
    get resolution() {
      return config.resolution;
    },
    get duration() {
      return config.duration;
    },
    get transitionDuration() {
      return config.transitionDuration;
    },
    get current(): ConversionConfig {
      return config;
    },

    setFps(fps: number) {
      config = { ...config, fps };
      saveConfigToStorage(config);
    },

    setResolution(resolution: number) {
      config = { ...config, resolution };
      saveConfigToStorage(config);
    },

    setDuration(duration: number) {
      config = { ...config, duration };
      saveConfigToStorage(config);
    },

    setTransitionDuration(transitionDuration: number) {
      config = { ...config, transitionDuration };
      saveConfigToStorage(config);
    },

    reset() {
      config = { ...DEFAULT_CONFIG };
      saveConfigToStorage(config);
    },
  };
}

export const configStore = createConfigStore();
