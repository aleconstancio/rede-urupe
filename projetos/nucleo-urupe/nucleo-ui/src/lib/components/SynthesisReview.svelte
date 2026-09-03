<script>
  import { Button, Input, Select, Checkbox, Switch, Dialog, EmptyState, Skeleton, Badge, Alert, Spinner } from '@talos/ui';

    import { onMount } from 'svelte';
    import { Check, Clock, BrainCircuit, Hash, User, AlertCircle } from 'lucide-svelte';
    import { fade, slide } from 'svelte/transition';

    let pendingMemories = [];
    let loading = true;
    let error = null;
    let processingId = null;

    async function fetchPending() {
        try {
            loading = true;
            const res = await fetch('/api/synthesis/pending');
            if (res.ok) {
                pendingMemories = await res.json();
            } else {
                error = "Falha ao carregar sínteses pendentes.";
            }
        } catch (err) {
            error = "Erro de conexão com o servidor.";
        } finally {
            loading = false;
        }
    }

    async function approve(segmentKey) {
        try {
            processingId = segmentKey;
            const res = await fetch('/api/synthesis/approve', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ segment_key: segmentKey })
            });

            if (res.ok) {
                pendingMemories = pendingMemories.filter(m => m.SegmentKey !== segmentKey);
            } else {
                alert("Erro ao aprovar síntese.");
            }
        } catch (err) {
            console.error(err);
        } finally {
            processingId = null;
        }
    }

    onMount(fetchPending);
</script>

<div class="review-container">
    <header class="header">
        <div class="title-group">
            <p class="eyebrow">Cognição Hierárquica</p>
            <h1>Revisão de Inteligência</h1>
        </div>
        <Button class="refresh-btn" onclick={fetchPending} disabled={loading}>
            <span class:spinning={loading}>
                <Clock size={16} />
            </span>
            Atualizar
        </Button
    </header>

    {#if loading && pendingMemories.length === 0}
        <div class="loading-state" in:fade>
            <div class="spinner"></div>
            <p>Sincronizando grafo de contexto...</p>
        </div>
    {:else if error}
        <div class="error-state" in:fade>
            <AlertCircle size={48} color="var(--destructive)" />
            <p>{error}</p>
        </div>
    {:else if pendingMemories.length === 0}
        <div class="empty-state" in:fade>
            <BrainCircuit size={64} color="var(--muted-foreground)" />
            <h2>Sem Pendências</h2>
            <p>O grafo de contexto está totalmente validado. Novas impressões aparecerão conforme o bot processar novas conversas.</p>
        </div>
    {:else}
        <div class="grid">
            {#each pendingMemories as memory (memory.SegmentKey)}
                <article class="memory-card" in:slide out:fade>
                    <div class="card-header">
                        <div class="topic-info">
                            <span class="macro">#{memory.Topic}</span>
                            <h3>{memory.Topic}</h3>
                        </div>
                        <div class="status-badge">PENDENTE</div>
                    </div>

                    <div class="card-body">
                        <div class="detail-row">
                            <Hash size={14} />
                            <span>Segmento: <code>{memory.SegmentKey}</code></span>
                        </div>
                        
                        {#if memory.Note}
                            <p class="note">{memory.Note}</p>
                        {/if}

                        {#if memory.Positions}
                            <div class="positions">
                                <p class="label">Posições Detectadas:</p>
                                <div class="positions-list">
                                    {memory.Positions}
                                </div>
                            </div>
                        {/if}

                        {#if memory.MetaAnalysis}
                            <div class="meta-analysis">
                                <p class="label">Padrão Meta-Analítico:</p>
                                <div class="meta-content">{memory.MetaAnalysis}</div>
                            </div>
                        {/if}
                    </div>

                    <div class="card-footer">
                        <div class="dynamics">
                            <span class="tag">{memory.Dynamics}</span>
                            <span class="tag">{memory.State}</span>
                        </div>
                        <Button 
                            class="approve-btn" 
                            disabled={processingId === memory.SegmentKey}
                            onclick={() => approve(memory.SegmentKey)}
                        >
                            {#if processingId === memory.SegmentKey}
                                <span class="spinner-small"></span>
                            {:else}
                                <Check size={18} />
                            {/if}
                            Validar Memória
                        </Button
                    </div>
                </article>
            {/each}
        </div>
    {/if}
</div>

<style>
    .review-container {
        padding: 2rem;
        max-width: 1400px;
        margin: 0 auto;
    }

    .header {
        display: flex;
        justify-content: space-between;
        align-items: flex-end;
        margin-bottom: 3rem;
    }

    h1 {
        font-family: var(--font-serif);
        font-size: 2.5rem;
        color: var(--foreground);
        margin: 0;
    }

    .refresh-btn {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.6rem 1rem;
        background: var(--glass-surface);
        border: 1px solid var(--border);
        border-radius: 8px;
        color: var(--muted-foreground);
        font-size: 0.9rem;
        cursor: pointer;
        transition: all 0.2s;
    }

    .refresh-btn:hover {
        color: var(--foreground);
        border-color: var(--primary);
    }

    .grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(450px, 1fr));
        gap: 2rem;
    }

    .memory-card {
        background: var(--glass-surface);
        border: 1px solid var(--border);
        border-radius: 20px;
        overflow: hidden;
        display: flex;
        flex-direction: column;
        transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        backdrop-filter: blur(10px);
    }

    .memory-card:hover {
        transform: translateY(-4px);
        border-color: var(--muted);
    }

    .card-header {
        padding: 1.5rem;
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        border-bottom: 1px solid var(--border);
    }

    .macro {
        font-size: 0.75rem;
        font-weight: 700;
        text-transform: uppercase;
        color: var(--primary);
        letter-spacing: 0.1em;
        display: block;
        margin-bottom: 0.25rem;
    }

    h3 {
        margin: 0;
        font-size: 1.25rem;
        color: var(--foreground);
    }

    .status-badge {
        font-size: 0.7rem;
        font-weight: 800;
        padding: 0.3rem 0.6rem;
        background: rgba(234, 179, 8, 0.1);
        color: #eab308;
        border-radius: 6px;
        border: 1px solid rgba(234, 179, 8, 0.2);
    }

    .card-body {
        padding: 1.5rem;
        flex: 1;
    }

    .detail-row {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.8rem;
        color: var(--muted-foreground);
        margin-bottom: 1.2rem;
    }

    code {
        background: var(--card);
        padding: 0.1rem 0.3rem;
        border-radius: 4px;
        font-family: var(--font-mono);
    }

    .note {
        font-size: 1rem;
        line-height: 1.6;
        color: var(--muted-foreground);
        margin-bottom: 1.5rem;
    }

    .label {
        font-size: 0.75rem;
        font-weight: 700;
        text-transform: uppercase;
        color: var(--muted-foreground);
        margin-bottom: 0.5rem;
    }

    .positions-list {
        font-size: 0.9rem;
        color: var(--muted-foreground);
        padding: 0.75rem;
        background: var(--card);
        border-radius: 10px;
        margin-bottom: 1.5rem;
    }

    .meta-content {
        font-size: 0.9rem;
        color: var(--primary);
        padding: 0.75rem;
        border-left: 2px solid var(--primary);
        background: var(--muted);
        border-radius: 0 8px 8px 0;
    }

    .card-footer {
        padding: 1.25rem 1.5rem;
        background: rgba(255, 255, 255, 0.02);
        display: flex;
        justify-content: space-between;
        align-items: center;
        border-top: 1px solid var(--border);
    }

    .dynamics {
        display: flex;
        gap: 0.5rem;
    }

    .tag {
        font-size: 0.75rem;
        padding: 0.2rem 0.5rem;
        background: var(--border);
        color: var(--muted-foreground);
        border-radius: 4px;
    }

    .approve-btn {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.6rem 1.2rem;
        background: var(--primary);
        color: white;
        border: none;
        border-radius: 10px;
        font-weight: 600;
        font-size: 0.9rem;
        cursor: pointer;
        transition: all 0.2s;
    }

    .approve-btn:hover:not(:disabled) {
        filter: brightness(1.1);
        transform: scale(1.02);
    }

    .approve-btn:disabled {
        opacity: 0.7;
        cursor: not-allowed;
    }

    .loading-state, .empty-state, .error-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 8rem 2rem;
        text-align: center;
    }

    .spinner {
        width: 40px;
        height: 40px;
        border: 3px solid var(--border);
        border-top-color: var(--primary);
        border-radius: 50%;
        animation: spinning 1s linear infinite;
        margin-bottom: 1.5rem;
    }

    .spinning {
        display: inline-flex;
        animation: spinning 1s linear infinite;
    }

    @keyframes spinning {
        to { transform: rotate(360deg); }
    }

    @media (max-width: 600px) {
        .grid {
            grid-template-columns: 1fr;
        }
    }
</style>
