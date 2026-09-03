<script lang="ts">
  import { Card, MetricCard, Skeleton } from 'bindrunes';
  import { createLiveQuery } from '$lib/api/hooks';
  import { api } from '$lib/api/client';

  const tokens = createLiveQuery(() => api.tokens(30), ['metrics']);
</script>

<svelte:head>
  <title>Tokens — Analytics</title>
</svelte:head>

<div class="grid gap-4 sm:grid-cols-3">
  {#if $tokens.loading}
    {#each Array(3) as _}
      <Card variant="glass" class="!p-4"><Skeleton lines={2} /></Card>
    {/each}
  {:else}
    <MetricCard label="TOKENS TOTAIS (30d)" value={String($tokens.data?.totalTokens ?? 0)} />
    <MetricCard label="CHAMADAS TOTAIS" value={String($tokens.data?.totalCalls ?? 0)} />
    <MetricCard label="MODELOS ÚNICOS" value={String($tokens.data?.byModel?.length ?? 0)} variant="info" />
  {/if}
</div>

<!-- By Model -->
<Card variant="glass" class="!p-6">
  <h3 class="font-semibold mb-4">Uso por Modelo</h3>
  {#if $tokens.loading}
    <Skeleton lines={4} />
  {:else if $tokens.data?.byModel?.length > 0}
    <div class="space-y-3">
      {#each $tokens.data.byModel as m}
        {@const pct = $tokens.data.totalTokens > 0 ? (m.total_tokens / $tokens.data.totalTokens * 100) : 0}
        <div class="space-y-1">
          <div class="flex justify-between text-sm">
            <span class="font-mono">{m.model}</span>
            <span class="text-muted-foreground">{m.total_tokens.toLocaleString()} tokens ({m.call_count} chamadas)</span>
          </div>
          <div class="h-2 rounded-full bg-surface-1 overflow-hidden">
            <div class="h-full rounded-full bg-primary" style="width: {pct}%"></div>
          </div>
        </div>
      {/each}
    </div>
  {:else}
    <p class="text-sm text-muted-foreground">Sem dados de tokens disponíveis.</p>
  {/if}
</Card>

<!-- Daily Trend -->
<Card variant="glass" class="!p-6">
  <h3 class="font-semibold mb-4">Tendência Diária</h3>
  {#if $tokens.loading}
    <Skeleton lines={4} />
  {:else if $tokens.data?.trend?.length > 0}
    <div class="space-y-2">
      {#each $tokens.data.trend.slice(-14) as t}
        <div class="flex items-center gap-4 text-sm">
          <span class="text-muted-foreground w-16 shrink-0">{String(t.date).slice(5)}</span>
          <div class="flex-1">
            <div class="h-4 rounded bg-primary/70" style="width: {Math.min((t.total_tokens / Math.max(...$tokens.data.trend.map((x: any) => x.total_tokens))) * 100, 100)}%"></div>
          </div>
          <span class="text-muted-foreground w-24 text-right">{t.total_tokens.toLocaleString()}</span>
        </div>
      {/each}
    </div>
  {:else}
    <p class="text-sm text-muted-foreground">Sem dados de tendência disponíveis.</p>
  {/if}
</Card>
