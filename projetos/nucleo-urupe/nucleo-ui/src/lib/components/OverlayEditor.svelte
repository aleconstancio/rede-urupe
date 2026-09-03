<script>
  import { Button, Input, Select, Checkbox, Switch, Dialog, EmptyState, Skeleton, Badge, Alert, Spinner } from '@talos/ui';

    import { onMount } from 'svelte';

    /** @type {any[]} */
    let overlays = [];
    /** @type {any[]} */
    let identities = [];
    let loading = true;
    /** @type {any} */
    let editing = null;
    let dialogOpen = $state(false);

    function openEditor(data) { editing = data; dialogOpen = true; }
    function closeEditor() { editing = null; dialogOpen = false; }

    async function fetchData() {
        loading = true;
        try {
            const [ovRes, idRes] = await Promise.all([
                fetch('/api/personas'),
                fetch('/api/identities')
            ]);
            if (ovRes.ok) overlays = (await ovRes.json()).items || [];
            if (idRes.ok) identities = (await idRes.json()).items || [];
        } catch (err) {
            console.error("Failed to fetch data", err);
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
        
        // @ts-ignore
        data.is_enabled = formData.get('is_enabled') === 'on';
        // @ts-ignore
        data.is_default = formData.get('is_default') === 'on';
        // @ts-ignore
        data.sort_order = parseInt(/** @type {string} */(formData.get('sort_order'))) || 0;
        
        // @ts-ignore
        try { data.traits = JSON.parse(data.traits || "{}"); } catch(e) { data.traits = {}; }
        // @ts-ignore
        data.allowed_intents = data.allowed_intents ? data.allowed_intents.split(',').map(v => v.trim()).filter(v => v !== "") : [];

        try {
            const res = await fetch('/api/personas', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
            if (res.ok) {
                editing = null;
                dialogOpen = false;
                closeEditor();
                await fetchData();
            }
        } catch (err) {
            console.error("Failed to save overlay", err);
        }
    }

    onMount(fetchData);
</script>

<div class="manager">
    <div class="section-header">
        <h3>Overlays de Estilo</h3>
        <Button class="add-btn" onclick={() => openEditor({ id: "", identity_id: identities[0]?.id || "", name: "", style_prompt: "", is_enabled: true, sort_order: 0 })}>
            + Novo Overlay
        </Button>
    </div>

    {#if loading}
        <PageLoading type="text" lines={3} />
    {:else}
        <div class="overlay-list">
            {#each overlays as o}
                <div class="overlay-item {o.is_default ? 'default' : ''}">
                    <div class="item-info">
                        <strong>{o.name}</strong>
                        <span class="sub-text">Asssociado a: {identities.find(id => id.id === o.identity_id)?.name || o.identity_id}</span>
                    </div>
                    <div class="item-actions">
                        <Button class="edit-btn" onclick={() => openEditor({...o})}>Editar</Button>
                    </div>
                </div>
            {/each}
        </div>
    {/if}

    <Dialog bind:open={dialogOpen} title={editing?.id ? 'Editar Overlay' : 'Novo Overlay'}>
      {#if editing}
        <form onsubmit={handleSave}>
          <input type="hidden" name="id" value={editing.id} />
          <div class="field">
            <label for="ov-identity">Identidade Pai</label>
            <select id="ov-identity" name="identity_id" value={editing.identity_id}>
              {#each identities as id}
                <option value={id.id}>{id.name}</option>
              {/each}
            </select>
          </div>
          <div class="field">
            <label for="ov-name">Nome do Estilo</label>
            <input id="ov-name" name="name" value={editing.name} required />
          </div>
          <div class="field">
            <label for="ov-desc">Descrição</label>
            <input id="ov-desc" name="description" value={editing.description || ""} />
          </div>
          <div class="field">
            <label for="ov-prompt">Style Prompt</label>
            <textarea id="ov-prompt" name="style_prompt" required>{editing.style_prompt}</textarea>
          </div>
          <div class="field">
            <label for="ov-traits">Traits (JSON)</label>
            <input id="ov-traits" name="traits" value={JSON.stringify(editing.traits || {})} />
          </div>
          <div class="field">
            <label for="ov-intents">Intenções Permitidas (vírgula)</label>
            <input id="ov-intents" name="allowed_intents" value={editing.allowed_intents?.join(', ') || ""} />
          </div>
          <div class="field">
            <label for="ov-order">Ordem de Prioridade</label>
            <input id="ov-order" type="number" name="sort_order" value={editing.sort_order} />
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

    .overlay-list {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .overlay-item {
        background: rgba(255,255,255,0.02);
        border: 1px solid var(--border);
        border-radius: 8px;
        padding: 0.75rem 1rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .overlay-item.default {
        border-left: 3px solid var(--primary);
    }

    .item-info {
        display: flex;
        flex-direction: column;
    }

    .sub-text {
        font-size: 0.75rem;
        color: var(--muted-foreground);
    }

    .edit-btn {
        background: none;
        border: none;
        color: var(--primary);
        font-size: 0.8rem;
        font-weight: 600;
        cursor: pointer;
    }

    .loading { padding: 1rem; color: var(--muted-foreground); font-size: 0.8rem; }
</style>
