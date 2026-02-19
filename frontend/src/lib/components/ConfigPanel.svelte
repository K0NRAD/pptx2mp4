<script lang="ts">
  import { configStore } from '../stores/config.store.svelte';

  const resolutionOptions = [
    { value: 720, label: '720p (HD)' },
    { value: 1080, label: '1080p (Full HD)' },
    { value: 1440, label: '1440p (2K)' },
    { value: 2160, label: '2160p (4K)' },
  ];
</script>

<div>
  <h3 class="h5 fw-semibold text-dark mb-4">Konvertierungs-Einstellungen</h3>

  <div class="row g-4">
    <div class="col-12 col-md-4">
      <label for="fps" class="form-label d-flex justify-content-between">
        FPS (Frames per Second)
        <span class="text-secondary fw-normal small">1-60</span>
      </label>
      <input
        id="fps"
        type="number"
        min="1"
        max="60"
        value={configStore.fps}
        oninput={(event: Event) => configStore.setFps(parseInt((event.currentTarget as HTMLInputElement).value))}
        class="form-control"
      />
    </div>

    <div class="col-12 col-md-4">
      <label for="resolution" class="form-label">
        Auflösung
      </label>
      <select
        id="resolution"
        value={configStore.resolution}
        onchange={(event: Event) => configStore.setResolution(parseInt((event.currentTarget as HTMLSelectElement).value))}
        class="form-select"
      >
        {#each resolutionOptions as option}
          <option value={option.value}>{option.label}</option>
        {/each}
      </select>
    </div>

    <div class="col-12 col-md-4">
      <label for="duration" class="form-label d-flex justify-content-between">
        Dauer pro Slide (Sekunden)
        <span class="text-secondary fw-normal small">1-60</span>
      </label>
      <input
        id="duration"
        type="number"
        min="1"
        max="60"
        value={configStore.duration}
        oninput={(event: Event) => configStore.setDuration(parseInt((event.currentTarget as HTMLInputElement).value))}
        class="form-control"
      />
    </div>

    <div class="col-12 col-md-4">
      <label for="transitionDuration" class="form-label d-flex justify-content-between">
        Überblendung (Sekunden)
        <span class="text-secondary fw-normal small">0–3</span>
      </label>
      <input
        id="transitionDuration"
        type="number"
        min="0"
        max="3"
        step="0.5"
        value={configStore.transitionDuration}
        oninput={(event: Event) => configStore.setTransitionDuration(parseFloat((event.currentTarget as HTMLInputElement).value))}
        class="form-control"
      />
    </div>
  </div>

  <button onclick={() => configStore.reset()} class="btn btn-outline-secondary mt-4">
    Auf Standardwerte zurücksetzen
  </button>
</div>
