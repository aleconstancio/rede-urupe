import { onMount, onDestroy } from 'svelte';
import { api } from './client';
import { sseClient } from './sse';

export function createLiveQuery<T>(fetcher: () => Promise<T>, refreshOn?: string[]) {
  let data = $state<T | null>(null);
  let loading = $state(true);
  let error = $state<string | null>(null);

  async function refresh() {
    try {
      loading = true;
      data = await fetcher();
      error = null;
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    refresh();
    if (refreshOn?.length) {
      const unsub = sseClient.subscribe((event) => {
        if (refreshOn.includes(event)) refresh();
      });
      onDestroy(unsub);
    }
  });

  return {
    get data() { return data; },
    get loading() { return loading; },
    get error() { return error; },
    refresh,
  };
}
