<script>
    import { Compass, Eye, ShieldAlert, Satellite, MapPin, Layers, Trees } from 'lucide-svelte';

    let satelliteFeeds = [
        { collection: 'CBERS-4A', sensor: 'WPM (8m)', target: 'Bacia Amazônica', status: 'Ativo' },
        { collection: 'Sentinel-2', sensor: 'MSI (10m)', target: 'Mata Atlântica / SP', status: 'Sincronizado' },
        { collection: 'Amazonia-1', sensor: 'WFC (64m)', target: 'Cerrado Central', status: 'Ativo' }
    ];

    let alerts = [
        { type: 'Alerta de Queimada', location: 'Região Norte SP', level: 'Alto (NBR Anomalia)', date: 'Hoje' },
        { type: 'Estoque de Carbono', location: 'Assentamento Rural #4', level: 'Estável (+2.4% NDVI)', date: 'Ontem' }
    ];
</script>

<div class="guara-view">
    <div class="view-header">
        <div>
            <p class="eyebrow">Módulo Guará 🪶</p>
            <h2 class="font-serif text-2xl font-bold">Inteligência Geoespacial & Sensoriamento Ambiental</h2>
        </div>
        <button class="btn btn-primary">
            <Satellite size={18} />
            <span>Consultar Satélites INPE/Sentinel</span>
        </button>
    </div>

    <div class="guara-grid">
        <!-- Satellites Stream -->
        <div class="panel">
            <div class="panel-header">
                <h3>Coleções de Satélites & Constelações</h3>
            </div>
            <div class="panel-body">
                <div class="sats-list">
                    {#each satelliteFeeds as s}
                        <div class="sat-card">
                            <div class="sat-icon">
                                <Satellite size={24} class="text-emerald" />
                            </div>
                            <div class="sat-details">
                                <h4>{s.collection} ({s.sensor})</h4>
                                <p>Alvo: <b>{s.target}</b></p>
                            </div>
                            <span class="status-badge">{s.status}</span>
                        </div>
                    {/each}
                </div>
            </div>
        </div>

        <!-- Environmental Alerts -->
        <div class="panel">
            <div class="panel-header">
                <h3>Alertas Territoriais & Análise Espectral</h3>
            </div>
            <div class="panel-body">
                <div class="alerts-list">
                    {#each alerts as a}
                        <div class="alert-card">
                            <div class="alert-header">
                                <ShieldAlert size={18} class="text-amber" />
                                <span class="alert-title">{a.type}</span>
                                <span class="date">{a.date}</span>
                            </div>
                            <p class="loc"><MapPin size={14} /> {a.location}</p>
                            <span class="lvl">{a.level}</span>
                        </div>
                    {/each}
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    .guara-view {
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

    .guara-grid {
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
    }

    .panel-header h3 {
        font-weight: 600;
    }

    .panel-body {
        padding: 1.5rem;
    }

    .sats-list, .alerts-list {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .sat-card {
        background: var(--background);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 1.25rem;
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    .sat-details h4 {
        font-size: 0.95rem;
        font-weight: 600;
    }

    .sat-details p {
        font-size: 0.8rem;
        color: var(--muted-foreground);
    }

    .status-badge {
        margin-left: auto;
        font-size: 0.75rem;
        background: oklch(from var(--primary) l c h / 0.15);
        color: var(--primary);
        padding: 0.2rem 0.6rem;
        border-radius: 9999px;
        font-weight: 600;
    }

    .alert-card {
        background: var(--background);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 1.25rem;
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .alert-header {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.9rem;
        font-weight: 600;
    }

    .date {
        margin-left: auto;
        font-size: 0.75rem;
        color: var(--muted-foreground);
    }

    .loc {
        display: flex;
        align-items: center;
        gap: 0.3rem;
        font-size: 0.85rem;
        color: var(--muted-foreground);
    }

    .lvl {
        font-size: 0.8rem;
        font-weight: 600;
        color: var(--primary);
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
</style>
