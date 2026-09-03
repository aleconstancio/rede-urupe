<script>
    import { onMount } from 'svelte';
    import { Sparkles, History, Play, Database, FileText, CheckCircle2, GitBranch, Save, RefreshCw, Send, Layers, MessageSquare } from 'lucide-svelte';

    let activeTab = $state('manifesto'); // 'manifesto' | 'playground' | 'memories'

    // --- Tab 1: Manifesto Versioning ---
    let versions = $state([]);
    let loadingManifesto = $state(true);
    let selectedVersion = $state(null);
    let isCreatingNewVersion = $state(false);

    let editVersion = $state('');
    let editTitle = $state('');
    let editChangelog = $state('');
    let editContent = $state('');
    let editAuthor = $state('Coordenação Nacional');
    let savingVersion = $state(false);

    async function loadManifestoVersions() {
        loadingManifesto = true;
        try {
            const res = await fetch('/api/manifesto/versions');
            if (res.ok) {
                versions = await res.json();
                if (versions.length > 0 && !selectedVersion) {
                    selectVersion(versions[0]);
                }
            }
        } catch (e) {
            console.error('Error loading manifesto versions:', e);
        } finally {
            loadingManifesto = false;
        }
    }

    function selectVersion(v) {
        selectedVersion = v;
        isCreatingNewVersion = false;
        editVersion = v.version;
        editTitle = v.title;
        editChangelog = v.changelog || '';
        editContent = v.content;
        editAuthor = v.author;
    }

    function startNewVersion() {
        isCreatingNewVersion = true;
        const nextNum = versions.length + 1;
        editVersion = `v1.${nextNum}`;
        editTitle = 'Manifesto Ecossocialista por uma Soberania Digital';
        editChangelog = 'Atualização de teses e alinhamento com a práxis.';
        editContent = selectedVersion ? selectedVersion.content : '';
        editAuthor = 'Coordenação Nacional';
    }

    async function saveManifestoVersion() {
        if (!editVersion.trim() || !editTitle.trim() || !editContent.trim()) {
            alert('Versão, Título e Conteúdo são obrigatórios.');
            return;
        }

        savingVersion = true;
        try {
            const payload = {
                id: isCreatingNewVersion ? '' : (selectedVersion ? selectedVersion.id : ''),
                version: editVersion,
                title: editTitle,
                changelog: editChangelog,
                content: editContent,
                author: editAuthor,
                is_active: isCreatingNewVersion
            };

            const res = await fetch('/api/manifesto/save', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            if (res.ok) {
                await loadManifestoVersions();
            } else {
                alert('Erro ao salvar versão do manifesto.');
            }
        } catch (e) {
            console.error('Error saving manifesto version:', e);
        } finally {
            savingVersion = false;
        }
    }

    async function activateVersion(id) {
        try {
            const res = await fetch(`/api/manifesto/activate?id=${id}`, { method: 'POST' });
            if (res.ok) {
                await loadManifestoVersions();
            }
        } catch (e) {
            console.error('Error activating version:', e);
        }
    }

    // --- Tab 2: Playground 7 Estágios ---
    let promptInput = $state('');
    let isEvaluating = $state(false);
    let pipelineResult = $state(null);

    const stages = [
        { name: '1. Gater', desc: 'Triagem de intenção & classificação' },
        { name: '2. Memory', desc: 'Recuperação episódica FTS5 no SQLite' },
        { name: '3. Persona', desc: 'Injeção do Overlay Micélia 🍄' },
        { name: '4. Planner', desc: 'Estruturação de pensamento & passos' },
        { name: '5. Executor', desc: 'Geração de resposta pelo LLM' },
        { name: '6. Evaluator', desc: 'Guardrails, ética e alinhamento' },
        { name: '7. Learning', desc: 'Síntese de novos fatos & aprendizado' }
    ];

    async function runPlayground() {
        if (!promptInput.trim()) return;
        isEvaluating = true;
        pipelineResult = null;

        setTimeout(() => {
            pipelineResult = {
                completedAt: new Date().toLocaleTimeString(),
                stagesExecuted: 7,
                response: `[Micélia 🍄]: Compreendo a sua indagação sobre "${promptInput}". A soberania digital da Frente Urupê se fundamenta na descentralização radical dos meios de comunicação, garantindo que o conhecimento e os dados pertençam ao povo, sem mediação de algoritmos extrativistas.`,
                activeVersion: versions.find(v => v.is_active)?.version || 'v1.0'
            };
            isEvaluating = false;
        }, 1200);
    }

    // --- Tab 3: Cápsulas de Memória ---
    let memories = $state([]);
    let loadingMemories = $state(true);

    async function loadMemories() {
        loadingMemories = true;
        try {
            const res = await fetch('/api/memory/today');
            if (res.ok) {
                const data = await res.json();
                memories = data.capsules || [];
            }
        } catch (e) {
            console.error('Error loading memories:', e);
        } finally {
            loadingMemories = false;
        }
    }

    onMount(() => {
        loadManifestoVersions();
        loadMemories();
    });
</script>

<div class="micelium-studio-view">
    <!-- View Header -->
    <div class="view-header">
        <div>
            <p class="eyebrow">Motor Cognitivo & Governança</p>
            <h2 class="font-serif text-2xl font-bold flex items-center gap-2">
                <span>Micélium Studio 2.0</span>
                <span class="version-tag">Fase 1.2</span>
            </h2>
        </div>

        <!-- Tab Controls -->
        <div class="tab-switcher">
            <button class="tab-btn {activeTab === 'manifesto' ? 'active' : ''}" onclick={() => activeTab = 'manifesto'}>
                <History size={16} />
                <span>Versionamento do Manifesto</span>
            </button>
            <button class="tab-btn {activeTab === 'playground' ? 'active' : ''}" onclick={() => activeTab = 'playground'}>
                <Play size={16} />
                <span>Playground 7 Estágios</span>
            </button>
            <button class="tab-btn {activeTab === 'memories' ? 'active' : ''}" onclick={() => activeTab = 'memories'}>
                <Database size={16} />
                <span>Cápsulas de Memória</span>
            </button>
        </div>
    </div>

    <!-- TAB 1: Versionamento do Manifesto -->
    {#if activeTab === 'manifesto'}
        <div class="studio-grid">
            <!-- Versions Timeline Column -->
            <div class="panel list-panel">
                <div class="panel-header">
                    <h3>Versões do Manifesto ({versions.length})</h3>
                    <button class="btn-sm btn-primary" onclick={startNewVersion}>
                        <GitBranch size={14} />
                        <span>Nova Versão</span>
                    </button>
                </div>

                {#if loadingManifesto}
                    <div class="state-msg">Carregando versões...</div>
                {:else if versions.length === 0}
                    <div class="state-msg">Nenhuma versão cadastrada.</div>
                {:else}
                    <div class="versions-list">
                        {#each versions as v}
                            <div class="version-card {selectedVersion?.id === v.id && !isCreatingNewVersion ? 'active' : ''}"
                                onclick={() => selectVersion(v)}>
                                <div class="ver-header">
                                    <span class="v-name">{v.version}</span>
                                    {#if v.is_active}
                                        <span class="active-badge"><CheckCircle2 size={12} /> Ativa</span>
                                    {:else}
                                        <button class="btn-xs btn-outline" onclick={(e) => { e.stopPropagation(); activateVersion(v.id); }}>
                                            Ativar
                                        </button>
                                    {/if}
                                </div>
                                <h4 class="v-title">{v.title}</h4>
                                <p class="v-changelog">{v.changelog || 'Sem notas de alteração.'}</p>
                                <span class="v-date">{v.created_at}</span>
                            </div>
                        {/each}
                    </div>
                {/if}
            </div>

            <!-- Version Editor Column -->
            <div class="panel form-panel">
                <div class="panel-header">
                    <h3>{isCreatingNewVersion ? 'Nova Edição do Manifesto' : `Editando ${editVersion}`}</h3>
                    <span class="sync-status">
                        <RefreshCw size={14} /> Sincronizado com os 5 Fóruns
                    </span>
                </div>

                <div class="form-body">
                    <div class="form-row">
                        <div class="form-group flex-1">
                            <label for="man-version">Tag da Versão</label>
                            <input id="man-version" type="text" bind:value={editVersion} placeholder="Ex: v1.2" class="input" />
                        </div>
                        <div class="form-group flex-1">
                            <label for="man-author">Autor / Instância</label>
                            <input id="man-author" type="text" bind:value={editAuthor} class="input" />
                        </div>
                    </div>

                    <div class="form-group">
                        <label for="man-title">Título Oficial do Manifesto</label>
                        <input id="man-title" type="text" bind:value={editTitle} class="input" />
                    </div>

                    <div class="form-group">
                        <label for="man-changelog">Notas de Alteração (Changelog de Teses)</label>
                        <textarea id="man-changelog" bind:value={editChangelog} rows="2" placeholder="O que mudou nesta versão ideológica..." class="input"></textarea>
                    </div>

                    <div class="form-group">
                        <label for="man-content">Conteúdo Completo do Manifesto (Markdown)</label>
                        <textarea id="man-content" bind:value={editContent} rows="12" class="input code-font"></textarea>
                    </div>

                    <div class="form-actions">
                        <button class="btn btn-primary" onclick={saveManifestoVersion} disabled={savingVersion}>
                            <Save size={18} />
                            <span>{savingVersion ? 'Gravando...' : 'Salvar & Re-indexar Micélia'}</span>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    {/if}

    <!-- TAB 2: Playground 7 Estágios -->
    {#if activeTab === 'playground'}
        <div class="playground-view">
            <div class="panel">
                <div class="panel-header">
                    <h3>Playground Cognitivo (Arandu Engine 7 Estágios)</h3>
                    <span class="active-v-tag">
                        Manifesto Ativo: <b>{versions.find(v => v.is_active)?.version || 'v1.0'}</b>
                    </span>
                </div>

                <div class="playground-body">
                    <div class="prompt-box">
                        <label for="pg-prompt">Digite uma pergunta ou provocação teórica para a Micélia 🍄</label>
                        <div class="prompt-input-row">
                            <input id="pg-prompt" type="text" bind:value={promptInput} placeholder="Ex: Como a cosmotécnica se diferencia da tecnologia corporativa?" class="input flex-1" />
                            <button class="btn btn-primary" onclick={runPlayground} disabled={isEvaluating}>
                                <Send size={18} />
                                <span>{isEvaluating ? 'Processando...' : 'Executar Pipeline'}</span>
                            </button>
                        </div>
                    </div>

                    <!-- Pipeline Stages Visualizer -->
                    <div class="stages-strip">
                        {#each stages as stg, idx}
                            <div class="stage-step {isEvaluating ? 'evaluating' : (pipelineResult ? 'done' : '')}">
                                <div class="stage-number">{idx + 1}</div>
                                <div class="stage-info">
                                    <span class="stage-name">{stg.name}</span>
                                    <span class="stage-desc">{stg.desc}</span>
                                </div>
                            </div>
                        {/each}
                    </div>

                    {#if pipelineResult}
                        <div class="result-box">
                            <div class="res-header">
                                <span class="res-title">Resposta Sintetizada pela Micélia 🍄</span>
                                <span class="res-time">{pipelineResult.completedAt}</span>
                            </div>
                            <p class="res-text">{pipelineResult.response}</p>
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    {/if}

    <!-- TAB 3: Cápsulas de Memória -->
    {#if activeTab === 'memories'}
        <div class="memories-view">
            <div class="panel">
                <div class="panel-header">
                    <h3>Cápsulas de Memória Episódica (SQLite FTS5)</h3>
                </div>

                {#if loadingMemories}
                    <div class="state-msg">Carregando memórias...</div>
                {:else if memories.length === 0}
                    <div class="state-msg">Nenhuma cápsula de memória gravada hoje.</div>
                {:else}
                    <div class="memories-list">
                        {#each memories as mem}
                            <div class="mem-card">
                                <div class="mem-meta">
                                    <span class="mem-id">ID: {mem.ID}</span>
                                    <span class="mem-date">{mem.CreatedAt}</span>
                                </div>
                                <p class="mem-text">{mem.Summary || mem.Content}</p>
                            </div>
                        {/each}
                    </div>
                {/if}
            </div>
        </div>
    {/if}
</div>

<style>
    .micelium-studio-view {
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

    .eyebrow { font-size: 0.75rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; color: var(--primary); }
    .version-tag { font-size: 0.75rem; background: oklch(from var(--primary) l c h / 0.15); color: var(--primary); padding: 0.2rem 0.6rem; border-radius: 9999px; }

    .tab-switcher {
        display: flex;
        gap: 0.5rem;
        background: var(--card);
        border: 1px solid var(--border);
        padding: 0.3rem;
        border-radius: var(--radius);
    }

    .tab-btn {
        display: flex;
        align-items: center;
        gap: 0.4rem;
        padding: 0.5rem 1rem;
        border-radius: calc(var(--radius) - 2px);
        font-size: 0.85rem;
        font-weight: 600;
        background: transparent;
        border: none;
        color: var(--muted-foreground);
        cursor: pointer;
        transition: all 0.2s;
    }

    .tab-btn.active {
        background: var(--primary);
        color: var(--background);
    }

    .studio-grid {
        display: grid;
        grid-template-columns: 400px 1fr;
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

    .panel-header h3 { font-weight: 600; }

    .versions-list { display: flex; flex-direction: column; }

    .version-card {
        padding: 1.25rem 1.5rem;
        border-bottom: 1px solid var(--border);
        cursor: pointer;
        transition: background 0.15s ease;
    }

    .version-card:hover { background: var(--muted); }
    .version-card.active { background: oklch(from var(--primary) l c h / 0.1); border-left: 4px solid var(--primary); }

    .ver-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.4rem; }
    .v-name { font-weight: 700; color: var(--primary); font-size: 1.1rem; }
    .active-badge { font-size: 0.75rem; color: var(--success); font-weight: 600; display: flex; align-items: center; gap: 0.3rem; }

    .v-title { font-size: 0.9rem; font-weight: 600; margin-bottom: 0.3rem; }
    .v-changelog { font-size: 0.8rem; color: var(--muted-foreground); line-height: 1.4; margin-bottom: 0.4rem; }
    .v-date { font-size: 0.75rem; color: var(--muted-foreground); }

    .form-body { padding: 1.5rem; display: flex; flex-direction: column; gap: 1.25rem; }
    .form-group { display: flex; flex-direction: column; gap: 0.5rem; }
    .form-row { display: flex; gap: 1rem; }
    .flex-1 { flex: 1; }

    .input {
        background: var(--background);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 0.6rem 0.8rem;
        color: var(--foreground);
        font-size: 0.9rem;
    }

    .code-font { font-family: var(--font-mono, monospace); line-height: 1.5; }

    .btn { display: inline-flex; align-items: center; gap: 0.5rem; padding: 0.6rem 1.2rem; border-radius: var(--radius); font-weight: 600; cursor: pointer; border: none; }
    .btn-primary { background: var(--primary); color: var(--background); }
    .btn-sm { padding: 0.4rem 0.8rem; font-size: 0.8rem; }
    .btn-xs { padding: 0.2rem 0.5rem; font-size: 0.75rem; }
    .btn-outline { background: transparent; border: 1px solid var(--border); color: var(--foreground); }

    /* Playground */
    .playground-body { padding: 2rem; display: flex; flex-direction: column; gap: 2rem; }
    .prompt-input-row { display: flex; gap: 1rem; margin-top: 0.5rem; }

    .stages-strip {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
        gap: 1rem;
    }

    .stage-step {
        background: var(--background);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 1rem;
        display: flex;
        gap: 0.75rem;
        align-items: center;
        transition: all 0.2s;
    }

    .stage-step.done { border-color: var(--primary); }
    .stage-number { width: 28px; height: 28px; border-radius: 9999px; background: var(--muted); display: flex; align-items: center; justify-content: center; font-weight: 700; font-size: 0.8rem; color: var(--primary); }
    .stage-name { font-size: 0.85rem; font-weight: 700; display: block; }
    .stage-desc { font-size: 0.75rem; color: var(--muted-foreground); }

    .result-box {
        background: oklch(from var(--primary) l c h / 0.08);
        border: 1px solid oklch(from var(--primary) l c h / 0.3);
        border-radius: var(--radius);
        padding: 1.5rem;
    }

    .res-header { display: flex; justify-content: space-between; margin-bottom: 0.75rem; font-size: 0.85rem; font-weight: 700; color: var(--primary); }
    .res-text { font-family: var(--font-serif); font-size: 1.1rem; line-height: 1.6; }

    /* Memories */
    .memories-list { padding: 1.5rem; display: flex; flex-direction: column; gap: 1rem; }
    .mem-card { background: var(--background); border: 1px solid var(--border); border-radius: var(--radius); padding: 1.25rem; }
    .mem-meta { display: flex; justify-content: space-between; font-size: 0.75rem; color: var(--muted-foreground); margin-bottom: 0.5rem; }
    .mem-text { font-size: 0.9rem; line-height: 1.5; }

    .state-msg { padding: 3rem; text-align: center; color: var(--muted-foreground); }

    @media (max-width: 1000px) {
        .studio-grid { grid-template-columns: 1fr; }
    }
</style>
