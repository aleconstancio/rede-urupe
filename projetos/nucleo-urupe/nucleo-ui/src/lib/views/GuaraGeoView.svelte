<script>
    import { onMount } from 'svelte';
    import { Compass, Satellite, Flame, Droplets, MapPin, Layers, AlertTriangle, ShieldCheck, Activity, RefreshCw } from 'lucide-svelte';

    let activeLayer = $state('ndvi'); // 'ndvi' | 'fires' | 'water' | 'territories'

    const territories = [
        { name: 'Assentamento Terra Livre (PR)', area: '1.240 ha', ndvi: 0.78, status: 'Saudável', alerts: 0 },
        { name: 'Comunidade Rio Vermelho (BA)', area: '850 ha', ndvi: 0.62, status: 'Atenção Hídrica', alerts: 1 },
        { name: 'Horta Comunitária Zona Sul (SP)', area: '12 ha', ndvi: 0.81, status: 'Excelente', alerts: 0 },
        { name: 'Reserva Manejada Tapajós (PA)', area: '4.500 ha', ndvi: 0.85, status: 'Monitorando Focos', alerts: 2 }
    ];

    const satelliteFeeds = [
        { satellite: 'Sentinel-2A', passTime: 'Hoje, 09:30', status: 'Processado 100%' },
        { satellite: 'Landsat 9', passTime: 'Hoje, 06:15', status: 'Indexado FTS5' },
        { satellite: 'CBERS-4A', passTime: 'Ontem, 14:00', status: 'Nuvem Limpa' }
    ];
</script>

<div class="guara-geo-view">
    <div class="view-header">
        <div>
            <p class="eyebrow text-cyan">Inteligência Geoespacial & Satélites</p>
            <h2 class="font-serif text-2xl font-bold flex items-center gap-2">
                <span>Guará Geo 🪶</span>
                <span class="badge-tag">Monitoramento Territorial</span>
            </h2>
        </div>

        <div class="layer-selector">
            <button class="layer-btn {activeLayer === 'ndvi' ? 'active' : ''}" onclick={() => activeLayer = 'ndvi'}>
                <Layers size={16} /> Índice NDVI
            </button>
            <button class="layer-btn {activeLayer === 'fires' ? 'active' : ''}" onclick={() => activeLayer = 'fires'}>
                <Flame size={16} /> Focos de Calor
            </button>
            <button class="layer-btn {activeLayer === 'water' ? 'active' : ''}" onclick={() => activeLayer = 'water'}>
                <Droplets size={16} /> Estresse Hídrico
            </button>
        </div>
    </div>

    <div class="geo-grid">
        <!-- Main Map Simulation Canvas -->
        <div class="panel map-panel">
            <div class="panel-header">
                <span class="map-title">Visão de Satélite em Tempo Real ({activeLayer.toUpperCase()})</span>
                <span class="live-badge"><Activity size={12} /> Feed Sentinel-2</span>
            </div>

            <div class="map-canvas">
                <div class="map-overlay">
                    <div class="geo-pin pin-1" title="Terra Livre">
                        <MapPin size={24} class="text-cyan" />
                        <span class="pin-label">Terra Livre</span>
                    </div>
                    <div class="geo-pin pin-2" title="Rio Vermelho">
                        <MapPin size={24} class="text-amber" />
                        <span class="pin-label">Rio Vermelho</span>
                    </div>
                    <div class="geo-pin pin-3" title="Tapajós">
                        <MapPin size={24} class="text-emerald" />
                        <span class="pin-label">Tapajós</span>
                    </div>
                </div>

                <div class="map-legend">
                    <span class="leg-title">Índice de Vigor Vegetal (NDVI)</span>
                    <div class="legend-bar">
                        <span>0.0 (Solo Exposto)</span>
                        <div class="gradient-strip"></div>
                        <span>1.0 (Massa Densa)</span>
                    </div>
                </div>
            </div>
        </div>

        <!-- Right Side: Territories Status & Feeds -->
        <div class="side-col">
            <div class="panel">
                <div class="panel-header">
                    <h3>Territórios Monitorados ({territories.length})</h3>
                </div>
                <div class="territories-list">
                    {#each territories as ter}
                        <div class="ter-card">
                            <div class="ter-header">
                                <span class="ter-name">{ter.name}</span>
                                <span class="ter-status {ter.alerts > 0 ? 'alert' : 'ok'}">{ter.status}</span>
                            </div>
                            <div class="ter-stats">
                                <span>Área: <b>{ter.area}</b></span>
                                <span>NDVI: <b>{ter.ndvi}</b></span>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>

            <div class="panel">
                <div class="panel-header">
                    <h3>Passagens de Satélites</h3>
                </div>
                <div class="feeds-list">
                    {#each satelliteFeeds as sat}
                        <div class="sat-item">
                            <Satellite size={18} class="text-cyan" />
                            <div class="sat-info">
                                <span class="sat-name">{sat.satellite}</span>
                                <span class="sat-time">{sat.passTime}</span>
                            </div>
                            <span class="sat-status">{sat.status}</span>
                        </div>
                    {/each}
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    .guara-geo-view { padding: 2rem; max-width: 1600px; margin: 0 auto; }
    .view-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
    .eyebrow { font-size: 0.75rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; }
    .text-cyan { color: #06b6d4; }
    .badge-tag { font-size: 0.75rem; background: rgba(6, 182, 212, 0.15); color: #06b6d4; padding: 0.2rem 0.6rem; border-radius: 9999px; }

    .layer-selector { display: flex; gap: 0.5rem; background: var(--card); border: 1px solid var(--border); padding: 0.3rem; border-radius: var(--radius); }
    .layer-btn { display: flex; align-items: center; gap: 0.4rem; padding: 0.5rem 1rem; border-radius: calc(var(--radius) - 2px); font-size: 0.85rem; font-weight: 600; background: transparent; border: none; color: var(--muted-foreground); cursor: pointer; }
    .layer-btn.active { background: #06b6d4; color: var(--background); }

    .geo-grid { display: grid; grid-template-columns: 1fr 380px; gap: 2rem; }
    .panel { background: var(--card); border: 1px solid var(--border); border-radius: var(--radius); overflow: hidden; }
    .panel-header { padding: 1.25rem 1.5rem; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center; }
    .panel-header h3 { font-weight: 600; font-size: 1rem; }
    .map-title { font-weight: 700; font-size: 0.9rem; }

    .live-badge { font-size: 0.75rem; color: #06b6d4; font-weight: 700; display: flex; align-items: center; gap: 0.4rem; }

    /* Map Canvas Simulation */
    .map-canvas {
        height: 520px;
        background: radial-gradient(circle at center, #0f172a 0%, #020617 100%);
        position: relative;
        overflow: hidden;
    }

    .map-overlay { position: absolute; inset: 0; }
    .geo-pin { position: absolute; display: flex; flex-direction: column; align-items: center; cursor: pointer; transition: transform 0.2s; }
    .geo-pin:hover { transform: scale(1.15); }
    .pin-label { font-size: 0.75rem; font-weight: 700; background: rgba(0, 0, 0, 0.75); color: white; padding: 0.2rem 0.5rem; border-radius: 4px; border: 1px solid var(--border); margin-top: 0.2rem; }

    .pin-1 { top: 30%; left: 40%; }
    .pin-2 { top: 55%; left: 65%; }
    .pin-3 { top: 25%; left: 25%; }

    .map-legend {
        position: absolute;
        bottom: 1.5rem;
        left: 1.5rem;
        background: rgba(15, 23, 42, 0.85);
        backdrop-filter: blur(10px);
        border: 1px solid var(--border);
        padding: 1rem;
        border-radius: var(--radius);
        width: 320px;
    }

    .leg-title { font-size: 0.75rem; font-weight: 700; color: #06b6d4; display: block; margin-bottom: 0.5rem; }
    .legend-bar { font-size: 0.7rem; color: var(--muted-foreground); display: flex; flex-direction: column; gap: 0.3rem; }
    .gradient-strip { height: 8px; border-radius: 4px; background: linear-gradient(to right, #ef4444, #f59e0b, #10b981, #06b6d4); }

    .side-col { display: flex; flex-direction: column; gap: 2rem; }

    .territories-list { padding: 1.25rem; display: flex; flex-direction: column; gap: 1rem; }
    .ter-card { background: var(--background); border: 1px solid var(--border); border-radius: var(--radius); padding: 1rem; }
    .ter-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem; }
    .ter-name { font-weight: 700; font-size: 0.9rem; }
    .ter-status { font-size: 0.75rem; padding: 0.15rem 0.5rem; border-radius: 4px; font-weight: 600; }
    .ter-status.ok { background: rgba(16, 185, 129, 0.15); color: #10b981; }
    .ter-status.alert { background: rgba(245, 158, 11, 0.15); color: #f59e0b; }
    .ter-stats { display: flex; justify-content: space-between; font-size: 0.8rem; color: var(--muted-foreground); }

    .feeds-list { padding: 1.25rem; display: flex; flex-direction: column; gap: 0.75rem; }
    .sat-item { display: flex; align-items: center; gap: 0.75rem; background: var(--background); border: 1px solid var(--border); border-radius: var(--radius); padding: 0.75rem; }
    .sat-info { display: flex; flex-direction: column; flex: 1; }
    .sat-name { font-weight: 700; font-size: 0.85rem; }
    .sat-time { font-size: 0.75rem; color: var(--muted-foreground); }
    .sat-status { font-size: 0.75rem; color: #06b6d4; font-weight: 600; }

    @media (max-width: 1000px) { .geo-grid { grid-template-columns: 1fr; } }
</style>
