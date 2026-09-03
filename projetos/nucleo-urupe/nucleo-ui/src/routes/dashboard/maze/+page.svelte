<script lang="ts">
  import { Card, MetricCard, Badge, Skeleton, Switch } from 'bindrunes';
  import { createLiveQuery } from '$lib/api/hooks';
  import { api } from '$lib/api/client';

  const status = createLiveQuery(() => api.status(), ['metrics', 'feed']);
  const metrics = createLiveQuery(() => api.metrics(), ['metrics', 'feed']);
  const persona = createLiveQuery(() => api.persona(), ['persona']);
  const feed = createLiveQuery(() => api.feed(), ['feed']);

  let passiveMode = $state(false);

  async function togglePassive() {
    passiveMode = !passiveMode;
    await api.setMode(passiveMode);
  }

  function uptimeLabel(minutes: number): string {
    if (minutes < 60) return `${minutes}min`;
    const h = Math.floor(minutes / 60);
    const m = minutes % 60;
    return m > 0 ? `${h}h ${m}min` : `${h}h`;
  }
</script>

<svelte:head>
  <title>AIrelius — Labirinto de Dédalo</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold">Labirinto</h1>
      <p class="text-muted-foreground">Bot cognitivo Discord e hub central</p>
    </div>
    <Badge variant={$status.data?.online ? 'success' : 'destructive'}>
      {$status.data?.online ? 'Online' : 'Offline'}
    </Badge>
  </div>

  <!-- Core Metrics -->
  <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
    {#if $status.loading || $metrics.loading}
      {#each Array(4) as _}
        <Card variant="glass" class="!p-4"><Skeleton lines={2} /></Card>
      {/each}
    {:else}
      <MetricCard
        label="UPTIME"
        value={uptimeLabel($status.data?.uptime_minutes ?? 0)}
        variant="success"
      />
      <MetricCard
        label="MENSAGENS HOJE"
        value={String($metrics.data?.bot_messages ?? 0)}
      />
      <MetricCard
        label="MEMBROS ATIVOS (1h)"
        value={String($metrics.data?.active_users ?? 0)}
      />
      <MetricCard
        label="MEMÓRIAS HOJE"
        value={String($metrics.data?.memories_today ?? 0)}
        variant="info"
      />
    {/if}
  </div>

  <!-- Persona -->
  <Card variant="glass" class="!p-6">
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-lg font-semibold">Persona Ativa</h2>
      <Badge variant="primary">
        {$persona.data?.ActiveProfile?.Name ?? 'Bot'}
      </Badge>
    </div>
    <div class="grid gap-4 sm:grid-cols-2">
      <div>
        <div class="text-sm text-muted-foreground">Identity</div>
        <div class="font-medium">{$persona.data?.ActiveIdentity?.name ?? '—'}</div>
      </div>
      <div>
        <div class="text-sm text-muted-foreground">Overlay</div>
        <div class="font-medium">{$persona.data?.ActivePersona?.name ?? '—'}</div>
      </div>
    </div>
  </Card>

  <!-- Passive Mode Toggle -->
  <Card variant="glass" class="!p-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-lg font-semibold">Modo Passivo</h2>
        <p class="text-sm text-muted-foreground">
          Quando ativo, o bot só responde a menções diretas. Pulso ambiente desativado.
        </p>
      </div>
      <Switch checked={passiveMode} onchange={togglePassive} />
    </div>
  </Card>

  <!-- Recent Feed -->
  <Card variant="glass" class="!p-6">
    <h2 class="text-lg font-semibold mb-4">Feed Recente</h2>
    {#if $feed.loading}
      <Skeleton lines={4} />
    {:else if $feed.data?.Messages?.length > 0}
      <div class="space-y-3">
        {#each $feed.data.Messages.slice(0, 10) as msg}
          <div class="flex items-start gap-3 text-sm">
            <span class="font-medium shrink-0">{msg.author_name ?? 'Unknown'}</span>
            <span class="text-muted-foreground line-clamp-2">{msg.content ?? ''}</span>
          </div>
        {/each}
      </div>
    {:else}
      <p class="text-sm text-muted-foreground">Nenhuma mensagem recente.</p>
    {/if}
  </Card>
</div>
