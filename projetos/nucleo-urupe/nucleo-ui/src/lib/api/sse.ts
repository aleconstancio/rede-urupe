type EventHandler = (event: string) => void;

const API_BASE = import.meta.env.VITE_MAZE_API_URL || 'http://localhost:9393';

export function createSSEClient() {
  let eventSource: EventSource | null = null;
  let handlers: EventHandler[] = [];
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  let reconnectDelay = 1000;
  const MAX_DELAY = 30000;

  function connect() {
    if (eventSource) return;
    eventSource = new EventSource(`${API_BASE}/events`);

    eventSource.onmessage = () => {}; // keepalive, ignore
    eventSource.onerror = () => {
      eventSource?.close();
      eventSource = null;
      reconnectTimer = setTimeout(() => {
        reconnectDelay = Math.min(reconnectDelay * 2, MAX_DELAY);
        connect();
      }, reconnectDelay);
    };

    // Listen for all known event types
    ['feed', 'metrics', 'persona', 'projects'].forEach((type) => {
      eventSource!.addEventListener(type, () => {
        reconnectDelay = 1000; // reset backoff on successful event
        handlers.forEach((h) => { h(type); });
      });
    });
  }

  function disconnect() {
    if (reconnectTimer) clearTimeout(reconnectTimer);
    eventSource?.close();
    eventSource = null;
  }

  function subscribe(handler: EventHandler) {
    handlers.push(handler);
    if (!eventSource) connect();
    return () => {
      handlers = handlers.filter((h) => h !== handler);
      if (handlers.length === 0) disconnect();
    };
  }

  return { subscribe, connect, disconnect };
}

export const sseClient = createSSEClient();
