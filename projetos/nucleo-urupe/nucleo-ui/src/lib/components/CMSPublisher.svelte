<script>
    import { onMount } from 'svelte';
    import { Newspaper, Send, Trash2, Plus, Edit3, CheckCircle, RefreshCw } from 'lucide-svelte';

    let articles = $state([]);
    let loading = $state(true);
    let saving = $state(false);
    let statusMsg = $state("");

    // Form State
    let id = $state("");
    let title = $state("");
    let slug = $state("");
    let summary = $state("");
    let content = $state("");
    let author = $state("Frente Urupê");
    let category = $state("Notícia");
    let isPublished = $state(true);
    let publishToDiscord = $state(false);

    async function loadArticles() {
        loading = true;
        try {
            const res = await fetch('/api/cms/articles');
            if (res.ok) {
                articles = await res.json();
            }
        } catch (e) {
            console.error("Erro ao carregar artigos do CMS:", e);
        } finally {
            loading = false;
        }
    }

    function resetForm() {
        id = "";
        title = "";
        slug = "";
        summary = "";
        content = "";
        author = "Frente Urupê";
        category = "Notícia";
        isPublished = true;
        publishToDiscord = false;
    }

    /** @param {any} art */
    function editArticle(art) {
        id = art.id;
        title = art.title;
        slug = art.slug;
        summary = art.summary || "";
        content = art.content;
        author = art.author || "Frente Urupê";
        category = art.category || "Notícia";
        isPublished = art.is_published;
    }

    async function saveArticle() {
        if (!title || !content) {
            alert("Título e conteúdo são obrigatórios!");
            return;
        }
        saving = true;

        if (!slug) {
            slug = title.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '');
        }

        const payload = {
            id,
            title,
            slug,
            summary,
            content,
            author,
            category,
            is_published: isPublished,
            publish_to_discord: publishToDiscord
        };

        try {
            const res = await fetch('/api/cms/article/save', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            if (res.ok) {
                statusMsg = "Artigo salvo e publicado com sucesso!";
                resetForm();
                await loadArticles();
                setTimeout(() => statusMsg = "", 3000);
            } else {
                alert("Erro ao salvar artigo");
            }
        } catch (e) {
            alert("Erro de conexão");
        } finally {
            saving = false;
        }
    }

    /** @param {string} artId */
    async function deleteArticle(artId) {
        if (!confirm("Deseja apagar esta notícia do CMS?")) return;
        try {
            const res = await fetch(`/api/cms/article/delete?id=${encodeURIComponent(artId)}`, {
                method: 'DELETE'
            });
            if (res.ok) {
                statusMsg = "Artigo excluído!";
                await loadArticles();
                setTimeout(() => statusMsg = "", 3000);
            }
        } catch (e) {
            alert("Erro ao excluir artigo");
        }
    }

    onMount(() => {
        loadArticles();
    });
</script>

<div class="cms-publisher p-6 bg-[#1d2021] text-[#ebdbb2] rounded-xl border border-[#3c3836]">
    <header class="flex justify-between items-center pb-4 mb-6 border-b border-[#3c3836]">
        <div class="flex items-center gap-3">
            <div class="p-2.5 bg-[#fe8019]/10 text-[#fe8019] rounded-lg border border-[#fe8019]/30">
                <Newspaper size={22} />
            </div>
            <div>
                <h2 class="font-serif text-2xl font-bold text-[#fabd2f]">Prensa Digital • CMS Urupê News</h2>
                <p class="text-xs text-[#a89984]">Publicação oficial de artigos, manifesto e informes comunitários</p>
            </div>
        </div>

        <button 
            onclick={resetForm}
            class="flex items-center gap-2 px-3 py-2 bg-[#b8bb26] hover:bg-[#b8bb26]/90 text-[#1d2021] rounded-lg text-xs font-bold font-mono transition-colors"
        >
            <Plus size={14} />
            Novo Artigo
        </button>
    </header>

    {#if statusMsg}
        <div class="mb-4 p-3 bg-[#b8bb26]/15 border border-[#b8bb26]/40 text-[#b8bb26] text-xs font-mono rounded-lg flex items-center gap-2">
            <CheckCircle size={14} />
            {statusMsg}
        </div>
    {/if}

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Editor Form (2 cols) -->
        <div class="lg:col-span-2 p-5 bg-[#282828] border border-[#3c3836] rounded-xl flex flex-col gap-4">
            <h3 class="font-serif text-lg font-bold text-[#fe8019] border-b border-[#3c3836] pb-2">
                {id ? "Editar Artigo" : "Criar Nova Publicação"}
            </h3>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                    <label class="block text-xs font-mono text-[#a89984] mb-1">Título da Notícia</label>
                    <input 
                        type="text" 
                        bind:value={title} 
                        placeholder="Ex: Soberania Digital no Território..." 
                        class="w-full px-3 py-2 bg-[#1d2021] border border-[#3c3836] rounded-lg text-sm text-[#ebdbb2] focus:border-[#b8bb26] outline-none"
                    />
                </div>
                <div>
                    <label class="block text-xs font-mono text-[#a89984] mb-1">Slug URL</label>
                    <input 
                        type="text" 
                        bind:value={slug} 
                        placeholder="soberania-digital-no-territorio" 
                        class="w-full px-3 py-2 bg-[#1d2021] border border-[#3c3836] rounded-lg text-sm text-[#8ec07c] font-mono focus:border-[#b8bb26] outline-none"
                    />
                </div>
            </div>

            <div>
                <label class="block text-xs font-mono text-[#a89984] mb-1">Resumo Executivo</label>
                <input 
                    type="text" 
                    bind:value={summary} 
                    placeholder="Síntese de até 2 linhas da notícia..." 
                    class="w-full px-3 py-2 bg-[#1d2021] border border-[#3c3836] rounded-lg text-sm text-[#ebdbb2] focus:border-[#b8bb26] outline-none"
                />
            </div>

            <div>
                <label class="block text-xs font-mono text-[#a89984] mb-1">Conteúdo (Markdown)</label>
                <textarea 
                    bind:value={content} 
                    rows="8"
                    placeholder="Escreva a notícia em formato Markdown..."
                    class="w-full px-3 py-2 bg-[#1d2021] border border-[#3c3836] rounded-lg text-sm text-[#ebdbb2] font-mono focus:border-[#b8bb26] outline-none"
                ></textarea>
            </div>

            <div class="flex items-center justify-between pt-2 border-t border-[#3c3836]">
                <label class="flex items-center gap-2 text-xs font-mono text-[#b8bb26] cursor-pointer">
                    <input type="checkbox" bind:checked={publishToDiscord} class="accent-[#b8bb26]" />
                    <span>Transmitir para o Discord (#📰│noticias)</span>
                </label>

                <button 
                    onclick={saveArticle}
                    disabled={saving}
                    class="flex items-center gap-2 px-4 py-2 bg-[#fe8019] hover:bg-[#fe8019]/90 text-[#1d2021] rounded-lg text-xs font-bold font-mono transition-colors"
                >
                    <Send size={14} />
                    {saving ? "Salvando..." : "Publicar Notícia"}
                </button>
            </div>
        </div>

        <!-- Articles List (1 col) -->
        <div class="p-5 bg-[#282828] border border-[#3c3836] rounded-xl flex flex-col gap-3">
            <h3 class="font-serif text-lg font-bold text-[#fabd2f] border-b border-[#3c3836] pb-2 flex justify-between items-center">
                <span>Acervo CMS</span>
                <span class="text-xs font-mono text-[#a89984]">{articles.length} posts</span>
            </h3>

            <div class="flex flex-col gap-2 max-h-[420px] overflow-y-auto pr-1">
                {#each articles as art}
                    <div class="article-item p-3 bg-[#1d2021] border border-[#3c3836] rounded-lg flex justify-between items-start hover:border-[#b8bb26]/50 transition-colors">
                        <div>
                            <h4 class="font-serif text-sm font-bold text-[#ebdbb2] line-clamp-1">{art.title}</h4>
                            <p class="text-[10px] font-mono text-[#a89984] mt-0.5">/{art.slug}</p>
                        </div>

                        <div class="flex items-center gap-1">
                            <button 
                                onclick={() => editArticle(art)}
                                title="Editar"
                                class="p-1 text-[#a89984] hover:text-[#b8bb26] transition-colors"
                            >
                                <Edit3 size={13} />
                            </button>
                            <button 
                                onclick={() => deleteArticle(art.id)}
                                title="Excluir"
                                class="p-1 text-[#a89984] hover:text-[#fe8019] transition-colors"
                            >
                                <Trash2 size={13} />
                            </button>
                        </div>
                    </div>
                {:else}
                    <div class="text-center py-8 text-[#a89984] font-mono text-xs">
                        Nenhum artigo cadastrado.
                    </div>
                {/each}
            </div>
        </div>
    </div>
</div>
