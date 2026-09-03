<script lang="ts">
  import { Card, MetricCard, Skeleton } from 'bindrunes';
  import { createLiveQuery } from '$lib/api/hooks';
  import { api } from '$lib/api/client';

  const modlog = createLiveQuery(() => api.modlog(), ['feed']);
  const audit = createLiveQuery(() => api.audit(), ['feed']);
</script>

<svelte:head>
  <title>Moderação — Labirinto de Dédalo</title>
</svelte:head>

<div class="space-y-6">
  <div>
    <h1 class="text-2xl font-bold">Moderação</h1>
    <p class="text-muted-foreground">Warnings, ações e log de auditoria</p>
  </div>

  <!-- Mod Log -->
  <Card variant="glass" class="!p-6">
    <h2 class="text-lg font-semibold mb-4">Log de Moderação</h2>
    {#if $modlog.loading}
      <Skeleton lines={4} />
    {:else if $modlog.data?.entries?.length > 0}
      <div class="space-y-3">
        {#each $modlog.data.entries as entry}
          <div class="flex items-center justify-between text-sm p-3 rounded-lg bg-surface-1/50">
            <div>
              <span class="font-medium">{entry.action}</span>
              <span class="text-muted-foreground"> — {entry.reason ?? 'Sem motivo'}</span>
            </div>
            <span class="text-xs text-muted-foreground">{entry.timestamp ?? ''}</span>
          </div>
        {/each}
      </div>
    {:else}
      <p class="text-sm text-muted-foreground">Nenhuma ação de moderação registrada.</p>
    {/if}
  </Card>

  <!-- Audit Log -->
  <Card variant="glass" class="!p-6">
    <h2 class="text-lg font-semibold mb-4">Auditoria</h2>
    {#if $audit.loading}
      <Skeleton lines={4} />
    {:else if $audit.data?.entries?.length > 0}
      <div class="space-y-3">
        {#each $audit.data.entries as entry}
          <div class="flex items-center justify-between text-sm p-3 rounded-lg bg-surface-1/50">
            <div>
              <span class="font-medium">{entry.action}</span>
              <span class="text-muted-foreground"> — {entry.actor ?? 'system'} → {entry.target ?? '—'}</span>
            </div>
            <span class="text-xs text-muted-foreground">{entry.timestamp ?? ''}</span>
          </div>
        {/each}
      </div>
    {:else}
      <p class="text-sm text-muted-foreground">Nenhum registro de auditoria.</p>
    {/if}
  </Card>
</div>
