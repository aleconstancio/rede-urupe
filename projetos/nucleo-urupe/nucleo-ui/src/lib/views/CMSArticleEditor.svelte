<script>
    import { onMount } from 'svelte';
    import { FileText, Plus, Trash2, Edit3, CheckCircle2, ShieldCheck, Globe, Save } from 'lucide-svelte';

    let articles = $state([]);
    let loading = $state(true);
    let saving = $state(false);

    // Form state
    let editingId = $state(null);
    let title = $state('');
    let slug = $state('');
    let summary = $state('');
    let content = $state('');
    let category = $state('Notícia');
    let author = $state('Frente Urupê');
    let isPublished = $state(true);

    async function loadArticles() {
        loading = true;
        try {
            const res = await fetch('/api/cms/articles');
            if (res.ok) {
                articles = await res.json();
            }
        } catch (e) {
            console.error('Error loading CMS articles:', e);
        } finally {
            loading = false;
        }
    }

    function startNew() {
        editingId = null;
        title = '';
        slug = '';
        summary = '';
        content = '';
        category = 'Notícia';
        author = 'Frente Urupê';
        isPublished = true;
    }

    function editArticle(item) {
        editingId = item.id;
        title = item.title;
        slug = item.slug;
        summary = item.summary;
        content = item.content;
        category = item.category;
        author = item.author;
        isPublished = item.is_published;
    }

    async function handleSave() {
        if (!title.trim() || !content.trim()) {
            alert('Título e conteúdo são obrigatórios.');
            return;
        }

        saving = true;
        try {
            const payload = {
                id: editingId || '',
                slug,
                title,
                summary,
                content,
                author,
                category,
                is_published: isPublished,
                p2p_signed: true
            };

            const res = await fetch('/api/cms/article/save', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            if (res.ok) {
                await loadArticles();
                startNew();
            } else {
                alert('Erro ao salvar artigo.');
            }
        } catch (e) {
            console.error('Error saving article:', e);
        } finally {
            saving = false;
        }
    }

    async function handleDelete(id) {
        if (!confirm('Deseja realmente remover este artigo do site público?')) return;

        try {
            const res = await fetch(`/api/cms/article/delete?id=${id}`, { method: 'DELETE' });
            if (res.ok) {
                await loadArticles();
                if (editingId === id) startNew();
            }
        } catch (e) {
            console.error('Error deleting article:', e);
        }
    }

    onMount(() => {
        loadArticles();
    });
</script>

<div class="cms-editor-view">
    <div class="view-header">
        <div>
            <p class="eyebrow">Gestão de Conteúdo Institucional</p>
            <h2 class="font-serif text-2xl font-bold">CMS & Publicações Públicas</h2>
        </div>
        <button class="btn btn-primary" onclick={startNew}>
            <Plus size={18} />
            <span>Novo Artigo</span>
        </button>
    </div>

    <div class="editor-grid">
        <!-- Form Column -->
        <div class="panel form-panel">
            <div class="panel-header">
                <h3>{editingId ? 'Editar Artigo' : 'Criar Novo Artigo'}</h3>
                <span class="p2p-badge">
                    <ShieldCheck size={14} />
                    Injeção P2P Ativa
                </span>
            </div>

            <div class="form-body">
                <div class="form-group">
                    <label for="cms-title">Título do Artigo</label>
                    <input id="cms-title" type="text" bind:value={title} placeholder="Ex: Manifesto pela Soberania Digital" class="input" />
                </div>

                <div class="form-row">
                    <div class="form-group flex-1">
                        <label for="cms-category">Categoria</label>
                        <select id="cms-category" bind:value={category} class="input">
                            <option value="Manifesto">Manifesto</option>
                            <option value="Notícia">Notícia</option>
                            <option value="Formação">Formação Política</option>
                            <option value="Editorial">Editorial</option>
                        </select>
                    </div>

                    <div class="form-group flex-1">
                        <label for="cms-author">Autor</label>
                        <input id="cms-author" type="text" bind:value={author} class="input" />
                    </div>
                </div>

                <div class="form-group">
                    <label for="cms-summary">Resumo / Subtítulo</label>
                    <textarea id="cms-summary" bind:value={summary} rows="2" placeholder="Resumo curto para o card da landing page..." class="input"></textarea>
                </div>

                <div class="form-group">
                    <label for="cms-content">Conteúdo Principal (Markdown)</label>
                    <textarea id="cms-content" bind:value={content} rows="10" placeholder="Escreva o artigo em Markdown..." class="input code-font"></textarea>
                </div>

                <div class="form-checkbox">
                    <label class="checkbox-label">
                        <input type="checkbox" bind:checked={isPublished} />
                        <span>Publicar imediatamente no Site Oficial (HTTPS) e inovar no Rizoma P2P</span>
                    </label>
                </div>

                <div class="form-actions">
                    <button class="btn btn-primary" onclick={handleSave} disabled={saving}>
                        <Save size={18} />
                        <span>{saving ? 'Salvando...' : 'Salvar & Publicar'}</span>
                    </button>
                    {#if editingId}
                        <button class="btn btn-outline" onclick={startNew}>Cancelar</button>
                    {/if}
                </div>
            </div>
        </div>

        <!-- Articles List Column -->
        <div class="panel list-panel">
            <div class="panel-header">
                <h3>Artigos Registrados ({articles.length})</h3>
            </div>

            {#if loading}
                <div class="state-msg">Carregando lista...</div>
            {:else if articles.length === 0}
                <div class="state-msg">Nenhum artigo cadastrado no banco SQLite.</div>
            {:else}
                <div class="articles-list">
                    {#each articles as item}
                        <div class="list-item {editingId === item.id ? 'active' : ''}">
                            <div class="item-info">
                                <div class="item-meta">
                                    <span class="cat">{item.category}</span>
                                    <span class="status {item.is_published ? 'pub' : 'draft'}">
                                        {item.is_published ? 'Publicado' : 'Rascunho'}
                                    </span>
                                </div>
                                <h4 class="item-title">{item.title}</h4>
                                <p class="item-summary">{item.summary || 'Sem resumo.'}</p>
                            </div>

                            <div class="item-actions">
                                <button class="icon-btn" onclick={() => editArticle(item)} title="Editar">
                                    <Edit3 size={16} />
                                </button>
                                <button class="icon-btn danger" onclick={() => handleDelete(item.id)} title="Excluir">
                                    <Trash2 size={16} />
                                </button>
                            </div>
                        </div>
                    {/each}
                </div>
            {/if}
        </div>
    </div>
</div>

<style>
    .cms-editor-view {
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

    .editor-grid {
        display: grid;
        grid-template-columns: 1fr 450px;
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
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .panel-header h3 {
        font-weight: 600;
    }

    .p2p-badge {
        display: flex;
        align-items: center;
        gap: 0.4rem;
        font-size: 0.75rem;
        background: oklch(from var(--primary) l c h / 0.15);
        color: var(--primary);
        padding: 0.25rem 0.6rem;
        border-radius: 9999px;
        font-weight: 600;
    }

    .form-body {
        padding: 1.5rem;
        display: flex;
        flex-direction: column;
        gap: 1.25rem;
    }

    .form-group {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .form-row {
        display: flex;
        gap: 1rem;
    }

    .flex-1 { flex: 1; }

    label {
        font-size: 0.85rem;
        font-weight: 600;
        color: var(--muted-foreground);
    }

    .input {
        background: var(--background);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 0.6rem 0.8rem;
        color: var(--foreground);
        font-family: inherit;
        font-size: 0.9rem;
    }

    .code-font {
        font-family: var(--font-mono, monospace);
        line-height: 1.5;
    }

    .checkbox-label {
        display: flex;
        align-items: center;
        gap: 0.6rem;
        font-size: 0.85rem;
        cursor: pointer;
    }

    .form-actions {
        display: flex;
        gap: 1rem;
        margin-top: 0.5rem;
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

    .btn-outline {
        background: transparent;
        border: 1px solid var(--border);
        color: var(--foreground);
    }

    /* List */
    .articles-list {
        display: flex;
        flex-direction: column;
    }

    .list-item {
        padding: 1.25rem 1.5rem;
        border-bottom: 1px solid var(--border);
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        gap: 1rem;
        transition: background 0.15s ease;
    }

    .list-item:hover {
        background: var(--muted);
    }

    .list-item.active {
        background: oklch(from var(--primary) l c h / 0.1);
        border-left: 4px solid var(--primary);
    }

    .item-meta {
        display: flex;
        gap: 0.5rem;
        align-items: center;
        margin-bottom: 0.4rem;
        font-size: 0.75rem;
    }

    .cat {
        font-weight: 700;
        color: var(--primary);
    }

    .status.pub {
        color: var(--success);
    }

    .status.draft {
        color: var(--muted-foreground);
    }

    .item-title {
        font-size: 0.95rem;
        font-weight: 600;
        margin-bottom: 0.3rem;
    }

    .item-summary {
        font-size: 0.8rem;
        color: var(--muted-foreground);
        line-height: 1.4;
    }

    .item-actions {
        display: flex;
        gap: 0.3rem;
    }

    .icon-btn {
        background: transparent;
        border: none;
        color: var(--muted-foreground);
        padding: 0.4rem;
        border-radius: 4px;
        cursor: pointer;
    }

    .icon-btn:hover {
        color: var(--foreground);
        background: var(--muted);
    }

    .icon-btn.danger:hover {
        color: var(--destructive);
    }

    .state-msg {
        padding: 2rem;
        text-align: center;
        color: var(--muted-foreground);
    }

    @media (max-width: 1100px) {
        .editor-grid {
            grid-template-columns: 1fr;
        }
    }
</style>
