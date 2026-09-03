<script lang="ts">
  import type { Snippet } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { AuthGuard, createAuth, Omnibar, deriveOmnibarOptions, createOmnibar } from 'bindrunes';
  import { DashboardShell } from 'bindrunes/dashboard';
  import { supabase } from '$lib/api';
  import { api } from '$lib/api/client';

  let { children }: { children: Snippet } = $props();

  const auth = createAuth({
    onLogout: () => goto('/login')
  });

  let channels = $state<string[]>([]);
  let selectedChannel = $state('');

  onMount(async () => {
    selectedChannel = localStorage.getItem('maze_channel') ?? '';
    try {
      const res = await api.listChannels();
      channels = (res as { channels?: string[] }).channels ?? [];
    } catch {
      channels = [];
    }
  });

  function selectChannel(channelId: string) {
    selectedChannel = channelId;
    if (channelId) {
      localStorage.setItem('maze_channel', channelId);
    } else {
      localStorage.removeItem('maze_channel');
    }
    window.location.reload();
  }

  const navigationGroups = [
    {
      label: 'Hub',
      items: [
        { title: 'Painel', to: '/dashboard', description: 'Visão geral dos produtos', icon: '🏠' },
        { title: 'Vico', to: '/dashboard/vico', description: 'Automação jurídica', icon: '⚖️' },
        { title: 'Orb', to: '/dashboard/orb', description: 'Prática mística', icon: '🔮' },
        { title: 'Labirinto', to: '/dashboard/maze', description: 'Bot Discord', icon: '🤖' },
        { title: 'Moderação', to: '/dashboard/moderation', description: 'Warnings e audit', icon: '🛡️' },
        { title: 'Analytics', to: '/dashboard/analytics', description: 'Métricas e tendências', icon: '📊' },
      ],
    },
    {
      label: 'Sistema',
      items: [
        { title: 'Configurações', to: '/dashboard/settings', description: 'Conta e preferências', icon: '⚙️' },
      ],
    },
  ];

  const omnibarOptions = deriveOmnibarOptions(navigationGroups);
  const omnibar = createOmnibar({ options: omnibarOptions });

  async function handleLogout() {
    await supabase.auth.signOut();
    auth.logout();
  }
</script>

<AuthGuard auth={auth} fallback="/login" navigate={(url) => goto(url)}>
  <DashboardShell
    appName="LABIRINTO"
    appSubtitle="de Dédalo"
    brandIcon="🌀"
    navigation={navigationGroups}
    pathname={$page.url.pathname}
    headerPrefix="HUB CENTRAL"
    sidebarCollapsible="icon"
    statusChip={{ variant: 'success', label: 'Online', dot: true, animate: true }}
  >
    {#snippet headerActions()}
      {#if channels.length > 1}
        <select
          value={selectedChannel}
          onchange={(e) => selectChannel((e.target as HTMLSelectElement).value)}
          class="rounded-md px-2 py-1.5 text-sm text-muted-foreground bg-transparent border border-border hover:bg-accent hover:text-accent-foreground transition-colors"
        >
          <option value="">Todos os canais</option>
          {#each channels as ch}
            <option value={ch}>{ch}</option>
          {/each}
        </select>
      {/if}
      <button
        onclick={() => omnibar.open()}
        class="rounded-md px-3 py-1.5 text-sm text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors"
      >
        ⌘K Buscar
      </button>
      <button
        onclick={handleLogout}
        class="rounded-md px-3 py-1.5 text-sm text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors"
      >
        Sair
      </button>
    {/snippet}

    {@render children()}
  </DashboardShell>

  <Omnibar state={omnibar} />
</AuthGuard>
