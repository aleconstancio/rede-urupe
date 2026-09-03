<script>
    import { Radio, Share2, Send, Flame, MessageSquare, ExternalLink, Filter } from 'lucide-svelte';

    let trendingTopics = [
        { topic: 'Soberania Alimentar', volume: '1.4k mídias', heat: 'Alto', sentiment: 'Positivo' },
        { topic: 'Passe Livre Estudantil', volume: '3.2k mídias', heat: 'Crítico', sentiment: 'Mobilização' },
        { topic: 'Denúncia de Desmatamento', volume: '890 mídias', heat: 'Urgente', sentiment: 'Alerta' }
    ];

    let actionPackages = [
        { id: 1, title: 'Campanha Passe Livre Já', category: 'Mobilização Urbana', status: 'Disparado P2P', reaches: 4200 },
        { id: 2, title: 'Pack de Memes: Contra a Especulação de Terras', category: 'Memética', status: 'Em Produção', reaches: 1800 }
    ];
</script>

<div class="spore-ops-view">
    <div class="view-header">
        <div>
            <p class="eyebrow">Módulo Spore Ops 🍄</p>
            <h2 class="font-serif text-2xl font-bold">Central de Agitprop & Guerrilha Digital</h2>
        </div>
        <button class="btn btn-primary">
            <Send size={18} />
            <span>Disparar Nova Campanha P2P</span>
        </button>
    </div>

    <div class="ops-grid">
        <!-- Clipping & Trends Panel -->
        <div class="panel">
            <div class="panel-header">
                <h3>Radar de Pautas & Clipping de Mídias</h3>
                <span class="badge">Ao Vivo</span>
            </div>
            <div class="panel-body">
                <div class="topics-list">
                    {#each trendingTopics as t}
                        <div class="topic-card">
                            <div class="topic-info">
                                <div class="topic-heat">
                                    <Flame size={16} class="text-amber" />
                                    <span class="heat-tag">{t.heat}</span>
                                </div>
                                <h4>{t.topic}</h4>
                                <span class="vol">{t.volume}</span>
                            </div>
                            <button class="btn-sm btn-outline">Gerar Pauta</button>
                        </div>
                    {/each}
                </div>
            </div>
        </div>

        <!-- Action Packages Panel -->
        <div class="panel">
            <div class="panel-header">
                <h3>Pacotes de Ação em Andamento</h3>
            </div>
            <div class="panel-body">
                <div class="packages-list">
                    {#each actionPackages as pkg}
                        <div class="package-card">
                            <div class="pkg-header">
                                <span class="pkg-cat">{pkg.category}</span>
                                <span class="pkg-status">{pkg.status}</span>
                            </div>
                            <h4>{pkg.title}</h4>
                            <div class="pkg-footer">
                                <span>Alcançados: <b>{pkg.reaches} militantes</b></span>
                                <Share2 size={16} />
                            </div>
                        </div>
                    {/each}
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    .spore-ops-view {
        padding: 2rem;
        max-width: 1600px;
        margin: 0 auto;
    }

    .view-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 2rem;
    }

    .eyebrow {
        font-size: 0.75rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        color: var(--primary);
    }

    .ops-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 2rem;
    }

    .panel {
        background: var(--card);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        overflow: hidden;
    }

    .panel-header {
        padding: 1.25rem 1.5rem;
        border-bottom: 1px solid var(--border);
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .panel-header h3 {
        font-weight: 600;
    }

    .badge {
        font-size: 0.75rem;
        background: oklch(from var(--primary) l c h / 0.15);
        color: var(--primary);
        padding: 0.2rem 0.6rem;
        border-radius: 9999px;
        font-weight: 600;
    }

    .panel-body {
        padding: 1.5rem;
    }

    .topics-list, .packages-list {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .topic-card, .package-card {
        background: var(--background);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 1.25rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .package-card {
        flex-direction: column;
        align-items: flex-start;
        gap: 0.75rem;
    }

    .topic-heat {
        display: flex;
        align-items: center;
        gap: 0.4rem;
        margin-bottom: 0.3rem;
    }

    .heat-tag {
        font-size: 0.75rem;
        font-weight: 700;
        color: var(--primary);
    }

    .topic-info h4, .package-card h4 {
        font-weight: 600;
        font-size: 1rem;
    }

    .vol {
        font-size: 0.8rem;
        color: var(--muted-foreground);
    }

    .pkg-header {
        display: flex;
        justify-content: space-between;
        width: 100%;
        font-size: 0.75rem;
    }

    .pkg-cat {
        color: var(--primary);
        font-weight: 600;
    }

    .pkg-status {
        background: var(--muted);
        padding: 0.2rem 0.5rem;
        border-radius: 4px;
    }

    .pkg-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;
        font-size: 0.85rem;
        color: var(--muted-foreground);
        border-top: 1px solid var(--border);
        padding-top: 0.6rem;
    }

    .btn {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.6rem 1.2rem;
        border-radius: var(--radius);
        font-weight: 600;
        cursor: pointer;
        border: none;
    }

    .btn-primary {
        background: var(--primary);
        color: var(--background);
    }

    .btn-sm {
        padding: 0.4rem 0.8rem;
        font-size: 0.8rem;
    }

    .btn-outline {
        background: transparent;
        border: 1px solid var(--border);
        color: var(--foreground);
    }
</style>
