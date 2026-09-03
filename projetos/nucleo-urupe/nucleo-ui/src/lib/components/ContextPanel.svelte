<script>
  import { Button, Input, Select, Checkbox, Switch, Dialog, EmptyState, Skeleton, Badge, Alert, Spinner } from '@talos/ui';

    import { context } from '../stores';

    async function clearContext() {
        if (!confirm("Tem certeza que deseja zerar o contexto persistente? O bot perderá a memória de longo prazo deste canal.")) return;
        
        try {
            await fetch('/api/context/clear', { method: 'POST' });
        } catch(err) {
            console.error("Erro ao limpar contexto", err);
        }
    }
</script>

<div class="panel">
    <div class="panel-header">
        <h2>Janela de Contexto</h2>
        <Button class="danger-button" onclick={clearContext}>Zerar Resumo</Button>
    </div>
    
    <div class="context-grid">
        {#if $context}
            <div class="context-card">
                <p class="eyebrow">Contexto Base</p>
                <div class="context-copy">{$context.BaseContextDetail}</div>
                <div class="context-foot">Atualizado: {$context.LastBaseContextAt}</div>
            </div>
            
            <div class="context-card">
                <p class="eyebrow">Resumo Rolante ({$context.SummaryLabel})</p>
                <div class="context-copy">{$context.SummaryDetail}</div>
                <div class="context-foot">Atualizado: {$context.LastSummaryAt}</div>
            </div>

            <div class="context-card">
                <p class="eyebrow">Estado Operacional</p>
                <ul class="context-list">
                    <li>Última Resposta: {$context.LastResponseAt}</li>
                    <li>Modelo em Uso: {$context.LastModel}</li>
                    <li>Último Alerta: {$context.LastError} ({$context.LastErrorAt})</li>
                    <li>Sincronização: {$context.LastProfileSync}</li>
                </ul>
            </div>
        {:else}
            <p class="empty-state">Carregando contexto...</p>
        {/if}
    </div>
</div>

<style>
    .panel-header {
        padding: 1.5rem 1.6rem 1.2rem;
        border-bottom: 1px solid var(--border);
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .context-grid {
        padding: 1.5rem;
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
        gap: 1.5rem;
    }

    .context-card {
        background: rgba(0,0,0,0.1);
        border: 1px solid var(--border);
        border-radius: 12px;
        padding: 1.2rem;
    }

    .context-copy {
        color: var(--muted-foreground);
        font-size: 0.95rem;
        line-height: 1.5;
        margin-bottom: 1rem;
    }

    .context-foot {
        font-size: 0.8rem;
        color: var(--muted-foreground);
        font-family: var(--font-mono);
    }

    .context-list {
        list-style: none;
        padding: 0;
        margin: 0;
        color: var(--muted-foreground);
        font-size: 0.95rem;
        line-height: 1.8;
    }

    .danger-button {
        background: transparent;
        border: 1px solid var(--destructive);
        color: var(--destructive);
        padding: 0.5rem 1rem;
        border-radius: 8px;
        cursor: pointer;
        transition: all 0.2s ease;
    }
    .danger-button:hover {
        background: var(--destructive);
        color: #fff;
    }
</style>
