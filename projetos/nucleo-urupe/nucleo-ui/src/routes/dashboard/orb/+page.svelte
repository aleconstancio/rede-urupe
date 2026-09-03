<script lang="ts">
  import { Card, MetricCard, Badge, Skeleton, Alert } from 'bindrunes';
  import { onMount } from 'svelte';

  const orbUrl = import.meta.env.VITE_ORB_URL || 'http://localhost:8080';

  let data = $state<Record<string, unknown> | null>(null);
  let loading = $state(true);
  let reachable = $state(false);

  onMount(async () => {
    try {
      const res = await fetch(`${orbUrl}/api/v1/status`);
      if (res.ok) {
        data = await res.json();
        reachable = true;
      }
    } catch {
      reachable = false;
    } finally {
      loading = false;
    }
  });
</script>

<svelte:head>
  <title>Orb — Labirinto de Dédalo</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Orb</h1>
      <p class="text-muted-foreground">Prática mística local-first</p>
    </div>
    <Badge variant={reachable ? 'success' : 'info'}>
      {reachable ? 'Conectado' : 'Local-first'}
    </Badge>
  </div>

  {#if !reachable && !loading}
    <Alert variant="info" title="Modo local-first" description="Orb roda localmente. Dados sempre disponíveis no dispositivo." />
  {/if}

  <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
    {#if loading}
      {#each Array(3) as _}
        <Card variant="glass" class="!p-4"><Skeleton lines={2} /></Card>
      {/each}
    {:else}
      <MetricCard label="STATUS" value={reachable ? 'Online' : 'Local-first'} variant={reachable ? 'success' : 'info'} />
      <MetricCard label="TIPO" value="Local-first" />
      <MetricCard label="SYNC" value={reachable ? 'Ativo' : 'Offline'} />
    {/if}
  </div>
</div>
