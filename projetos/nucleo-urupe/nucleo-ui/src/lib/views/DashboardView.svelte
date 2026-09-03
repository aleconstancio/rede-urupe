<script>
    import MetricsPanel from '../components/MetricsPanel.svelte';
    import FeedPanel from '../components/FeedPanel.svelte';
    import { persona } from '../stores';
</script>

<div class="view-content">
    <section class="section">
        <MetricsPanel />
    </section>

    <div class="dashboard-grid">
        <div class="main-col">
            {#if $persona?.LastThought}
            <section class="section cognitive-status">
                <div class="status-bar">
                    <div class="pulse-indicator"></div>
                    <span class="status-label">PENSAMENTO ATIVO:</span>
                    <span class="status-text">"{$persona.LastThought}"</span>
                </div>
            </section>
            {/if}

            <section class="section panel">
                <div class="panel-header">
                    <p class="eyebrow">Atividade Recente</p>
                    <h3>Resumo do Feed</h3>
                </div>
                <div class="preview-container">
                    <!-- We'll reuse FeedPanel but with a limit if possible, or just show it as is -->
                    <FeedPanel limit={15} />
                </div>
            </section>
        </div>

        <div class="side-col">
            {#if $persona?.LastThought}
            <section class="section panel thought-panel">
                <div class="panel-header">
                    <p class="eyebrow">Estado Interno</p>
                    <h3>Última Reflexão</h3>
                </div>
                <div class="thought-content">
                    "{$persona.LastThought}"
                </div>
            </section>
            {/if}

            <section class="section panel">
                <div class="panel-header">
                    <p class="eyebrow">Dicas</p>
                    <h3>Atalhos</h3>
                </div>
                <div class="shortcuts">
                    <p class="note">Os <b>Arquivos Cognitivos</b> contêm o histórico de longo prazo e fatos consolidados.</p>
                    <p class="note">No <b>Estúdio de Persona</b>, você pode ajustar a identidade e aprovar mudanças de postura.</p>
                </div>
            </section>
        </div>
    </div>
</div>

<style>
    .view-content {
        padding: 2rem;
        max-width: 1600px;
        margin: 0 auto;
    }

    .section {
        margin-bottom: 2rem;
    }

    .panel-header {
        padding: 1.5rem 1.5rem 0.5rem;
    }

    .dashboard-grid {
        display: grid;
        grid-template-columns: 1fr 350px;
        gap: 2rem;
    }

    .preview-container {
        height: 500px;
        overflow: hidden;
    }

    .thought-panel {
        border-color: var(--primary);
        background: var(--muted);
    }

    .thought-content {
        padding: 1rem 1.5rem 1.5rem;
        font-style: italic;
        color: var(--muted-foreground);
        line-height: 1.6;
    }

    .cognitive-status {
        margin-bottom: 1.5rem;
        background: rgba(255, 191, 0, 0.05);
        border: 1px solid rgba(255, 191, 0, 0.2);
        border-radius: 12px;
        padding: 0.75rem 1.25rem;
    }

    .status-bar {
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    .pulse-indicator {
        width: 8px;
        height: 8px;
        background: var(--primary);
        border-radius: 50%;
        box-shadow: 0 0 10px var(--primary);
        animation: status-pulse 2s infinite;
    }

    @keyframes status-pulse {
        0% { transform: scale(1); opacity: 1; }
        50% { transform: scale(1.3); opacity: 0.6; }
        100% { transform: scale(1); opacity: 1; }
    }

    .status-label {
        font-size: 0.65rem;
        font-weight: 800;
        letter-spacing: 0.1em;
        color: var(--primary);
    }

    .status-text {
        font-size: 0.9rem;
        font-style: italic;
        color: var(--muted-foreground);
        font-family: var(--font-serif);
    }

    .shortcuts {
        padding: 1rem 1.5rem 1.5rem;
    }

    .note {
        font-size: 0.9rem;
        color: var(--muted-foreground);
        margin-bottom: 1rem;
        line-height: 1.4;
    }

    @media (max-width: 1200px) {
        .dashboard-grid {
            grid-template-columns: 1fr;
        }
    }
</style>
