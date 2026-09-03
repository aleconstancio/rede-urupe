<script>
    import { metrics } from '../stores';

    /** @type {Record<string, string>} */
    const labelMap = {
        active_users: 'Participantes (1h)',
        messages_1h: 'Mensagens (1h)',
        channel_name: 'Canal Ativo',
        active_persona: 'Persona Ativa',
        bot_messages: 'Pulso de Resposta',
        memories_today: 'Arquivos Criados'
    };

    /** @type {Record<string, string>} */
    const detailMap = {
        active_users: 'Pessoas observadas na última hora.',
        messages_1h: 'Volume de mensagens no canal (60m).',
        channel_name: 'Identificador do canal de operações.',
        active_persona: 'Identidade core em controle agora.',
        bot_messages: 'Total de respostas enviadas hoje.',
        memories_today: 'Novas memórias consolidadas hoje.'
    };

    /** @type {Record<string, string>} */
    const toneMap = {
        active_users: 'accent',
        messages_1h: 'success',
        channel_name: 'dim',
        active_persona: 'warning',
        bot_messages: 'success',
        memories_today: 'accent'
    };

    /** @type {Record<string, number>} */
    const orderMap = {
        active_users: 0,
        messages_1h: 1,
        channel_name: 2,
        active_persona: 3,
        bot_messages: 4,
        memories_today: 5
    };

    /** @param {string} value */
    function toTitleCase(value) {
        return value
            .replace(/[_-]+/g, ' ')
            .replace(/\b\w/g, (letter) => letter.toUpperCase());
    }

    /** 
     * @param {string} key 
     * @param {any} value 
     */
    function formatMetricValue(key, value) {
        const numeric = Number(value);
        if (!Number.isFinite(numeric)) {
            return String(value);
        }

        if (key === 'avg_latency') {
            return `${Math.round(numeric)}ms`;
        }

        if (key.includes('usage')) {
            return `${Math.round(numeric * 100)}%`;
        }

        return numeric.toLocaleString('pt-BR');
    }

    /** 
     * @param {string} key 
     * @param {any} value 
     */
    function getProgress(key, value) {
        const numeric = Number(value);
        if (!Number.isFinite(numeric)) {
            return null;
        }

        if (key === 'total_tokens_24h') {
            return Math.max(0, Math.min(100, (numeric / 1000000) * 100)); // 1M token budget
        }

        if (key === 'active_users') {
            return Math.max(0, Math.min(100, numeric * 10));
        }

        if (key === 'api_calls') {
            return Math.max(0, Math.min(100, (numeric / 500) * 100));
        }

        if (key === 'avg_latency') {
            return Math.max(0, Math.min(100, (numeric / 2500) * 100));
        }

        return Math.max(0, Math.min(100, numeric));
    }

    /** @param {any} metric */
    function normalizeLegacyMetric(metric) {
        const key = metric.Key || metric.key || metric.Label || metric.label || 'metric';
        const value = metric.Value ?? metric.value ?? 0;
        const tone = (metric.Tone || metric.tone || /** @type {Record<string, string>} */(toneMap)[key]) || 'accent';

        return {
            key,
            label: metric.Label || metric.label || (/** @type {Record<string, string>} */(labelMap))[key] || toTitleCase(key),
            value: metric.Value ?? metric.value ?? formatMetricValue(key, value),
            detail: metric.Detail || metric.detail || (/** @type {Record<string, string>} */(detailMap))[key] || 'Métrica operacional recebida do backend.',
            tone,
            featured: Boolean(metric.Featured || metric.featured),
            mono: Boolean(metric.Mono || metric.mono),
            progress: metric.Progress ?? metric.progress ?? getProgress(key, value)
        };
    }

    /** @param {any} payload */
    function normalizeMetrics(payload) {
        if (!payload) {
            return [];
        }

        if (Array.isArray(payload)) {
            return payload.map(normalizeLegacyMetric);
        }

        if (typeof payload === 'object') {
            return Object.entries(payload)
                .sort(([left], [right]) => ((/** @type {Record<string, number>} */(orderMap))[left] ?? 99) - ((/** @type {Record<string, number>} */(orderMap))[right] ?? 99))
                .map(([key, value]) => ({
                    key,
                    label: (/** @type {Record<string, string>} */(labelMap))[key] || toTitleCase(key),
                    value: formatMetricValue(key, value),
                    detail: (/** @type {Record<string, string>} */(detailMap))[key] || 'Métrica operacional recebida do backend.',
                    tone: (/** @type {Record<string, string>} */(toneMap))[key] || 'accent',
                    featured: false,
                    mono: true,
                    progress: getProgress(key, value)
                }));
        }

        return [];
    }

    $: cards = normalizeMetrics($metrics);
</script>

<div class="metrics-grid">
    {#if cards.length}
        {#each cards as metric}
            <div
                class="metric-card"
                class:featured={metric.featured}
                class:tone-success={metric.tone === 'success'}
                class:tone-warning={metric.tone === 'warning'}
                class:tone-accent={metric.tone === 'accent'}
            >
                <div class="card-head">
                    <p class="metric-label">{metric.label}</p>
                    <span class="metric-key">{metric.key}</span>
                </div>

                <div class="metric-value" class:mono={metric.mono}>{metric.value}</div>
                <p class="metric-detail">{metric.detail}</p>

                {#if metric.progress !== null}
                    <div class="metric-rail" aria-hidden="true">
                        <span style={`width: ${metric.progress}%`}></span>
                    </div>
                {/if}
            </div>
        {/each}
    {:else}
        <div class="loading-state">
            <p>Sincronizando pulso operacional...</p>
        </div>
    {/if}
</div>

<style>
    .metrics-grid {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 1.25rem;
    }

    .metric-card {
        position: relative;
        overflow: hidden;
        min-height: 10.5rem;
        padding: 1.35rem;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        background:
            linear-gradient(180deg, rgba(255, 255, 255, 0.035), transparent 34%),
            var(--card);
        border: 1px solid var(--border);
        border-top: 1px solid var(--border);
        border-radius: 18px;
        box-shadow: inset 0 1px 1px var(--border);
        transition: transform 0.25s ease, border-color 0.25s ease, box-shadow 0.25s ease;
    }

    .metric-card:hover {
        transform: translateY(-3px);
        border-color: var(--border);
        box-shadow: 0 16px 30px rgba(0, 0, 0, 0.18), inset 0 1px 1px var(--border);
    }

    .metric-card.featured {
        background:
            linear-gradient(135deg, rgba(255, 255, 255, 0.05), transparent 45%),
            var(--card);
        border-color: var(--muted);
    }

    .metric-card.tone-success {
        border-left: 3px solid var(--success);
    }

    .metric-card.tone-warning {
        border-left: 3px solid var(--warning);
    }

    .metric-card.tone-accent {
        border-left: 3px solid var(--primary);
    }

    .card-head {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
        gap: 0.75rem;
        margin-bottom: 1rem;
    }

    .metric-label {
        margin: 0;
        text-transform: uppercase;
        letter-spacing: 0.16em;
        font-size: 0.65rem;
        font-weight: 700;
        color: var(--muted-foreground);
    }

    .metric-key {
        font-family: var(--font-mono);
        font-size: 0.65rem;
        color: var(--muted-foreground);
        background: rgba(255, 255, 255, 0.04);
        border: 1px solid var(--border);
        border-radius: 999px;
        padding: 0.2rem 0.45rem;
        white-space: nowrap;
    }

    .metric-value {
        font-family: var(--font-serif);
        font-size: 2.35rem;
        line-height: 1;
        font-weight: 700;
        letter-spacing: -0.03em;
        color: var(--foreground);
        margin-bottom: 0.55rem;
    }

    .metric-detail {
        margin: 0 0 1rem;
        font-size: 0.83rem;
        line-height: 1.45;
        color: var(--muted-foreground);
    }

    .metric-rail {
        width: 100%;
        height: 6px;
        border-radius: 999px;
        overflow: hidden;
        background: rgba(255, 255, 255, 0.05);
        border: 1px solid var(--border);
    }

    .metric-rail span {
        display: block;
        height: 100%;
        border-radius: inherit;
        background: linear-gradient(90deg, var(--muted), var(--primary));
    }

    .loading-state {
        grid-column: 1 / -1;
        padding: 2.5rem;
        text-align: center;
        background: var(--glass-surface);
        border-radius: var(--radius);
        border: 1px dashed var(--border);
        color: var(--muted-foreground);
        font-style: italic;
    }
</style>
