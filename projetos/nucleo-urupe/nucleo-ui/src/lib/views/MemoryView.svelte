<script>
  import { Button, Input, Select, Checkbox, Switch, Dialog, EmptyState, Skeleton, Badge, Alert, Spinner } from '@talos/ui';

    import { onMount } from 'svelte';
    import { memory } from '../stores';
    import { BrainCircuit, Clock, Users, MessageSquare, HelpCircle, Thermometer } from 'lucide-svelte';

    /** @type {string} */
    let selectedDate = new Date().toISOString().split('T')[0];

    async function fetchMemoryByDate(date) {
        try {
            const res = await fetch(`/api/memory/today?date=${date}`); // Backend uses date param or defaults to today
            if (res.ok) {
                const data = await res.json();
                memory.set(data);
            }
        } catch (err) {
            console.error('Failed to fetch memory capsules:', err);
        }
    }

    function handleDateChange(e) {
        selectedDate = e.target.value;
        fetchMemoryByDate(selectedDate);
    }

    onMount(() => {
        fetchMemoryByDate(selectedDate);
    });
</script>

<div class="view-content">
    <div class="header-actions">
        <div class="date-picker">
            <Clock size={16} />
            <input type="date" value={selectedDate} onchange={handleDateChange} />
        </div>
    </div>

    {#if $memory?.capsules?.length}
        <div class="capsules-timeline">
            {#each $memory.capsules as capsule}
                <div class="capsule-card" class:compacted={capsule.is_merged}>
                    <div class="capsule-aside">
                        <div class="time-stamp">
                            <span class="span-icon"><Clock size={14} /></span>
                            {capsule.time_span}
                        </div>
                        <div class="capsule-line"></div>
                    </div>
                    
                    <div class="capsule-body">
                        <div class="capsule-header">
                            <h3 class="topic">{capsule.main_topic}</h3>
                            <div class="mood-badge" data-mood={capsule.mood.toLowerCase()}>
                                <Thermometer size={14} />
                                {capsule.mood}
                            </div>
                        </div>

                        <div class="participants">
                            <Users size={14} />
                            {#each capsule.participants as p}
                                <span class="participant-tag">{p}</span>
                            {/each}
                        </div>

                        <div class="facts-section">
                            <h4>Fatos Relevantes</h4>
                            <ul>
                                {#each capsule.key_facts as fact}
                                    <li>{fact}</li>
                                {/each}
                            </ul>
                        </div>

                        {#if capsule.unresolved_questions?.length}
                        <div class="questions-section">
                            <h4><HelpCircle size={14} /> Questões em Aberto</h4>
                            <ul>
                                {#each capsule.unresolved_questions as q}
                                    <li>{q}</li>
                                {/each}
                            </ul>
                        </div>
                        {/if}
                    </div>
                </div>
            {/each}
        </div>
    {:else}
        <div class="empty-state">
            <BrainCircuit size={48} />
            <p>Nenhuma cápsula de memória encontrada para esta data.</p>
            <span class="detail">O CapsuleWorker gera novos registros a cada hora.</span>
        </div>
    {/if}
</div>

<style>
    .view-content {
        padding: 2rem;
        max-width: 1000px;
        margin: 0 auto;
    }

    .header-actions {
        display: flex;
        justify-content: flex-end;
        margin-bottom: 3rem;
    }

    .date-picker {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        background: var(--glass-surface);
        padding: 0.5rem 1rem;
        border-radius: 12px;
        border: 1px solid var(--border);
    }

    .date-picker input {
        background: transparent;
        border: none;
        color: var(--foreground);
        font-family: var(--font-mono);
        color-scheme: dark;
    }

    .capsules-timeline {
        display: flex;
        flex-direction: column;
        gap: 0;
    }

    .capsule-card {
        display: grid;
        grid-template-columns: 140px 1fr;
        gap: 2rem;
    }

    .capsule-aside {
        display: flex;
        flex-direction: column;
        align-items: flex-end;
        padding-top: 0.5rem;
    }

    .time-stamp {
        font-family: var(--font-mono);
        font-size: 0.85rem;
        color: var(--muted-foreground);
        display: flex;
        align-items: center;
        gap: 0.5rem;
        background: var(--card);
        padding: 0.25rem 0.75rem;
        border-radius: 20px;
        border: 1px solid var(--border);
    }

    .capsule-line {
        width: 2px;
        flex: 1;
        background: var(--border);
        margin-right: 2.5rem;
        margin-top: 1rem;
        margin-bottom: 1rem;
    }

    .capsule-card:last-child .capsule-line {
        display: none;
    }

    .capsule-body {
        background: var(--card);
        border: 1px solid var(--border);
        border-radius: 20px;
        padding: 1.5rem 2rem;
        margin-bottom: 2rem;
        transition: transform 0.2s ease, border-color 0.2s ease;
    }

    .capsule-body:hover {
        border-color: var(--muted);
        transform: translateX(4px);
    }

    .capsule-card.compacted .capsule-body {
        border-left: 4px solid var(--primary);
        background: linear-gradient(90deg, var(--muted), transparent);
    }

    .capsule-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 1rem;
    }

    .topic {
        margin: 0;
        font-family: var(--font-serif);
        font-size: 1.4rem;
        color: var(--foreground);
    }

    .mood-badge {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.75rem;
        font-weight: 700;
        text-transform: uppercase;
        padding: 0.35rem 0.75rem;
        border-radius: 8px;
        background: rgba(255, 255, 255, 0.05);
    }

    .participants {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        gap: 0.5rem;
        margin-bottom: 1.5rem;
        color: var(--muted-foreground);
    }

    .participant-tag {
        font-size: 0.8rem;
        background: var(--background);
        padding: 0.15rem 0.5rem;
        border-radius: 6px;
        border: 1px solid var(--border);
    }

    h4 {
        font-size: 0.8rem;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        color: var(--muted-foreground);
        margin: 1.5rem 0 0.75rem;
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    ul {
        margin: 0;
        padding: 0;
        list-style: none;
    }

    li {
        font-size: 0.95rem;
        color: var(--muted-foreground);
        padding: 0.4rem 0;
        border-bottom: 1px solid rgba(255, 255, 255, 0.03);
        line-height: 1.5;
    }

    li:last-child {
        border-bottom: none;
    }

    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 5rem 2rem;
        text-align: center;
        color: var(--muted-foreground);
    }

    .empty-state p {
        font-size: 1.2rem;
        margin: 1.5rem 0 0.5rem;
    }

    .empty-state .detail {
        font-size: 0.9rem;
        opacity: 0.7;
    }
</style>
