<script>
    import { feed } from '../stores';
    let { limit = null } = $props();

    let feedValue = $derived($feed);
    let messages = $derived(limit && feedValue?.Messages ? feedValue.Messages.slice(0, limit) : feedValue?.Messages);
</script>

<div class="feed-panel-container">
    <div class="panel-header">
        <div>
            <h2>Feed Ao Vivo</h2>
            {#if feedValue}
                <span class="count-tag">{feedValue.CountLabel}</span>
            {/if}
        </div>
    </div>
    <div class="transcript-list">
        {#if feedValue && messages && messages.length > 0}
            {#each messages as msg}
                <div class="entry" class:entry-bot={msg.is_bot} class:entry-bot-latest={msg.is_latest_bot}>
                    <div class="entry-meta">
                        <span class="entry-author">@{msg.author || 'Desconhecido'}</span>
                        <span class="entry-badge" class:bot-badge={msg.is_bot}>{msg.is_bot ? 'TALOS' : 'MEMBRO'}</span>
                        {#if msg.category && msg.category !== 'uncategorized'}
                            <span class="entry-archive">📦 {msg.category}</span>
                        {/if}
                        <span class="entry-time">
                            {#if new Date(msg.timestamp).toDateString() !== new Date().toDateString()}
                                {new Date(msg.timestamp).toLocaleDateString('pt-BR', {day: '2-digit', month: 'short'})} 
                            {/if}
                            {new Date(msg.timestamp).toLocaleTimeString('pt-BR', {hour: '2-digit', minute:'2-digit'})}
                        </span>
                    </div>
                    <div class="entry-content">
                        {#if msg.content && msg.content.trim() !== ""}
                            {msg.content}
                        {:else if msg.attachments && msg.attachments.length > 0}
                            <span class="attachment-fallback">[Arquivo Anexo]</span>
                        {:else}
                            <span class="empty-fallback">(Mensagem sem texto)</span>
                        {/if}
                    </div>
                    
                    {#if msg.internal_monologue}
                        <div class="entry-monologue">
                            <span class="monologue-header">Cognição Interna</span>
                            {#if msg.internal_monologue.startsWith('{')}
                                {@const m = JSON.parse(msg.internal_monologue)}
                                <div class="monologue-grid">
                                    {#if m.surface_read}<div class="m-row"><b>Leitura:</b> {m.surface_read}</div>{/if}
                                    {#if m.subtext}<div class="m-row"><b>Subtexto:</b> {m.subtext}</div>{/if}
                                    {#if m.vibe_plan}<div class="m-row"><b>Vibe:</b> {m.vibe_plan}</div>{/if}
                                    {#if m.strategic_frame}<div class="m-row"><b>Frame:</b> {m.strategic_frame}</div>{/if}
                                </div>
                            {:else}
                                <div class="m-row">{msg.internal_monologue}</div>
                            {/if}
                        </div>
                    {/if}

                    {#if msg.grounding_ledger}
                        <div class="entry-grounding">
                            <span class="ledger-header">Ancoragem (Memória)</span>
                            {#if msg.grounding_ledger.startsWith('[')}
                                {@const ledger = JSON.parse(msg.grounding_ledger)}
                                <div class="ledger-tags">
                                    {#each ledger as item}
                                        <span class="ledger-tag">{item}</span>
                                    {/each}
                                </div>
                            {:else}
                                <div class="ledger-text">{msg.grounding_ledger}</div>
                            {/if}
                        </div>
                    {/if}
                </div>
            {/each}
        {:else if feedValue}
            <div class="empty-state">{feedValue.EmptyText}</div>
        {:else}
            <PageLoading type="text" lines={3} />
        {/if}
    </div>
</div>

<style>
    .feed-panel-container {
        display: flex;
        flex-direction: column;
        height: 100%;
        background: var(--card);
        border-radius: var(--radius);
        border: 1px solid var(--border);
        overflow: hidden;
    }

    .panel-header {
        padding: 1.5rem 1.6rem 1.2rem;
        border-bottom: 1px solid var(--border);
        display: flex;
        justify-content: space-between;
        align-items: center;
        background: rgba(150, 150, 150, 0.02);
    }

    .count-tag {
        font-size: 0.75rem;
        color: var(--muted-foreground);
        background: var(--border);
        padding: 0.2rem 0.6rem;
        border-radius: 999px;
        margin-left: 0.75rem;
    }

    .transcript-list {
        padding: 1.5rem;
        display: grid;
        gap: 1rem;
        flex: 1;
        overflow-y: auto;
    }

    .entry {
        background: rgba(150, 150, 150, 0.05);
        border: 1px solid var(--border);
        border-radius: 12px;
        padding: 1.2rem;
        transition: transform 0.2s ease;
    }

    .entry:hover {
        transform: translateX(4px);
        background: rgba(150, 150, 150, 0.08);
    }

    .entry-bot {
        background: var(--muted);
        border-color: var(--primary);
    }

    .entry-meta {
        display: flex;
        align-items: center;
        gap: 0.7rem;
        margin-bottom: 0.75rem;
    }

    .entry-author {
        font-weight: 700;
        color: var(--foreground);
        font-size: 0.9rem;
    }

    .entry-badge {
        font-size: 0.65rem;
        padding: 0.2rem 0.5rem;
        border-radius: 4px;
        background: rgba(255,255,255,0.08);
        text-transform: uppercase;
        letter-spacing: 0.05em;
        font-weight: 600;
        color: var(--muted-foreground);
    }

    .bot-badge {
        background: var(--primary);
        color: #fff;
    }

    .entry-archive {
        font-size: 0.65rem;
        color: var(--muted-foreground);
        background: rgba(255, 255, 255, 0.05);
        padding: 0.15rem 0.5rem;
        border-radius: 4px;
        font-weight: 500;
        text-transform: capitalize;
    }

    .entry-time {
        font-size: 0.75rem;
        color: var(--muted-foreground);
        font-family: var(--font-mono);
        margin-left: auto;
    }

    .entry-content {
        white-space: pre-wrap;
        line-height: 1.6;
        color: var(--foreground);
        font-size: 0.95rem;
    }

    .attachment-fallback, .empty-fallback {
        color: var(--primary);
        font-style: italic;
        font-size: 0.85rem;
        opacity: 0.8;
    }

    .entry-monologue {
        margin-top: 1rem;
        padding: 1rem;
        background: rgba(0, 0, 0, 0.3);
        border: 1px solid rgba(255, 191, 0, 0.1);
        border-radius: 10px;
        font-size: 0.85rem;
        color: var(--muted-foreground);
        border-left: 3px solid var(--primary);
    }

    .monologue-header {
        display: block;
        font-size: 0.65rem;
        text-transform: uppercase;
        font-weight: 800;
        letter-spacing: 0.1em;
        color: var(--primary);
        margin-bottom: 0.75rem;
        opacity: 0.8;
    }

    .m-row {
        margin-bottom: 0.5rem;
        line-height: 1.5;
    }

    .m-row b {
        color: var(--foreground);
        font-weight: 600;
        margin-right: 0.5rem;
    }

    .entry-grounding {
        margin-top: 0.75rem;
        padding: 0.75rem 1rem;
        background: rgba(255, 255, 255, 0.02);
        border-radius: 8px;
        border: 1px dashed var(--border);
    }

    .ledger-header {
        display: block;
        font-size: 0.6rem;
        text-transform: uppercase;
        font-weight: 700;
        color: var(--muted-foreground);
        margin-bottom: 0.5rem;
    }

    .ledger-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 0.4rem;
    }

    .ledger-tag {
        font-size: 0.75rem;
        background: rgba(255, 191, 0, 0.05);
        color: var(--primary);
        padding: 2px 8px;
        border-radius: 4px;
        border: 1px solid rgba(255, 191, 0, 0.2);
    }

    .ledger-text {
        font-size: 0.8rem;
        color: var(--muted-foreground);
        font-style: italic;
    }

    .empty-state {
        padding: 3rem;
        text-align: center;
        color: var(--muted-foreground);
        font-style: italic;
    }
</style>

