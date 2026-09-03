<script lang="ts">
  import { Card, MetricCard, Skeleton } from 'bindrunes';
  import { createLiveQuery } from '$lib/api/hooks';
  import { api } from '$lib/api/client';

  const sentiment = createLiveQuery(() => api.sentiment(30), ['feed', 'metrics']);
</script>

<svelte:head>
  <title>Sentimento — Analytics</title>
</svelte:head>

<div class="grid gap-4 sm:grid-cols-3">
  {#if $sentiment.loading}
    {#each Array(3) as _}
      <Card variant="glass" class="!p-4"><Skeleton lines={2} /></Card>
    {/each}
  {:else}
    <MetricCard label="INTENSIDADE MÉDIA" value={$sentiment.data?.avgIntensity ?? '—'} />
    <MetricCard label="CONFLITO MÉDIO" value={$sentiment.data?.avgConflict ?? '—'} variant="warning" />
    <MetricCard label="MENSAGENS" value={String($sentiment.data?.totalMessages ?? 0)} />
  {/if}
</div>

<Card variant="glass" class="!p-6">
  <h3 class="font-semibold mb-4">Tendência de Sentimento (30 dias)</h3>
  {#if $sentiment.loading}
    <Skeleton lines={6} />
  {:else if $sentiment.data?.trends?.length > 0}
    <div class="space-y-2">
      {#each $sentiment.data.trends.slice(-14) as t}
        <div class="flex items-center gap-4 text-sm">
          <span class="text-muted-foreground w-16 shrink-0">{String(t.date).slice(5)}</span>
          <div class="flex-1 flex gap-1">
            <div class="h-4 rounded" style="width: {Math.min((t.avg_intensity ?? 0) * 100, 100)}%; background: oklch(0.6 0.2 250)"></div>
            <div class="h-4 rounded" style="width: {Math.min((t.avg_euphoria ?? 0) * 100, 100)}%; background: oklch(0.65 0.2 145)"></div>
            <div class="h-4 rounded" style="width: {Math.min((t.avg_conflict ?? 0) * 100, 100)}%; background: oklch(0.6 0.2 25)"></div>
          </div>
        </div>
      {/each}
    </div>
    <div class="flex gap-4 mt-4 text-xs text-muted-foreground">
      <span class="flex items-center gap-1"><span class="w-3 h-3 rounded" style="background: oklch(0.6 0.2 250)"></span> Intensidade</span>
      <span class="flex items-center gap-1"><span class="w-3 h-3 rounded" style="background: oklch(0.65 0.2 145)"></span> Euforia</span>
      <span class="flex items-center gap-1"><span class="w-3 h-3 rounded" style="background: oklch(0.6 0.2 25)"></span> Conflito</span>
    </div>
  {:else}
    <p class="text-sm text-muted-foreground">Sem dados de sentimento disponíveis.</p>
  {/if}
</Card>
