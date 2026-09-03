<script>
    import { onMount } from 'svelte';
    import { Cpu, DollarSign, Activity, TrendingUp, ShieldCheck, Zap, Server, BarChart3 } from 'lucide-svelte';

    const b2bAgents = [
        { name: 'Vico ⚖️ (Assistente Jurídico B2B)', client: 'Escritórios Aliados', status: 'Operacional', tokensToday: '142k', revenue: 'R$ 8.500/mês' },
        { name: 'Minos RH 👔 (Triagem Corporativa)', client: 'Empresas Parceiras', status: 'Operacional', tokensToday: '98k', revenue: 'R$ 6.200/mês' },
        { name: 'Jataí Support 🤖 (Atendimento B2B)', client: 'Cooperativas', status: 'Manutenção', tokensToday: '12k', revenue: 'R$ 3.800/mês' }
    ];

    const financialMetrics = {
        totalRevenueMonth: 'R$ 18.500,00',
        infraCostMonth: 'R$ 3.200,00',
        netFundReinvested: 'R$ 15.300,00',
        activeFleetCount: 3
    };
</script>

<div class="jatai-ops-view">
    <div class="view-header">
        <div>
            <p class="eyebrow text-purple">Plataforma Industrial & Autofinanciamento</p>
            <h2 class="font-serif text-2xl font-bold flex items-center gap-2">
                <span>Jataí Ops 🐝</span>
                <span class="badge-tag">Frotas B2B & Fundo Urupê</span>
            </h2>
        </div>

        <div class="revenue-badge">
            <DollarSign size={16} />
            <span>Faturamento Mensal Reinvestido: <b>{financialMetrics.totalRevenueMonth}</b></span>
        </div>
    </div>

    <!-- Top Financial Overview Grid -->
    <div class="fin-grid">
        <div class="fin-card">
            <div class="fin-icon text-purple"><DollarSign size={24} /></div>
            <div class="fin-info">
                <span class="fin-label">Receita Bruta B2B</span>
                <span class="fin-val">{financialMetrics.totalRevenueMonth}</span>
            </div>
        </div>
        <div class="fin-card">
            <div class="fin-icon text-amber"><Server size={24} /></div>
            <div class="fin-info">
                <span class="fin-label">Custo de Servidores & Tokens</span>
                <span class="fin-val">{financialMetrics.infraCostMonth}</span>
            </div>
        </div>
        <div class="fin-card">
            <div class="fin-icon text-emerald"><ShieldCheck size={24} /></div>
            <div class="fin-info">
                <span class="fin-label">Fundo Urupê (Soberania)</span>
                <span class="fin-val">{financialMetrics.netFundReinvested}</span>
            </div>
        </div>
    </div>

    <!-- Fleet Status Grid -->
    <div class="fleet-section">
        <div class="panel">
            <div class="panel-header">
                <h3>Frota de Agentes Industriais em Operação</h3>
                <span class="active-count">{financialMetrics.activeFleetCount} Agentes Ativos</span>
            </div>

            <div class="fleet-list">
                {#each b2bAgents as ag}
                    <div class="agent-card">
                        <div class="ag-main">
                            <Cpu size={24} class="text-purple" />
                            <div>
                                <h4 class="ag-name">{ag.name}</h4>
                                <span class="ag-client">Cliente: {ag.client}</span>
                            </div>
                        </div>
                        <div class="ag-stats">
                            <div class="stat-box">
                                <span class="lbl">Consumo de Tokens</span>
                                <span class="val"><Zap size={12} /> {ag.tokensToday}</span>
                            </div>
                            <div class="stat-box">
                                <span class="lbl">Receita Recorrente</span>
                                <span class="val text-emerald">{ag.revenue}</span>
                            </div>
                            <span class="ag-status {ag.status === 'Operacional' ? 'ok' : 'maint'}">{ag.status}</span>
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    </div>
</div>

<style>
    .jatai-ops-view { padding: 2rem; max-width: 1600px; margin: 0 auto; }
    .view-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
    .eyebrow { font-size: 0.75rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; }
    .text-purple { color: #8b5cf6; }
    .badge-tag { font-size: 0.75rem; background: rgba(139, 92, 246, 0.15); color: #8b5cf6; padding: 0.2rem 0.6rem; border-radius: 9999px; }

    .revenue-badge { display: flex; align-items: center; gap: 0.4rem; background: rgba(16, 185, 129, 0.12); border: 1px solid rgba(16, 185, 129, 0.3); color: #10b981; padding: 0.5rem 1rem; border-radius: 9999px; font-size: 0.85rem; }

    .fin-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 1.5rem; margin-bottom: 2rem; }
    .fin-card { background: var(--card); border: 1px solid var(--border); border-radius: var(--radius); padding: 1.5rem; display: flex; align-items: center; gap: 1.25rem; }
    .fin-icon { width: 48px; height: 48px; border-radius: var(--radius); background: var(--background); display: flex; align-items: center; justify-content: center; }
    .fin-info { display: flex; flex-direction: column; }
    .fin-label { font-size: 0.8rem; color: var(--muted-foreground); font-weight: 600; }
    .fin-val { font-family: var(--font-serif); font-size: 1.5rem; font-weight: 700; }

    .panel { background: var(--card); border: 1px solid var(--border); border-radius: var(--radius); overflow: hidden; }
    .panel-header { padding: 1.25rem 1.5rem; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center; }
    .panel-header h3 { font-weight: 600; font-size: 1rem; }
    .active-count { font-size: 0.8rem; font-weight: 700; color: #8b5cf6; }

    .fleet-list { padding: 1.5rem; display: flex; flex-direction: column; gap: 1rem; }
    .agent-card { background: var(--background); border: 1px solid var(--border); border-radius: var(--radius); padding: 1.25rem; display: flex; justify-content: space-between; align-items: center; }
    .ag-main { display: flex; align-items: center; gap: 1rem; }
    .ag-name { font-family: var(--font-serif); font-size: 1.15rem; font-weight: 700; }
    .ag-client { font-size: 0.8rem; color: var(--muted-foreground); }

    .ag-stats { display: flex; align-items: center; gap: 2rem; }
    .stat-box { display: flex; flex-direction: column; }
    .lbl { font-size: 0.75rem; color: var(--muted-foreground); }
    .val { font-size: 0.9rem; font-weight: 700; display: flex; align-items: center; gap: 0.3rem; }
    .text-emerald { color: #10b981; }

    .ag-status { font-size: 0.75rem; padding: 0.2rem 0.6rem; border-radius: 4px; font-weight: 700; }
    .ag-status.ok { background: rgba(16, 185, 129, 0.15); color: #10b981; }
    .ag-status.maint { background: rgba(245, 158, 11, 0.15); color: #f59e0b; }

    @media (max-width: 1000px) { .agent-card { flex-direction: column; align-items: flex-start; gap: 1rem; } .ag-stats { width: 100%; justify-content: space-between; } }
</style>
