<script lang="ts">
  import { Card, MetricCard, Skeleton } from 'bindrunes';
  import { createLiveQuery } from '$lib/api/hooks';
  import { api } from '$lib/api/client';

  const channels = createLiveQuery(() => api.channels(), ['metrics', 'feed']);
  const engagement = createLiveQuery(() => api.engagement(7), ['metrics']);
</script>

<svelte:head>
  <title>Canais — Analytics</title>
</svelte:head>

<div class="grid gap-4 sm:grid-cols-3">
  {#if $channels.loading}
    {#each Array(3) as _}
      <Card variant="glass" class="!p-4"><Skeleton lines={2} /></Card>
    {/each}
  {:else}
    <MetricCard label="TOTAL MENSAGENS" value={String($channels.data?.totalMessages ?? 0)} />
    <MetricCard label="AUTORES ÚNICOS" value={String($channels.data?.uniqueAuthors ?? 0)} />
    <MetricCard label="ENGAGEMENT SCORE" value={String($engagement.data?.score?.toFixed(1) ?? '—')} variant="info" />
  {/if}
</div>

<Card variant="glass" class="!p-6">
  <h3 class="font-semibold mb-4">Saúde do Canal</h3>
  {#if $channels.loading}
    <Skeleton lines={4} />
  {:else if $channels.data}
    <div class="grid gap-4 sm:grid-cols-2">
      <div class="space-y-2">
        <div class="text-sm text-muted-foreground">Canal</div>
        <div class="font-medium">{$channels.data.channelName ?? '—'}</div>
      </div>
      <div class="space-y-2">
        <div class="text-sm text-muted-foreground">Mensagens Bot</div>
        <div class="font-medium">{$channels.data.botMessages ?? 0}</div>
      </div>
      <div class="space-y-2">
        <div class="text-sm text-muted-foreground">Reações Médias</div>
        <div class="font-medium">{$channels.data.avgReactions?.toFixed(1) ?? '0'}</div>
      </div>
      <div class="space-y-2">
        <div class="text-sm text-muted-foreground">Velocidade</div>
        <div class="font-medium">{$channels.data.messageVelocity?.toFixed(1) ?? '0'} msg/h</div>
      </div>
    </div>
  {:else}
    <p class="text-sm text-muted-foreground">Sem dados de canal disponíveis.</p>
  {/if}
</Card>
