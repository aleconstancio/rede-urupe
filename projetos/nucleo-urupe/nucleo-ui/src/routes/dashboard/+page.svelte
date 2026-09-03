<script lang="ts">
  import { Card, MetricCard, Badge, Skeleton } from 'bindrunes';
  import { createLiveQuery } from '$lib/api/hooks';
  import { api } from '$lib/api/client';

  const status = createLiveQuery(() => api.status(), ['metrics', 'feed']);
  const metrics = createLiveQuery(() => api.metrics(), ['metrics', 'feed']);
  const projects = createLiveQuery(() => api.projects());

  function uptimeLabel(minutes: number): string {
    if (minutes < 60) return `${minutes}min`;
    const h = Math.floor(minutes / 60);
    const m = minutes % 60;
    return m > 0 ? `${h}h ${m}min` : `${h}h`;
  }
</script>

<svelte:head>
  <title>Painel — Labirinto de Dédalo</title>
</svelte:head>

<div class="space-y-6">
  <div>
    <h1 class="text-2xl font-bold">Painel</h1>
    <p class="text-muted-foreground">Visão geral dos seus produtos</p>
  </div>

  <!-- System Status -->
  <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
    {#if $status.loading}
      {#each Array(4) as _}
        <Card variant="glass" class="!p-4"><Skeleton lines={2} /></Card>
      {/each}
    {:else if $status.data}
      <MetricCard
        label="STATUS"
        value={$status.data.online ? 'Online' : 'Offline'}
        variant={$status.data.online ? 'success' : 'destructive'}
      />
      <MetricCard
        label="UPTIME"
        value={uptimeLabel($status.data.uptime_minutes ?? 0)}
      />
      <MetricCard
        label="MEMENSAGENS TOTAL"
        value={String($status.data.total_messages ?? 0)}
      />
      <MetricCard
        label="MEMBROS (7d)"
        value={String($status.data.active_members_7d ?? 0)}
      />
    {/if}
  </div>

  <!-- Channel Metrics -->
  <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
    {#if $metrics.loading}
      {#each Array(4) as _}
        <Card variant="glass" class="!p-4"><Skeleton lines={2} /></Card>
      {/each}
    {:else if $metrics.data}
      <MetricCard
        label="USUÁRIOS ATIVOS (1h)"
        value={String($metrics.data.active_users ?? 0)}
        variant="success"
      />
      <MetricCard
        label="MENSAGENS (1h)"
        value={String($metrics.data.messages_1h ?? 0)}
      />
      <MetricCard
        label="BOT HOJE"
        value={String($metrics.data.bot_messages ?? 0)}
      />
      <MetricCard
        label="MEMÓRIAS HOJE"
        value={String($metrics.data.memories_today ?? 0)}
        variant="info"
      />
    {/if}
  </div>

  <!-- Active Persona -->
  {#if $status.data}
    <Card variant="glass" class="!p-6">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-lg font-semibold">Persona Ativa</h2>
          <p class="text-muted-foreground">Identity and style currently in use</p>
        </div>
        <Badge variant="primary">
          {$metrics.data?.active_persona ?? 'Bot'}
        </Badge>
      </div>
    </Card>
  {/if}

  <!-- Projects -->
  {#if $projects.data?.projects?.length > 0}
    <Card variant="glass" class="!p-6">
      <h2 class="text-lg font-semibold mb-4">Projetos</h2>
      <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {#each $projects.data.projects as p}
          <div class="flex items-center gap-3 p-3 rounded-lg bg-surface-1/50">
            <span class="text-2xl">{p.icon}</span>
            <div>
              <div class="font-medium">{p.name}</div>
              <div class="text-xs text-muted-foreground">{p.category}</div>
            </div>
          </div>
        {/each}
      </div>
    </Card>
  {/if}
</div>
