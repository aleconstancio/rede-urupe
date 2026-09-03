<script>
  import { Button, Input, Select, Checkbox, Switch, Dialog, EmptyState, Skeleton, Badge, Alert, Spinner, Tabs, TabsList, TabsTrigger, TabsContent } from '@talos/ui';

    import { onMount } from 'svelte';
    import { Activity, TrendingUp, Users, MessageSquare, Coins } from 'lucide-svelte';
    import Chart from 'chart.js/auto';

    let activeTab = $state('overview');
    let sentimentData = [];
    let tokenData = { by_model: [], by_reason: [], trend: [] };
    let growthData = [];
    let channelHealth = null;
    let engagementScore = 0;
    let loading = false;

    let sentimentChart = null;
    let tokenChart = null;
    let growthChart = null;

    const tabs = [
        { id: 'overview', label: 'Visao Geral', icon: Activity },
        { id: 'tokens', label: 'Tokens', icon: Coins },
        { id: 'growth', label: 'Crescimento', icon: TrendingUp }
    ];

    onMount(() => {
        loadAll();
    });

    async function loadAll() {
        loading = true;
        await Promise.all([
            loadSentiment(),
            loadTokens(),
            loadGrowth(),
            loadChannelHealth(),
            loadEngagement()
        ]);
        loading = false;
        setTimeout(() => {
            renderCharts();
        }, 100);
    }

    async function loadSentiment() {
        try {
            const res = await fetch('/api/analytics/sentiment?days=30');
            if (res.ok) {
                const data = await res.json();
                sentimentData = data.sentiment || [];
            }
        } catch (err) {
            console.error('Failed to load sentiment:', err);
        }
    }

    async function loadTokens() {
        try {
            const res = await fetch('/api/analytics/tokens?days=30');
            if (res.ok) {
                tokenData = await res.json();
            }
        } catch (err) {
            console.error('Failed to load tokens:', err);
        }
    }

    async function loadGrowth() {
        try {
            const res = await fetch('/api/analytics/growth?days=30');
            if (res.ok) {
                const data = await res.json();
                growthData = data.growth || [];
            }
        } catch (err) {
            console.error('Failed to load growth:', err);
        }
    }

    async function loadChannelHealth() {
        try {
            const res = await fetch('/api/analytics/channels');
            if (res.ok) {
                const data = await res.json();
                channelHealth = data.channel;
            }
        } catch (err) {
            console.error('Failed to load channel health:', err);
        }
    }

    async function loadEngagement() {
        try {
            const res = await fetch('/api/analytics/engagement?days=7');
            if (res.ok) {
                const data = await res.json();
                engagementScore = data.score || 0;
            }
        } catch (err) {
            console.error('Failed to load engagement:', err);
        }
    }

    function renderCharts() {
        renderSentimentChart();
        renderTokenChart();
        renderGrowthChart();
    }

    function renderSentimentChart() {
        const canvas = document.getElementById('sentiment-chart');
        if (!canvas || sentimentData.length === 0) return;

        if (sentimentChart) sentimentChart.destroy();

        sentimentChart = new Chart(canvas, {
            type: 'line',
            data: {
                labels: sentimentData.map(d => d.date),
                datasets: [
                    {
                        label: 'Intensidade',
                        data: sentimentData.map(d => d.avg_intensity),
                        borderColor: '#f59e0b',
                        backgroundColor: 'rgba(245, 158, 11, 0.1)',
                        fill: true,
                        tension: 0.4
                    },
                    {
                        label: 'Conflito',
                        data: sentimentData.map(d => d.avg_conflict),
                        borderColor: '#ef4444',
                        backgroundColor: 'rgba(239, 68, 68, 0.1)',
                        fill: true,
                        tension: 0.4
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: { legend: { labels: { color: '#94a3b8' } } },
                scales: {
                    x: { ticks: { color: '#64748b' }, grid: { color: 'rgba(100,116,139,0.1)' } },
                    y: { ticks: { color: '#64748b' }, grid: { color: 'rgba(100,116,139,0.1)' }, min: 0, max: 1 }
                }
            }
        });
    }

    function renderTokenChart() {
        const canvas = document.getElementById('token-chart');
        if (!canvas || tokenData.by_model.length === 0) return;

        if (tokenChart) tokenChart.destroy();

        tokenChart = new Chart(canvas, {
            type: 'doughnut',
            data: {
                labels: tokenData.by_model.map(d => d.model || 'unknown'),
                datasets: [{
                    data: tokenData.by_model.map(d => d.total_tokens),
                    backgroundColor: ['#f59e0b', '#3b82f6', '#10b981', '#8b5cf6', '#ef4444']
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: { legend: { labels: { color: '#94a3b8' } } }
            }
        });
    }

    function renderGrowthChart() {
        const canvas = document.getElementById('growth-chart');
        if (!canvas || growthData.length === 0) return;

        if (growthChart) growthChart.destroy();

        growthChart = new Chart(canvas, {
            type: 'bar',
            data: {
                labels: growthData.map(d => d.date),
                datasets: [
                    {
                        label: 'Entradas',
                        data: growthData.map(d => d.joins),
                        backgroundColor: '#10b981'
                    },
                    {
                        label: 'Saidas',
                        data: growthData.map(d => -d.leaves),
                        backgroundColor: '#ef4444'
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: { legend: { labels: { color: '#94a3b8' } } },
                scales: {
                    x: { stacked: true, ticks: { color: '#64748b' }, grid: { color: 'rgba(100,116,139,0.1)' } },
                    y: { stacked: true, ticks: { color: '#64748b' }, grid: { color: 'rgba(100,116,139,0.1)' } }
                }
            }
        });
    }

    function formatNumber(n) {
        if (!n) return '0';
        return n.toLocaleString('pt-BR');
    }
</script>

<div class="analytics-view">
    <div class="analytics-header">
        <Activity size={24} />
        <h2>Analitica</h2>
    </div>

    <Tabs bind:value={activeTab}>
    <TabsList>
      <TabsTrigger value="overview"><Activity size={16} /> Visao Geral</TabsTrigger>
      <TabsTrigger value="tokens"><Coins size={16} /> Tokens</TabsTrigger>
      <TabsTrigger value="growth"><TrendingUp size={16} /> Crescimento</TabsTrigger>
    </TabsList>

    <TabsContent value="overview">
      {#if channelHealth}
        <div class="overview-grid">
          <div class="stat-card"><MessageSquare size={20} /><div class="stat-value">{formatNumber(channelHealth.total_messages)}</div><div class="stat-label">Total Mensagens</div></div>
          <div class="stat-card"><Users size={20} /><div class="stat-value">{formatNumber(channelHealth.unique_authors)}</div><div class="stat-label">Autores Unicos</div></div>
          <div class="stat-card"><Activity size={20} /><div class="stat-value">{engagementScore.toFixed(1)}</div><div class="stat-label">Score Engajamento</div></div>
          <div class="stat-card"><Coins size={20} /><div class="stat-value">{formatNumber(tokenData.trend?.reduce((a, b) => a + b.total_tokens, 0) || 0)}</div><div class="stat-label">Tokens (30d)</div></div>
        </div>
      {/if}
      <div class="chart-container"><h3>Tendencia de Sentimento</h3><div class="chart-wrapper"><canvas id="sentiment-chart"></canvas></div></div>
    </TabsContent>

    <TabsContent value="tokens">
      <div class="charts-row">
        <div class="chart-container half"><h3>Uso por Modelo</h3><div class="chart-wrapper doughnut"><canvas id="token-chart"></canvas></div></div>
        <div class="chart-container half">
          <h3>Tendencia Diaria</h3>
          <div class="token-list">
            {#each tokenData.trend || [] as t}
              <div class="token-item"><span class="token-date">{t.date}</span><span class="token-value">{formatNumber(t.total_tokens)} tokens</span><span class="token-calls">{t.call_count} chamadas</span></div>
            {/each}
            {#if (tokenData.trend || []).length === 0}<div class="empty-state">Nenhum dado de tokens.</div>{/if}
          </div>
        </div>
      </div>
      <div class="chart-container"><h3>Uso por Motivo</h3><div class="token-list">{#each tokenData.by_reason || [] as t}<div class="token-item"><span class="token-reason">{t.reason}</span><span class="token-value">{formatNumber(t.total_tokens)} tokens</span><span class="token-calls">{t.call_count} chamadas</span></div>{/each}</div></div>
    </TabsContent>

    <TabsContent value="growth">
      <div class="chart-container"><h3>Crescimento de Membros</h3><div class="chart-wrapper"><canvas id="growth-chart"></canvas></div></div>
      <div class="growth-list">{#each growthData.slice(-10).reverse() as g}<div class="growth-item"><span class="growth-date">{g.date}</span><span class="growth-joins">+{g.joins}</span><span class="growth-leaves">-{g.leaves}</span><span class="growth-net" class:positive={g.net_change > 0} class:negative={g.net_change < 0}>{g.net_change > 0 ? '+' : ''}{g.net_change}</span></div>{/each}</div>
    </TabsContent>
  </Tabs>
</div>

<style>
    .analytics-view {
        padding: 1.5rem;
        max-width: 1200px;
        margin: 0 auto;
    }

    .analytics-header {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        margin-bottom: 1.5rem;
        color: var(--foreground);
    }

    .analytics-header h2 {
        margin: 0;
        font-size: 1.5rem;
    }
    
    .overview-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 1rem;
        margin-bottom: 1.5rem;
    }

    .stat-card {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.5rem;
        padding: 1.5rem;
        background: var(--glass-surface);
        border: 1px solid var(--border);
        border-radius: 12px;
        color: var(--primary);
    }

    .stat-value {
        font-size: 2rem;
        font-weight: 700;
        color: var(--foreground);
    }

    .stat-label {
        font-size: 0.85rem;
        color: var(--muted-foreground);
    }

    .chart-container {
        background: var(--glass-surface);
        border: 1px solid var(--border);
        border-radius: 12px;
        padding: 1.5rem;
        margin-bottom: 1.5rem;
    }

    .chart-container.half {
        flex: 1;
    }

    .chart-container h3 {
        margin: 0 0 1rem 0;
        font-size: 1rem;
        color: var(--foreground);
    }

    .chart-wrapper {
        height: 300px;
        position: relative;
    }

    .chart-wrapper.doughnut {
        height: 250px;
        max-width: 250px;
        margin: 0 auto;
    }

    .charts-row {
        display: flex;
        gap: 1.5rem;
        margin-bottom: 1.5rem;
    }

    .token-list {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        max-height: 300px;
        overflow-y: auto;
    }

    .token-item {
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 0.5rem 0.75rem;
        background: var(--background);
        border-radius: 6px;
        font-size: 0.85rem;
    }

    .token-date, .token-reason {
        color: var(--muted-foreground);
        min-width: 100px;
    }

    .token-value {
        color: var(--primary);
        font-weight: 600;
    }

    .token-calls {
        color: var(--muted-foreground);
        font-size: 0.75rem;
    }

    .growth-list {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .growth-item {
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 0.75rem 1rem;
        background: var(--glass-surface);
        border: 1px solid var(--border);
        border-radius: 8px;
    }

    .growth-date {
        color: var(--muted-foreground);
        min-width: 100px;
    }

    .growth-joins {
        color: #10b981;
        font-weight: 600;
    }

    .growth-leaves {
        color: #ef4444;
        font-weight: 600;
    }

    .growth-net {
        margin-left: auto;
        font-weight: 700;
    }

    .growth-net.positive { color: #10b981; }
    .growth-net.negative { color: #ef4444; }

    .empty-state {
        text-align: center;
        padding: 2rem;
        color: var(--muted-foreground);
    }
</style>
