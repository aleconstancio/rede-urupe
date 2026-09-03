<script>
  import { Button, Input, Select, Checkbox, Switch, Dialog, EmptyState, Skeleton, Badge, Alert, Spinner } from '@talos/ui';

    import { onMount } from 'svelte';
    import { persona } from '../stores';

    let { channelId = "" } = $props();
    
    let proposals = [];
    let loading = true;

    async function fetchProposals() {
        loading = true;
        try {
            const res = await fetch(`/api/persona-proposals?channel_id=${channelId}&status=pending`);
            if (res.ok) {
                const data = await res.json();
                proposals = data.items || [];
            }
        } catch (err) {
            console.error("Failed to fetch proposals", err);
        } finally {
            loading = false;
        }
    }

    async function handleAction(id, action) {
        try {
            const res = await fetch(`/api/persona-proposals/${id}/${action}`, { method: 'POST' });
            if (res.ok) {
                await fetchProposals();
            }
        } catch (err) {
            console.error(`Failed to ${action} proposal`, err);
        }
    }

    onMount(fetchProposals);
</script>

<div class="reviewer">
    <div class="section-header">
        <h3>Propostas de Adaptação (System 3)</h3>
        <Button class="refresh-btn" onclick={fetchProposals} title="Atualizar">
            <svg viewBox="0 0 24 24" width="16" height="16"><path fill="currentColor" d="M17.65 6.35A7.958 7.958 0 0 0 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0 1 12 18c-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"/></svg>
        </Button>
    </div>

    {#if loading}
        <PageLoading type="text" lines={3} />
    {:else if proposals.length === 0}
        <div class="empty">Nenhuma proposta pendente. O System 3 está em harmonia com o estilo atual.</div>
    {:else}
        <div class="proposal-list">
            {#each proposals as p}
                <div class="proposal-card">
                    <div class="card-header">
                        <span class="target-tag">{p.Target}</span>
                        <span class="confidence">{(p.Confidence * 100).toFixed(0)}% confiança</span>
                    </div>
                    
                    <div class="reason">
                        <strong>Motivo:</strong> {p.Reason}
                    </div>

                    <div class="changes">
                        <strong>Mudanças sugeridas:</strong>
                        <pre>{JSON.stringify(p.ProposedChanges, null, 2)}</pre>
                    </div>

                    {#if p.EvidenceMessageIDs && p.EvidenceMessageIDs.length > 0}
                        <div class="evidence">
                            <strong>Evidência:</strong> {p.EvidenceMessageIDs.length} mensagens citadas
                        </div>
                    {/if}

                    <div class="actions">
                        <Button class="btn approve" onclick={() => handleAction(p.ID, 'approve')}>Aprovar</Button>
                        <Button class="btn reject" onclick={() => handleAction(p.ID, 'reject')}>Rejeitar</Button>
                        <Button class="btn apply" onclick={() => handleAction(p.ID, 'apply')}>Aplicar Agora</Button>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

<style>
    .reviewer {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .section-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    h3 {
        margin: 0;
        font-size: 0.9rem;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        color: var(--muted-foreground);
    }

    .refresh-btn {
        background: none;
        border: none;
        color: var(--muted-foreground);
        cursor: pointer;
        padding: 4px;
        display: flex;
        align-items: center;
        border-radius: 4px;
    }

    .refresh-btn:hover {
        color: var(--primary);
        background: var(--muted);
    }

    .proposal-list {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .proposal-card {
        background: rgba(255, 255, 255, 0.03);
        border: 1px solid var(--border);
        border-radius: 12px;
        padding: 1rem;
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }

    .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .target-tag {
        background: var(--muted);
        color: var(--primary);
        padding: 2px 8px;
        border-radius: 4px;
        font-size: 0.75rem;
        font-weight: 700;
        text-transform: uppercase;
    }

    .confidence {
        font-size: 0.8rem;
        color: var(--muted-foreground);
    }

    .reason {
        font-size: 0.9rem;
        line-height: 1.4;
    }

    .changes pre {
        background: rgba(0,0,0,0.2);
        padding: 0.5rem;
        border-radius: 6px;
        font-size: 0.8rem;
        margin: 0.5rem 0 0;
        overflow-x: auto;
    }

    .actions {
        display: flex;
        gap: 0.5rem;
        margin-top: 0.5rem;
    }

    .btn {
        flex: 1;
        padding: 0.5rem;
        border-radius: 6px;
        font-size: 0.85rem;
        font-weight: 600;
        cursor: pointer;
        border: 1px solid transparent;
        transition: all 0.2s;
    }

    .approve {
        background: rgba(83, 212, 143, 0.1);
        color: var(--success);
        border-color: rgba(83, 212, 143, 0.2);
    }

    .approve:hover {
        background: var(--success);
        color: white;
    }

    .reject {
        background: rgba(235, 109, 98, 0.1);
        color: var(--destructive);
        border-color: rgba(235, 109, 98, 0.2);
    }

    .reject:hover {
        background: var(--destructive);
        color: white;
    }

    .apply {
        background: var(--primary);
        color: var(--background);
    }

    .apply:hover {
        filter: brightness(1.1);
    }

    .empty, .loading {
        padding: 2rem;
        text-align: center;
        color: var(--muted-foreground);
        font-size: 0.9rem;
        background: rgba(255,255,255,0.02);
        border-radius: 12px;
        border: 1px dashed var(--border);
    }
</style>
