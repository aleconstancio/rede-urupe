import { writable } from 'svelte/store';

export const status = writable(null);
export const metrics = writable(null);
export const feed = writable(null);
export const memory = writable({ capsules: [] });
export const persona = writable(null);
export const context = writable(null);
export const annotations = writable({ items: [] });

// Theme/colorScheme removed — handled by mode-watcher

export function initSSE() {
    const eventSource = new EventSource('/events');

    const fetchAndUpdate = async (endpoint, store) => {
        try {
            const res = await fetch(`/api/${endpoint}`);
            if (res.ok) {
                const data = await res.json();
                store.set(data);
            }
        } catch (err) {
            console.error(`Failed to fetch ${endpoint}:`, err);
        }
    };

    eventSource.addEventListener('status', () => fetchAndUpdate('status', status));
    eventSource.addEventListener('metrics', () => fetchAndUpdate('metrics', metrics));
    eventSource.addEventListener('feed', () => fetchAndUpdate('feed', feed));
    eventSource.addEventListener('memory', () => fetchAndUpdate('memory/today', memory));
    eventSource.addEventListener('persona', () => fetchAndUpdate('persona', persona));
    eventSource.addEventListener('context', () => fetchAndUpdate('context', context));
    eventSource.addEventListener('annotations', () => fetchAndUpdate('annotations', annotations));

    // Initial fetch
    fetchAndUpdate('status', status);
    fetchAndUpdate('metrics', metrics);
    fetchAndUpdate('feed', feed);
    fetchAndUpdate('memory/today', memory);
    fetchAndUpdate('persona', persona);
    fetchAndUpdate('context', context);
    fetchAndUpdate('annotations', annotations);

    return () => eventSource.close();
}
