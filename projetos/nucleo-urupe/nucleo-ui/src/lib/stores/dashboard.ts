import { writable } from 'svelte/store';

export interface ProductStatus {
  id: string;
  name: string;
  icon: string;
  status: 'online' | 'offline' | 'loading';
  url: string;
  stats: Record<string, string | number>;
}

function createDashboardStore() {
  const { subscribe, update } = writable<ProductStatus[]>([
    {
      id: 'vico',
      name: 'Vico',
      icon: 'scale',
      status: 'loading',
      url: import.meta.env.VITE_VICO_URL || 'http://localhost:5173',
      stats: {}
    },
    {
      id: 'orb',
      name: 'Orb',
      icon: 'circle-dot',
      status: 'loading',
      url: import.meta.env.VITE_ORB_URL || 'http://localhost:8080',
      stats: {}
    },
    {
      id: 'maze',
      name: 'Labirinto',
      icon: 'bot',
      status: 'loading',
      url: import.meta.env.VITE_MAZE_API_URL || 'http://localhost:9393',
      stats: {}
    }
  ]);

  return {
    subscribe,
    updateStatus: (id: string, status: ProductStatus['status'], stats?: Record<string, string | number>) =>
      update((products) =>
        products.map((p) =>
          p.id === id
            ? { ...p, status, stats: stats ?? p.stats }
            : p
        )
      ),
    fetchAll: async (token: string | null) => {
      const vicoUrl = import.meta.env.VITE_VICO_URL || 'http://localhost:5173';
      const mazeUrl = import.meta.env.VITE_MAZE_API_URL || 'http://localhost:9393';

      // Fetch Vico status
      try {
        const headers: Record<string, string> = {};
        if (token) headers['Authorization'] = `Bearer ${token}`;
        const resp = await fetch(`${vicoUrl}/api/v1/status`, { headers });
        if (resp.ok) {
          const data = await resp.json();
          update((products) =>
            products.map((p) =>
              p.id === 'vico'
                ? { ...p, status: 'online', stats: {
                    'Casos': data.case_count ?? 0,
                    'Plano': data.plan_tier ?? 'starter',
                  }}
                : p
            )
          );
        } else {
          update((products) => products.map((p) => p.id === 'vico' ? { ...p, status: 'offline' } : p));
        }
      } catch {
        update((products) => products.map((p) => p.id === 'vico' ? { ...p, status: 'offline' } : p));
      }

      // Fetch Maze/Labirinto status
      try {
        const resp = await fetch(`${mazeUrl}/api/status`);
        if (resp.ok) {
          const data = await resp.json();
          update((products) =>
            products.map((p) =>
              p.id === 'maze'
                ? { ...p, status: data.online ? 'online' : 'offline', stats: {
                    'Uptime': `${data.uptime_minutes ?? 0}min`,
                    'Mensagens': data.total_messages ?? 0,
                    'Membros': data.active_members_7d ?? 0,
                  }}
                : p
            )
          );
        } else {
          update((products) => products.map((p) => p.id === 'maze' ? { ...p, status: 'offline' } : p));
        }
      } catch {
        update((products) => products.map((p) => p.id === 'maze' ? { ...p, status: 'offline' } : p));
      }

      // Orb is Flutter local-first — mark as online (always available)
      update((products) =>
        products.map((p) =>
          p.id === 'orb'
            ? { ...p, status: 'online', stats: { 'Tipo': 'Local-first' } }
            : p
        )
      );
    }
  };
}

export const dashboard = createDashboardStore();
