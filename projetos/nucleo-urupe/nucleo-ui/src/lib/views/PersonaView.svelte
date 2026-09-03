<script>
    import PersonaPanel from '../components/PersonaPanel.svelte';
    import ContextPanel from '../components/ContextPanel.svelte';
    import CuratorBench from '../components/CuratorBench.svelte';
</script>

<div class="view-content">
    <div class="persona-grid">
        <div class="editor-col">
            <section class="section">
                <PersonaPanel />
            </section>
            <section class="section panel">
                <div class="panel-header">
                    <p class="eyebrow">Ajustes Dinâmicos</p>
                    <h3>Propostas de Mudança</h3>
                </div>
                <CuratorBench />
            </section>
        </div>
        <div class="context-col">
            <section class="section">
                <ContextPanel />
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

    .persona-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 2rem;
        align-items: start;
    }

    @media (max-width: 1200px) {
        .persona-grid {
            grid-template-columns: 1fr;
        }
    }
</style>
