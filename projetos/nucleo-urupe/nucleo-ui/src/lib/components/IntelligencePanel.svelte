<script>
    import { onMount, onDestroy } from 'svelte';
    import Chart from 'chart.js/auto';
    import { mode } from 'mode-watcher';
    import { Activity, Shapes, Users, Calendar, Hash } from 'lucide-svelte';

    let activityCanvas;
    let categoriesCanvas;
    let authorsCanvas;

    let activityChart;
    let categoriesChart;
    let authorsChart;

    let stats = null;
    let refreshTimer;
    let isBackfilling = false;
    let backfillDays = 7;

    let chartTokens = {
        accent: '#3b82f6',
        accentSoft: 'rgba(59, 130, 246, 0.14)',
        textDim: '#a1a1aa',
        borderSoft: 'rgba(255, 255, 255, 0.05)'
    };

    const dayOrder = [1, 2, 3, 4, 5, 6, 0];
    const dayLabels = ['Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sab', 'Dom'];

    function readChartTokens() {
        if (typeof document === 'undefined') {
            return;
        }

        const styles = getComputedStyle(document.documentElement);
        chartTokens = {
            accent: styles.getPropertyValue('--accent').trim() || '#3b82f6',
            accentSoft: styles.getPropertyValue('--accent-soft').trim() || 'rgba(59, 130, 246, 0.14)',
            textDim: styles.getPropertyValue('--text-dim').trim() || '#a1a1aa',
            borderSoft: styles.getPropertyValue('--border-soft').trim() || 'rgba(255, 255, 255, 0.05)'
        };
    }

    function toPoints(series) {
        return (series || []).map((item) => ({
            label: item.key ?? item.label ?? '',
            value: Number(item.value ?? item.count ?? 0) || 0
        }));
    }

    function total(points) {
        return points.reduce((sum, point) => sum + point.value, 0);
    }

    function formatHourLabel(label) {
        if (!label) {
            return '--';
        }

        const match = String(label).match(/^(\d{2})/);
        return match ? `${match[1]}h` : String(label);
    }

    function topPoint(points) {
        return points.reduce((best, point) => (point.value > (best?.value ?? -1) ? point : best), null);
    }

    function getHeatmapMax() {
        if (!stats?.heatmap?.length) {
            return 1;
        }

        return Math.max(...stats.heatmap.map((point) => point.value), 1);
    }

    function getHeatValue(day, hour) {
        return stats?.heatmap?.find((point) => point.day === day && point.hour === hour)?.value ?? 0;
    }

    function cellStyle(value) {
        if (!value) {
            return 'background: rgba(255,255,255,0.02)';
        }

        const max = getHeatmapMax();
        const intensity = Math.log1p(value) / Math.log1p(max);
        return `background: color-mix(in srgb, var(--primary) ${Math.min(intensity, 1) * 100}%, transparent)`;
    }

    function updateLineChart(chart, canvas, points) {
        if (!canvas) {
            return;
        }

        const labels = points.map((point) => point.label);
        const values = points.map((point) => point.value);

        if (chart) {
            chart.data.labels = labels;
            chart.data.datasets[0].data = values;
            chart.data.datasets[0].borderColor = chartTokens.accent;
            chart.data.datasets[0].pointBackgroundColor = chartTokens.accent;
            chart.update();
            return chart;
        }

        chart = new Chart(canvas.getContext('2d'), {
            type: 'line',
            data: {
                labels,
                datasets: [{
                    data: values,
                    borderColor: chartTokens.accent,
                    backgroundColor: (context) => {
                        const { chart, ctx, chartArea } = context;
                        if (!chartArea) {
                            return null;
                        }

                        const gradient = ctx.createLinearGradient(0, chartArea.top, 0, chartArea.bottom);
                        gradient.addColorStop(0, chartTokens.accentSoft);
                        gradient.addColorStop(1, 'rgba(0,0,0,0)');
                        return gradient;
                    },
                    borderWidth: 2,
                    fill: true,
                    tension: 0.36,
                    pointRadius: 0,
                    pointHoverRadius: 5,
                    pointBackgroundColor: chartTokens.accent
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: { display: false },
                    tooltip: {
                        intersect: false,
                        mode: 'index'
                    }
                },
                scales: {
                    x: {
                        grid: { display: false },
                        ticks: { color: chartTokens.textDim, font: { size: 10 } }
                    },
                    y: {
                        beginAtZero: true,
                        grid: { color: chartTokens.borderSoft },
                        ticks: { color: chartTokens.textDim, font: { size: 10 }, precision: 0 }
                    }
                }
            }
        });

        return chart;
    }

    function updateBarChart(chart, canvas, points) {
        if (!canvas) {
            return;
        }

        const labels = points.map((point) => point.label);
        const values = points.map((point) => point.value);

        if (chart) {
            chart.data.labels = labels;
            chart.data.datasets[0].data = values;
            chart.data.datasets[0].backgroundColor = chartTokens.accent;
            chart.update();
            return chart;
        }

        chart = new Chart(canvas.getContext('2d'), {
            type: 'bar',
            data: {
                labels,
                datasets: [{
                    data: values,
                    backgroundColor: chartTokens.accent,
                    borderRadius: 5,
                    barThickness: 12,
                    maxBarThickness: 14
                }]
            },
            options: {
                indexAxis: 'y',
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: { display: false },
                    tooltip: { intersect: false, mode: 'nearest' }
                },
                scales: {
                    x: {
                        beginAtZero: true,
                        grid: { display: false },
                        ticks: { display: false }
                    },
                    y: {
                        grid: { display: false },
                        ticks: {
                            color: chartTokens.textDim,
                            font: { size: 10 },
                            padding: 6
                        }
                    }
                }
            }
        });

        return chart;
    }

    async function fetchStats() {
        try {
            const res = await fetch('/api/stats');
            if (res.ok) {
                stats = await res.json();
            }
        } catch (error) {
            console.error(error);
        }
    }

    async function triggerBackfill() {
        if (!confirm(`Deseja iniciar a retroalimentação de ${backfillDays} dias? Isso pode levar algum tempo.`)) {
            return;
        }

        isBackfilling = true;
        try {
            const res = await fetch('/api/sync/backfill', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ days: backfillDays })
            });

            if (res.ok) {
                alert('Sincronização iniciada em segundo plano.');
            } else {
                alert('Erro ao iniciar sincronização.');
            }
        } catch (error) {
            console.error(error);
            alert('Falha na comunicação com o servidor.');
        } finally {
            isBackfilling = false;
        }
    }

    $: hourlyPoints = toPoints(stats?.hourly ?? []);
    $: categoryPoints = toPoints((stats?.categories ?? []).filter((item) => item.key !== 'Não Categorizado' && item.key !== 'Nao Categorizado')).slice(0, 8);
    $: authorPoints = toPoints(stats?.authors ?? []).slice(0, 8).map((item) => ({
        ...item,
        label: `@${item.label}`
    }));
    $: totalMessages = total(hourlyPoints);
    $: topCategory = topPoint(categoryPoints);
    $: topAuthor = topPoint(authorPoints);
    $: peakHour = topPoint(hourlyPoints);
    $: overviewCards = [
        {
            label: 'Mensagens 24h',
            value: totalMessages.toLocaleString('en-US'),
            detail: hourlyPoints.length ? `${hourlyPoints.length} janelas ativas` : 'Sem tráfego recente',
            tone: 'accent'
        },
        {
            label: 'Tema líder',
            value: topCategory ? topCategory.label : 'Sem categoria',
            detail: topCategory ? `${topCategory.value} mensagens` : 'Sem dados suficientes',
            tone: 'success'
        },
        {
            label: 'Voz líder',
            value: topAuthor ? topAuthor.label : 'Sem autor',
            detail: topAuthor ? `${topAuthor.value} mensagens` : 'Sem dados suficientes',
            tone: 'warning'
        },
        {
            label: 'Pico horário',
            value: peakHour ? formatHourLabel(peakHour.label) : '--',
            detail: peakHour ? `${peakHour.value} mensagens` : 'Sem dados suficientes',
            tone: 'accent'
        }
    ];

    function refreshCharts() {
        readChartTokens();
        activityChart = updateLineChart(activityChart, activityCanvas, hourlyPoints);
        categoriesChart = updateBarChart(categoriesChart, categoriesCanvas, categoryPoints);
        authorsChart = updateBarChart(authorsChart, authorsCanvas, authorPoints);
    }

    $: if (stats && typeof document !== 'undefined') {
        $mode;
        refreshCharts();
    }

    onMount(() => {
        readChartTokens();
        fetchStats();
        refreshTimer = setInterval(fetchStats, 10000);
    });

    onDestroy(() => {
        if (refreshTimer) {
            clearInterval(refreshTimer);
        }

        if (activityChart) {
            activityChart.destroy();
        }

        if (categoriesChart) {
            categoriesChart.destroy();
        }

        if (authorsChart) {
            authorsChart.destroy();
        }
    });
</script>

<div class="analytics-layout">
    <section class="panel hero-panel">
        <div class="section-header">
            <div>
                <div class="title-with-icon">
                    <Activity size={14} class="icon-accent" />
                    <p class="eyebrow">Registro de atividades</p>
                </div>
                <h3>Fluxo de Interações Normalizado</h3>
            </div>

            <div class="header-actions">
                <div class="status-pill">
                    <span class="dot"></span>
                    Atualiza a cada 10s
                </div>

                <div class="backfill-controls">
                    <select bind:value={backfillDays} class="backfill-select">
                        <option value={1}>1 dia</option>
                        <option value={7}>7 dias</option>
                        <option value={30}>30 dias</option>
                    </select>
                    <button class="backfill-btn" onclick={triggerBackfill} disabled={isBackfilling}>
                        {isBackfilling ? '⏳' : 'Retroalimentar'}
                    </button>
                </div>
            </div>
        </div>

        <div class="overview-grid">
            {#each overviewCards as card}
                <article class="overview-card" class:tone-success={card.tone === 'success'} class:tone-warning={card.tone === 'warning'} class:tone-accent={card.tone === 'accent'}>
                    <p class="card-label">{card.label}</p>
                    <div class="card-value">{card.value}</div>
                    <p class="card-detail">{card.detail}</p>
                </article>
            {/each}
        </div>

        <div class="chart-panel chart-panel-wide">
            <div class="chart-header">
                <div>
                    <div class="title-with-icon">
                        <Shapes size={14} class="icon-accent" />
                        <p class="eyebrow">Linha do tempo</p>
                    </div>
                    <h3>Volume de Mensagens</h3>
                </div>
                <p class="chart-caption">Mensagens por janela horária nas últimas 24 horas.</p>
            </div>
            <div class="chart-canvas canvas-large">
                <canvas bind:this={activityCanvas}></canvas>
            </div>
        </div>
    </section>

    <div class="split-row">
        <section class="panel chart-panel">
            <div class="chart-header">
                <div>
                    <div class="title-with-icon">
                        <Hash size={14} class="icon-accent" />
                        <p class="eyebrow">Categorias</p>
                    </div>
                    <h3>Top Temas</h3>
                </div>
                <p class="chart-caption">Categorias mais recorrentes no canal.</p>
            </div>
            <div class="chart-canvas canvas-compact">
                <canvas bind:this={categoriesCanvas}></canvas>
            </div>
        </section>

        <section class="panel chart-panel">
            <div class="chart-header">
                <div>
                    <div class="title-with-icon">
                        <Users size={14} class="icon-accent" />
                        <p class="eyebrow">Participantes</p>
                    </div>
                    <h3>Principais Vozes</h3>
                </div>
                <p class="chart-caption">Autores mais ativos no recorte atual.</p>
            </div>
            <div class="chart-canvas canvas-compact">
                <canvas bind:this={authorsCanvas}></canvas>
            </div>
        </section>
    </div>

    <section class="panel heatmap-section">
        <div class="chart-header heatmap-header">
            <div>
                <div class="title-with-icon">
                    <Calendar size={14} class="icon-accent" />
                    <p class="eyebrow">Mapa de calor</p>
                </div>
                <h3>Ritmo Semanal do Servidor</h3>
            </div>
            <p class="chart-caption">Densidade relativa de mensagens por dia e hora.</p>
        </div>

        <div class="heatmap-wrapper">
            <div class="day-labels">
                {#each dayLabels as dayLabel}
                    <div class="day-label">{dayLabel}</div>
                {/each}
            </div>

            <div class="heatmap-container">
                <div class="heatmap-grid">
                    {#if stats && stats.heatmap}
                        {#each dayOrder as day}
                            {#each Array.from({ length: 24 }) as _, hour}
                                {@const val = getHeatValue(day, hour)}
                                <div
                                    class="heatmap-cell"
                                    style={cellStyle(val)}
                                    title={`${dayLabels[day === 0 ? 6 : day - 1]}, ${hour}h: ${val} msg`}
                                ></div>
                            {/each}
                        {/each}
                    {/if}
                </div>

                <div class="hour-labels">
                    <span>0h</span>
                    <span>6h</span>
                    <span>12h</span>
                    <span>18h</span>
                    <span>23h</span>
                </div>
            </div>
        </div>
    </section>
</div>

<style>
    .analytics-layout {
        display: flex;
        flex-direction: column;
        gap: 1.25rem;
        padding: 1rem;
    }

    .hero-panel,
    .chart-panel,
    .heatmap-section {
        padding: 1.4rem;
    }

    .section-header,
    .chart-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        gap: 1rem;
    }

    .title-with-icon {
        display: flex;
        align-items: center;
        gap: 0.6rem;
        margin-bottom: 0.25rem;
    }

    :global(.icon-accent) {
        color: var(--primary);
    }

    .header-actions {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        flex-wrap: wrap;
        justify-content: flex-end;
    }

    .status-pill {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.45rem 0.75rem;
        border-radius: 999px;
        background: var(--muted);
        color: var(--primary);
        border: 1px solid var(--border);
        font-size: 0.8rem;
        font-weight: 600;
    }

    .dot {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: var(--primary);
        box-shadow: 0 0 0 4px rgba(255, 255, 255, 0.04);
    }

    .backfill-controls {
        display: flex;
        gap: 0.5rem;
        align-items: center;
    }

    .backfill-select {
        background: var(--card);
        color: var(--muted-foreground);
        border: 1px solid var(--border);
        border-radius: 8px;
        font-size: 0.75rem;
        padding: 0.35rem 0.5rem;
        outline: none;
    }

    .backfill-btn {
        background: var(--muted);
        color: var(--primary);
        border: 1px solid var(--primary);
        border-radius: 8px;
        padding: 0.35rem 0.75rem;
        font-size: 0.75rem;
        font-weight: 700;
        cursor: pointer;
        transition: all 0.2s ease;
    }

    .backfill-btn:hover:not(:disabled) {
        background: var(--primary);
        color: white;
    }

    .backfill-btn:disabled {
        opacity: 0.55;
        cursor: not-allowed;
    }

    .overview-grid {
        display: grid;
        grid-template-columns: repeat(4, minmax(0, 1fr));
        gap: 1rem;
    }

    .overview-card {
        min-height: 9.25rem;
        padding: 1.1rem;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        border: 1px solid var(--border);
        border-radius: 16px;
        background:
            linear-gradient(180deg, rgba(255, 255, 255, 0.03), transparent 42%),
            rgba(255, 255, 255, 0.015);
    }

    .overview-card.tone-success {
        border-left: 3px solid var(--success);
    }

    .overview-card.tone-warning {
        border-left: 3px solid var(--warning);
    }

    .overview-card.tone-accent {
        border-left: 3px solid var(--primary);
    }

    .card-label {
        margin: 0;
        text-transform: uppercase;
        letter-spacing: 0.16em;
        font-size: 0.64rem;
        font-weight: 700;
        color: var(--muted-foreground);
    }

    .card-value {
        font-family: var(--font-serif);
        font-size: 1.85rem;
        font-weight: 700;
        letter-spacing: -0.03em;
        line-height: 1.1;
        color: var(--foreground);
        margin: 0.7rem 0 0.35rem;
    }

    .card-detail,
    .chart-caption {
        margin: 0;
        font-size: 0.82rem;
        line-height: 1.45;
        color: var(--muted-foreground);
    }

    .chart-panel-wide {
        margin-top: 0.25rem;
    }

    .chart-canvas {
        margin-top: 1rem;
    }

    .canvas-large {
        height: 220px;
    }

    .canvas-compact {
        height: 235px;
    }

    .split-row {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1.25rem;
    }

    .heatmap-wrapper {
        display: flex;
        gap: 1rem;
        margin-top: 1.25rem;
    }

    .day-labels {
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        padding: 4px 0;
        height: 120px;
    }

    .day-label {
        font-size: 0.7rem;
        color: var(--muted-foreground);
        font-family: var(--font-mono);
        text-transform: uppercase;
    }

    .heatmap-container {
        flex: 1;
    }

    .heatmap-grid {
        display: grid;
        grid-template-columns: repeat(24, minmax(0, 1fr));
        grid-template-rows: repeat(7, minmax(0, 1fr));
        gap: 4px;
        height: 120px;
    }

    .heatmap-cell {
        border-radius: 2px;
        transition: transform 0.15s ease, box-shadow 0.15s ease;
    }

    .heatmap-cell:hover {
        transform: scale(1.4);
        z-index: 10;
        box-shadow: 0 0 10px var(--primary);
    }

    .hour-labels {
        display: flex;
        justify-content: space-between;
        margin-top: 0.5rem;
        font-size: 0.65rem;
        color: var(--muted-foreground);
        font-family: var(--font-mono);
    }

    @media (max-width: 1024px) {
        .overview-grid,
        .split-row {
            grid-template-columns: 1fr;
        }

        .section-header,
        .chart-header {
            flex-direction: column;
        }

        .header-actions {
            justify-content: flex-start;
        }
    }
</style>
