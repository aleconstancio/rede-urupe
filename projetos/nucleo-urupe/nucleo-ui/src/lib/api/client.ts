const API_BASE = import.meta.env.VITE_MAZE_API_URL || 'http://localhost:9393';

interface RequestOptions {
  method?: string;
  body?: unknown;
  headers?: Record<string, string>;
}

function getChannelParam(): string {
  if (typeof window === 'undefined') return '';
  const ch = localStorage.getItem('maze_channel');
  return ch ? `&channel_id=${ch}` : '';
}

async function request<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const { method = 'GET', body, headers = {} } = opts;
  const res = await fetch(`${API_BASE}${path}`, {
    method,
    headers: { 'Content-Type': 'application/json', ...headers },
    body: body ? JSON.stringify(body) : undefined,
  });
  if (!res.ok) throw new Error(`API ${res.status}: ${res.statusText}`);
  return res.json();
}

export const api = {
  status: () => request<Record<string, unknown>>('/api/status'),
  metrics: () => request<Record<string, unknown>>(`/api/metrics?_=${Date.now()}${getChannelParam()}`),
  feed: () => request<Record<string, unknown>>(`/api/feed?_=${Date.now()}${getChannelParam()}`),
  memoryToday: (date?: string) => request<Record<string, unknown>>(`/api/memory/today${date ? `?date=${date}` : ''}${getChannelParam()}`),
  persona: (channelId?: string) => request<Record<string, unknown>>(`/api/persona${channelId ? `?channel_id=${channelId}` : ''}`),
  annotations: () => request<Record<string, unknown>>(`/api/annotations${getChannelParam()}`),

  // Analytics
  sentiment: (days = 30) => request<Record<string, unknown>>(`/api/analytics/sentiment?days=${days}${getChannelParam()}`),
  tokens: (days = 30) => request<Record<string, unknown>>(`/api/analytics/tokens?days=${days}${getChannelParam()}`),
  growth: (days = 30) => request<Record<string, unknown>>(`/api/analytics/growth?days=${days}${getChannelParam()}`),
  channels: () => request<Record<string, unknown>>('/api/analytics/channels'),
  engagement: (days = 7) => request<Record<string, unknown>>(`/api/analytics/engagement?days=${days}${getChannelParam()}`),

  // Admin
  members: () => request<Record<string, unknown>>('/api/admin/members'),
  member: (id: string) => request<Record<string, unknown>>(`/api/admin/member/${id}`),
  audit: () => request<Record<string, unknown>>('/api/admin/audit'),
  modlog: () => request<Record<string, unknown>>('/api/admin/modlog'),
  welcome: () => request<Record<string, unknown>>('/api/admin/welcome'),

  // Projects
  projects: () => request<Record<string, unknown>>('/api/projects'),

  // Mode
  setMode: (passive: boolean) => request<Record<string, unknown>>('/api/mode', { method: 'POST', body: { passive } }),

  // Channels
  listChannels: () => request<Record<string, unknown>>('/api/channels'),
};
