<script>
  import { Button, Input, Select, Checkbox, Switch, Dialog, EmptyState, Skeleton, Badge, Alert, Spinner } from '@talos/ui';

    import { onMount } from 'svelte';

    /** @type {any[]} */
    let identities = [];
    let loading = true;
    /** @type {any} */
    let editing = null;
    let dialogOpen = $state(false);

    function openEditor(data) { editing = data; dialogOpen = true; }
    function closeEditor() { editing = null; dialogOpen = false; }

    async function fetchIdentities() {
        loading = true;
        try {
            const res = await fetch('/api/identities');
            if (res.ok) {
                const data = await res.json();
                identities = data.items || [];
            }
        } catch (err) {
            console.error("Failed to fetch identities", err);
        } finally {
            loading = false;
        }
    }

    /** @param {SubmitEvent} event */
    async function handleSave(event) {
        event.preventDefault();
        const form = /** @type {HTMLFormElement} */ (event.target);
        const formData = new FormData(form);
        const data = Object.fromEntries(formData.entries());
        
        // Simple validation/casting
        // @ts-ignore
        data.is_enabled = formData.get('is_enabled') === 'on';
        // @ts-ignore
        data.is_default = formData.get('is_default') === 'on';
        // @ts-ignore
        data.core_values = data.core_values ? data.core_values.split(',').map(v => v.trim()).filter(v => v !== "") : [];

        try {
            const res = await fetch('/api/identities', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
            if (res.ok) {
                editing = null;
        dialogOpen = false;
        closeEditor();
                await fetchIdentities();
            }
        } catch (err) {
            console.error("Failed to save identity", err);
        }
    }

    onMount(fetchIdentities);
</script>

<div class="manager">
    <div class="section-header">
        <h3>Identidades Centrais</h3>
        <Button class="add-btn" onclick={() => openEditor({ id: "", name: "", identity_prompt: "", is_enabled: true })}>
            + Nova Identidade
        </Button>
    </div>

    {#if loading}
        <PageLoading type="text" lines={3} />
    {:else}
        <div class="identity-grid">
            {#each identities as id}
                <div class="id-card {id.is_default ? 'default' : ''}">
                    <div class="card-body">
                        <div class="id-title">
                            <strong>{id.name}</strong>
                            {#if id.is_default}<span class="badge">PADRÃO</span>{/if}
                        </div>
                        <p class="id-desc">{id.description || 'Sem descrição.'}</p>
                    </div>
                    <div class="card-footer">
                        <Button class="edit-btn" onclick={() => openEditor({...id})}>Editar</Button>
                    </div>
                </div>
            {/each}
        </div>
    {/if}

    <Dialog bind:open={dialogOpen} title={editing?.id ? 'Editar Identidade' : 'Nova Identidade'}>
      {#if editing}
        <form onsubmit={handleSave}>
          <input type="hidden" name="id" value={editing.id} />
          <div class="field">
            <label for="id-name">Nome</label>
            <input id="id-name" name="name" value={editing.name} required />
          </div>
          <div class="field">
            <label for="id-desc">Descrição</label>
            <input id="id-desc" name="description" value={editing.description} />
          </div>
          <div class="field">
            <label for="id-prompt">Identity Prompt</label>
            <textarea id="id-prompt" name="identity_prompt" required>{editing.identity_prompt}</textarea>
          </div>
          <div class="field">
            <label for="id-values">Valores Centrais (separados por vírgula)</label>
            <input id="id-values" name="core_values" value={editing.core_values?.join(', ') || ''} placeholder="ex: humor seco, lealdade, precisao" />
          </div>
          {#snippet actions()}
            <Button variant="ghost" type="button" onclick={closeEditor}>Cancelar</Button>
            <Button type="submit">Salvar</Button>
          {/snippet}
        </form>
      {/if}
    </Dialog>
</div>

<style>
    .manager {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .section-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    h3 {
        margin: 0;
        font-size: 0.9rem;
        text-transform: uppercase;
        color: var(--muted-foreground);
    }

    .add-btn {
        background: var(--muted);
        color: var(--primary);
        border: none;
        padding: 4px 12px;
        border-radius: 4px;
        font-size: 0.8rem;
        font-weight: 700;
        cursor: pointer;
    }

    .identity-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
        gap: 1rem;
    }

    .id-card {
        background: rgba(255,255,255,0.02);
        border: 1px solid var(--border);
        border-radius: 8px;
        display: flex;
        flex-direction: column;
        transition: transform 0.2s;
    }

    .id-card.default {
        border-color: var(--primary);
        background: oklch(from var(--primary) l c h / 0.05);
    }

    .card-body {
        padding: 1rem;
        flex: 1;
    }

    .id-title {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 0.5rem;
    }

    .badge {
        font-size: 0.6rem;
        background: var(--primary);
        color: var(--background);
        padding: 1px 4px;
        border-radius: 2px;
        font-weight: 800;
    }

    .id-desc {
        font-size: 0.8rem;
        color: var(--muted-foreground);
        margin: 0;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }

    .card-footer {
        padding: 0.5rem 1rem;
        border-top: 1px solid var(--border);
    }

    .edit-btn {
        background: none;
        border: none;
        color: var(--primary);
        font-size: 0.8rem;
        font-weight: 600;
        cursor: pointer;
        padding: 0;
    }
    
    .loading { padding: 1rem; color: var(--muted-foreground); font-size: 0.8rem; }
</style>
