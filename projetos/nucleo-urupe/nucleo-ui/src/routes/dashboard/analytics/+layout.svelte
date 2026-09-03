<script lang="ts">
  import type { Snippet } from 'svelte';
  import { page } from '$app/stores';
  import { Tabs, TabsList, TabsTrigger } from 'bindrunes';

  let { children }: { children: Snippet } = $props();

  const tabs = [
    { value: 'sentiment', label: 'Sentimento', href: '/dashboard/analytics' },
    { value: 'tokens', label: 'Tokens', href: '/dashboard/analytics/tokens' },
    { value: 'growth', label: 'Crescimento', href: '/dashboard/analytics/growth' },
    { value: 'channels', label: 'Canais', href: '/dashboard/analytics/channels' },
  ];

  let activeTab = $derived(tabs.find(t => $page.url.pathname === t.href)?.value ?? 'sentiment');
</script>

<div class="space-y-6">
  <div>
    <h1 class="text-2xl font-bold">Analytics</h1>
    <p class="text-muted-foreground">Métricas, tendências e saúde do servidor</p>
  </div>

  <Tabs value={activeTab}>
    <TabsList>
      {#each tabs as t}
        <TabsTrigger value={t.value}>
          <a href={t.href}>{t.label}</a>
        </TabsTrigger>
      {/each}
    </TabsList>
  </Tabs>

  {@render children()}
</div>
