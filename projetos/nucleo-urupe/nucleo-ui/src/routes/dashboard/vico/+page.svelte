<script lang="ts">
  import { Card, MetricCard, Badge, Skeleton, Alert } from 'bindrunes';
  import { onMount } from 'svelte';

  const vicoUrl = import.meta.env.VITE_VICO_URL || 'http://localhost:5173';

  let data = $state<Record<string, unknown> | null>(null);
  let loading = $state(true);
  let reachable = $state(false);

  onMount(async () => {
    try {
      const res = await fetch(`${vicoUrl}/api/v1/status`);
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
  <title>Vico — Labirinto de Dédalo</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Vico</h1>
      <p class="text-muted-foreground">Automação jurídica com copilot IA</p>
    </div>
    <Badge variant={reachable ? 'success' : 'warning'}>
      {reachable ? 'Conectado' : 'Indisponível'}
    </Badge>
  </div>

  {#if !reachable && !loading}
    <Alert variant="warning" title="Serviço indisponível" description="Vico não está acessível. Verifique se o serviço está rodando." />
  {/if}

  <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
    {#if loading}
      {#each Array(3) as _}
        <Card variant="glass" class="!p-4"><Skeleton lines={2} /></Card>
      {/each}
    {:else}
      <MetricCard label="STATUS" value={reachable ? 'Online' : 'Offline'} variant={reachable ? 'success' : 'warning'} />
      <MetricCard label="CASOS" value={String(data?.case_count ?? 0)} />
      <MetricCard label="PLANO" value={String(data?.plan_tier ?? 'starter')} />
    {/if}
  </div>
</div>
