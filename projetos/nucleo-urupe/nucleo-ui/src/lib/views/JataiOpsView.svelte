<script>
    import { Cpu, DollarSign, Activity, Database, Server, CheckCircle2, Shield } from 'lucide-svelte';

    let fleetStats = {
        activeAgents: 14,
        totalTokens24h: '4.8M tokens',
        monthlyRevenue: 'R$ 18.400,00',
        urupeFundContribution: 'R$ 12.880,00 (70%)'
    };

    let b2bProducts = [
        { name: 'Vico ⚖️', domain: 'Copilot Jurídico', instances: 8, status: 'Operacional', mrr: 'R$ 11.200' },
        { name: 'Harness RH 💼', domain: 'Atendimento & RH', instances: 4, status: 'Operacional', mrr: 'R$ 5.600' },
        { name: 'Copilot ESG 🌿', domain: 'Gestão Ambiental', instances: 2, status: 'Operacional', mrr: 'R$ 1.600' }
    ];
</script>

<div class="jatai-view">
    <div class="view-header">
        <div>
            <p class="eyebrow">Plataforma Jataí 🐝</p>
            <h2 class="font-serif text-2xl font-bold">Frota B2B & Transparência do Autofinanciamento</h2>
        </div>
        <div class="fund-badge">
            <DollarSign size={16} />
            <span>Fundo Urupê: <b>{fleetStats.urupeFundContribution}</b></span>
        </div>
    </div>

    <!-- KPIs -->
    <div class="kpi-strip">
        <div class="kpi-card">
            <span class="label">Agentes Ativos na Frota</span>
            <span class="value">{fleetStats.activeAgents} robôs</span>
        </div>

        <div class="kpi-card">
            <span class="label">Consumo de Tokens (24h)</span>
            <span class="value">{fleetStats.totalTokens24h}</span>
        </div>

        <div class="kpi-card">
            <span class="label">Receita Recorrente B2B</span>
            <span class="value text-emerald">{fleetStats.monthlyRevenue}</span>
        </div>
    </div>

    <!-- Products Fleet -->
    <div class="panel">
        <div class="panel-header">
            <h3>Frotas de Agentes B2B em Produção</h3>
        </div>
        <div class="panel-body">
            <div class="fleet-table">
                <div class="table-header">
                    <span>Produto</span>
                    <span>Domínio</span>
                    <span>Instâncias</span>
                    <span>Status</span>
                    <span>Faturamento Mensal</span>
                </div>
                {#each b2bProducts as p}
                    <div class="table-row">
                        <span class="name"><b>{p.name}</b></span>
                        <span class="domain">{p.domain}</span>
                        <span class="instances">{p.instances} instâncias</span>
                        <span class="status"><CheckCircle2 size={14} class="text-emerald" /> {p.status}</span>
                        <span class="mrr">{p.mrr}</span>
                    </div>
                {/each}
            </div>
        </div>
    </div>
</div>

<style>
    .jatai-view {
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

    .fund-badge {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        background: oklch(from var(--primary) l c h / 0.15);
        color: var(--primary);
        padding: 0.5rem 1rem;
        border-radius: 9999px;
        font-size: 0.9rem;
    }

    .kpi-strip {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
        gap: 1.5rem;
        margin-bottom: 2rem;
    }

    .kpi-card {
        background: var(--card);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 1.5rem;
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .kpi-card .label {
        font-size: 0.8rem;
        font-weight: 600;
        color: var(--muted-foreground);
    }

    .kpi-card .value {
        font-size: 1.75rem;
        font-weight: 700;
        font-family: var(--font-serif);
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

    .fleet-table {
        display: flex;
        flex-direction: column;
    }

    .table-header {
        display: grid;
        grid-template-columns: 1.5fr 1.5fr 1fr 1fr 1.5fr;
        padding: 0.75rem 1rem;
        font-size: 0.8rem;
        font-weight: 700;
        color: var(--muted-foreground);
        border-bottom: 1px solid var(--border);
    }

    .table-row {
        display: grid;
        grid-template-columns: 1.5fr 1.5fr 1fr 1fr 1.5fr;
        padding: 1rem;
        align-items: center;
        border-bottom: 1px solid var(--border);
        font-size: 0.9rem;
    }

    .status {
        display: flex;
        align-items: center;
        gap: 0.4rem;
    }

    .mrr {
        font-weight: 700;
        color: var(--primary);
    }
</style>
