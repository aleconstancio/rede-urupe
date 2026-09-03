<script lang="ts">
  import { Card, MetricCard, Skeleton } from 'bindrunes';
  import { createLiveQuery } from '$lib/api/hooks';
  import { api } from '$lib/api/client';

  const growth = createLiveQuery(() => api.growth(30), ['feed']);
</script>

<svelte:head>
  <title>Crescimento — Analytics</title>
</svelte:head>

<div class="grid gap-4 sm:grid-cols-3">
  {#if $growth.loading}
    {#each Array(3) as _}
      <Card variant="glass" class="!p-4"><Skeleton lines={2} /></Card>
    {/each}
  {:else}
    <MetricCard label="CRESCIMENTO LÍQUIDO" value={String($growth.data?.netGrowth ?? 0)} variant={$growth.data?.netGrowth >= 0 ? 'success' : 'destructive'} />
    <MetricCard label="TOTAL MEMBROS" value={String($growth.data?.totalMembers ?? 0)} />
    <MetricCard label="Período" value="30 dias" />
  {/if}
</div>

<Card variant="glass" class="!p-6">
  <h3 class="font-semibold mb-4">Crescimento Diário (30 dias)</h3>
  {#if $growth.loading}
    <Skeleton lines={6} />
  {:else if $growth.data?.daily?.length > 0}
    <div class="space-y-2">
      {#each $growth.data.daily.slice(-14) as d}
        <div class="flex items-center gap-4 text-sm">
          <span class="text-muted-foreground w-16 shrink-0">{String(d.date).slice(5)}</span>
          <div class="flex gap-1 flex-1">
            <div class="h-4 rounded bg-success/70" style="width: {Math.min((d.joins / Math.max(...$growth.data.daily.map((x: any) => Math.max(x.joins, x.leaves)))) * 100, 100)}%"></div>
            <div class="h-4 rounded bg-destructive/70" style="width: {Math.min((d.leaves / Math.max(...$growth.data.daily.map((x: any) => Math.max(x.joins, x.leaves)))) * 100, 100)}%"></div>
          </div>
          <span class="text-muted-foreground w-24 text-right">+{d.joins} / -{d.leaves}</span>
        </div>
      {/each}
    </div>
    <div class="flex gap-4 mt-4 text-xs text-muted-foreground">
      <span class="flex items-center gap-1"><span class="w-3 h-3 rounded bg-success/70"></span> Entradas</span>
      <span class="flex items-center gap-1"><span class="w-3 h-3 rounded bg-destructive/70"></span> Saídas</span>
    </div>
  {:else}
    <p class="text-sm text-muted-foreground">Sem dados de crescimento disponíveis.</p>
  {/if}
</Card>
