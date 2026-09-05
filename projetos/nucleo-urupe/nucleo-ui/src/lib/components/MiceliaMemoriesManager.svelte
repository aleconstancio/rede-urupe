<script>
    import { onMount } from 'svelte';
    import { Shield, Trash2, RefreshCw, Search, Sparkles, Brain, CheckCircle } from 'lucide-svelte';

    let memories = $state([]);
    let loading = $state(true);
    let searchQuery = $state("");
    let statusMsg = $state("");

    async function loadMemories() {
        loading = true;
        try {
            const res = await fetch('/api/micelia/memories');
            if (res.ok) {
                const data = await res.json();
                memories = data.annotations || [];
            }
        } catch (e) {
            console.error("Erro ao carregar memórias da Micélia:", e);
        } finally {
            loading = false;
        }
    }

    /** @param {string} messageId */
    async function deleteMemory(messageId) {
        if (!confirm("Deseja apagar esta cápsula de memória da Micélia 🍄?")) return;
        try {
            const res = await fetch(`/api/micelia/memory/delete?message_id=${encodeURIComponent(messageId)}`, {
                method: 'DELETE'
            });
            if (res.ok) {
                statusMsg = "Cápsula de memória limpa com sucesso!";
                memories = memories.filter(m => m.message_id !== messageId);
                setTimeout(() => statusMsg = "", 3000);
            }
        } catch (e) {
            alert("Erro ao excluir memória");
        }
    }

    onMount(() => {
        loadMemories();
    });

    let filteredMemories = $derived(memories.filter(m => {
        const text = [
            m.author_id,
            m.message_id,
            ...(m.topic_tags || []),
            ...(m.stance_tags || []),
            m.evidence_type
        ].join(' ').toLowerCase();
        return text.includes(searchQuery.toLowerCase());
    }));
</script>

<div class="micelia-memories-manager p-6 bg-[#1d2021] text-[#ebdbb2] rounded-xl border border-[#3c3836]">
    <header class="flex justify-between items-center pb-4 mb-6 border-b border-[#3c3836]">
        <div class="flex items-center gap-3">
            <div class="p-2.5 bg-[#b8bb26]/10 text-[#b8bb26] rounded-lg border border-[#b8bb26]/30">
                <Brain size={22} />
            </div>
            <div>
                <h2 class="font-serif text-2xl font-bold text-[#fabd2f]">Memória Episódica da Micélia 🍄</h2>
                <p class="text-xs text-[#a89984]">Curadoria de fatos aprendidos e anotações nos canais do Discord</p>
            </div>
        </div>

        <div class="flex items-center gap-3">
            <div class="relative">
                <Search size={16} class="absolute left-3 top-3 text-[#a89984]" />
                <input 
                    type="text" 
                    placeholder="Filtrar memórias..." 
                    bind:value={searchQuery} 
                    class="pl-9 pr-4 py-2 bg-[#282828] border border-[#3c3836] rounded-lg text-sm text-[#ebdbb2] focus:border-[#b8bb26] outline-none"
                />
            </div>
            <button 
                onclick={loadMemories}
                class="flex items-center gap-2 px-3 py-2 bg-[#282828] hover:bg-[#3c3836] border border-[#3c3836] rounded-lg text-xs font-mono text-[#b8bb26] transition-colors"
            >
                <RefreshCw size={14} class={loading ? "animate-spin" : ""} />
                Atualizar
            </button>
        </div>
    </header>

    {#if statusMsg}
        <div class="mb-4 p-3 bg-[#b8bb26]/15 border border-[#b8bb26]/40 text-[#b8bb26] text-xs font-mono rounded-lg flex items-center gap-2">
            <CheckCircle size={14} />
            {statusMsg}
        </div>
    {/if}

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        {#each filteredMemories as m}
            <div class="memory-card p-4 bg-[#282828] border border-[#3c3836] rounded-xl flex flex-col justify-between hover:border-[#b8bb26]/50 transition-all">
                <div>
                    <div class="flex justify-between items-start mb-2">
                        <span class="text-xs font-mono text-[#fe8019] bg-[#fe8019]/10 px-2 py-0.5 rounded border border-[#fe8019]/30">
                            Autor: @{m.author_id.slice(-8)}
                        </span>
                        <button 
                            onclick={() => deleteMemory(m.message_id)}
                            title="Apagar esta memória"
                            class="p-1 text-[#a89984] hover:text-[#fe8019] hover:bg-[#fe8019]/10 rounded transition-colors"
                        >
                            <Trash2 size={14} />
                        </button>
                    </div>

                    <p class="text-xs text-[#d5c4a1] font-mono mb-3">
                        MSG: ID <span class="text-[#8ec07c]">{m.message_id}</span> • Canal <span class="text-[#fabd2f]">{m.channel_id}</span>
                    </p>

                    {#if m.topic_tags && m.topic_tags.length > 0}
                        <div class="flex flex-wrap gap-1.5 mb-2">
                            {#each m.topic_tags as tag}
                                <span class="text-[10px] bg-[#b8bb26]/15 text-[#b8bb26] px-2 py-0.5 rounded border border-[#b8bb26]/30">
                                    #{tag}
                                </span>
                            {/each}
                        </div>
                    {/if}
                </div>

                <div class="pt-3 border-t border-[#3c3836]/60 flex justify-between items-center text-[10px] font-mono text-[#a89984]">
                    <span>Durabilidade: <strong class="text-[#fabd2f]">{(m.durability_score * 100).toFixed(0)}%</strong></span>
                    <span>Tipo: {m.evidence_type || 'Geral'}</span>
                </div>
            </div>
        {:else}
            <div class="col-span-2 text-center py-12 text-[#a89984] font-mono text-sm">
                Nenhuma cápsula de memória encontrada.
            </div>
        {/each}
    </div>
</div>
