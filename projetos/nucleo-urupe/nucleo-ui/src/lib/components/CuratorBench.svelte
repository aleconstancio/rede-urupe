<script>
  import { Button, Input, Select, Checkbox, Switch, Dialog, EmptyState, Skeleton, Badge, Alert, Spinner } from '@talos/ui';

    import { annotations } from '../stores';
    import { 
        Shield, Tag, Fingerprint, MessageSquare, 
        AlertCircle, ExternalLink, Filter, Search
    } from 'lucide-svelte';
    import { fade } from 'svelte/transition';

    /** @type {any[]} */
    let items = $derived($annotations?.items || []);
    let searchQuery = $state("");

    /** @param {number} score */
    function getConfidenceColor(score) {
        if (score > 0.8) return 'var(--success)';
        if (score > 0.5) return 'var(--warning)';
        return 'var(--muted-foreground)';
    }

    /** @param {number} score */
    function formatScore(score) {
        return (score * 100).toFixed(0) + '%';
    }

    let filteredItems = $derived(items.filter(ann => {
        const text = [
            ann.author_id,
            ...(ann.topic_tags || []),
            ...(ann.stance_tags || []),
            ann.evidence_type
        ].join(' ').toLowerCase();
        return text.includes(searchQuery.toLowerCase());
    }));
</script>

<div class="curator-bench">
    <header class="bench-header">
        <div class="title-area">
            <div class="icon-box">
                <Shield size={20} />
            </div>
            <div>
                <h1>Evidence Layer</h1>
                <p class="subtitle">Inspeção de `message_annotations` (System 3)</p>
            </div>
        </div>

        <div class="search-bar">
            <Search size={16} class="search-icon" />
            <input type="text" placeholder="Filtrar por autor, tag ou evidência..." bind:value={searchQuery} />
        </div>
    </header>

    <div class="stats-row">
        <div class="stat-card">
            <span class="label">Total de Anotações</span>
            <span class="value">{items.length}</span>
        </div>
        <div class="stat-card">
            <span class="label">Pessoas Monitoradas</span>
            <span class="value">{new Set(items.map(a => a.author_id)).size}</span>
        </div>
        <div class="stat-card">
            <span class="label">Alta Confiança (>80%)</span>
            <span class="value">{items.filter(a => a.durability_score > 0.8).length}</span>
        </div>
    </div>

    <div class="evidence-grid">
        {#each filteredItems as ann}
            <div class="evidence-card" transition:fade>
                <div class="card-header">
                    <div class="author-info">
                        <Fingerprint size={14} class="icon-accent" />
                        <span class="author">@{ann.author_id.slice(-6)}</span>
                    </div>
                    <div class="confidence-badge" style="color: {getConfidenceColor(ann.durability_score)}">
                        {formatScore(ann.durability_score)} conf
                    </div>
                </div>

                <div class="tags">
                    {#each ann.topic_tags || [] as tag}
                        <span class="tag topic">#{tag}</span>
                    {/each}
                    {#each ann.stance_tags || [] as tag}
                        <span class="tag stance">{tag}</span>
                    {/each}
                </div>

                <div class="meta">
                    <div class="meta-item">
                        <Tag size={12} />
                        <span>{ann.evidence_type}</span>
                    </div>
                    <div class="meta-item">
                        <MessageSquare size={12} />
                        <span>MID: {ann.message_id.slice(-6)}</span>
                    </div>
                </div>

                <div class="scores-grid">
                    <div class="score">
                        <span class="score-label">Humor</span>
                        <div class="score-bar"><div class="fill" style="width: {ann.humor_score * 100}%"></div></div>
                    </div>
                    <div class="score">
                        <span class="score-label">Sarcasmo</span>
                        <div class="score-bar"><div class="fill" style="width: {ann.sarcasm_score * 100}%"></div></div>
                    </div>
                </div>

                <div class="card-footer">
                    <Button class="btn-inspect">
                        <ExternalLink size={14} />
                        Ver Contexto
                    </Button>
                </div>
            </div>
        {/each}
    </div>

    {#if filteredItems.length === 0}
        <div class="empty-state">
            <AlertCircle size={48} />
            <p>Nenhuma anotação encontrada para os critérios atuais.</p>
        </div>
    {/if}
</div>

<style>
    .curator-bench {
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
        padding: 2rem;
    }

    .bench-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .title-area {
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    .icon-box {
        width: 48px;
        height: 48px;
        background: var(--muted);
        color: var(--primary);
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    h1 {
        margin: 0;
        font-size: 1.5rem;
    }

    .subtitle {
        margin: 0;
        font-size: 0.9rem;
        color: var(--muted-foreground);
    }

    .search-bar {
        position: relative;
        width: 350px;
    }

    :global(.search-icon) {
        position: absolute;
        left: 1rem;
        top: 50%;
        transform: translateY(-50%);
        color: var(--muted-foreground);
    }

    input {
        width: 100%;
        background: var(--card);
        border: 1px solid var(--border);
        padding: 0.75rem 1rem 0.75rem 2.75rem;
        border-radius: 10px;
        color: var(--foreground);
        outline: none;
    }

    .stats-row {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 1.5rem;
    }

    .stat-card {
        background: rgba(255,255,255,0.02);
        border: 1px solid var(--border);
        padding: 1.25rem;
        border-radius: 12px;
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .stat-card .label {
        font-size: 0.8rem;
        text-transform: uppercase;
        color: var(--muted-foreground);
        letter-spacing: 0.05em;
    }

    .stat-card .value {
        font-size: 1.5rem;
        font-weight: 700;
        color: var(--primary);
    }

    .evidence-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
        gap: 1rem;
    }

    .evidence-card {
        background: var(--card);
        border: 1px solid var(--border);
        border-radius: 14px;
        padding: 1.25rem;
        display: flex;
        flex-direction: column;
        gap: 1rem;
        transition: transform 0.2s;
    }

    .evidence-card:hover {
        transform: translateY(-4px);
        border-color: var(--muted);
    }

    .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .author-info {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .author {
        font-family: var(--font-mono);
        font-weight: 700;
        font-size: 0.9rem;
    }

    .confidence-badge {
        font-size: 0.75rem;
        font-weight: 800;
        text-transform: uppercase;
    }

    .tags {
        display: flex;
        flex-wrap: wrap;
        gap: 0.4rem;
    }

    .tag {
        font-size: 0.75rem;
        font-weight: 600;
        padding: 2px 8px;
        border-radius: 4px;
    }

    .tag.topic { background: rgba(59, 130, 246, 0.1); color: #60a5fa; }
    .tag.stance { background: rgba(168, 85, 247, 0.1); color: #c084fc; }

    .meta {
        display: flex;
        gap: 1rem;
    }

    .meta-item {
        display: flex;
        align-items: center;
        gap: 0.4rem;
        font-size: 0.75rem;
        color: var(--muted-foreground);
    }

    .scores-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1rem;
    }

    .score-label {
        font-size: 0.7rem;
        color: var(--muted-foreground);
        margin-bottom: 0.25rem;
        display: block;
    }

    .score-bar {
        height: 4px;
        background: rgba(255,255,255,0.05);
        border-radius: 2px;
        overflow: hidden;
    }

    .fill {
        height: 100%;
        background: var(--primary);
    }

    .card-footer {
        margin-top: 0.5rem;
    }

    .btn-inspect {
        width: 100%;
        background: rgba(255,255,255,0.03);
        border: 1px solid var(--border);
        color: var(--muted-foreground);
        padding: 0.5rem;
        border-radius: 6px;
        font-size: 0.8rem;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        transition: all 0.2s;
    }

    .btn-inspect:hover {
        background: var(--muted);
        color: var(--primary);
        border-color: var(--primary);
    }

    .empty-state {
        text-align: center;
        padding: 5rem;
        color: var(--muted-foreground);
        opacity: 0.5;
    }
</style>
