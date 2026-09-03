<script>
  import { Button, Input, Select, Checkbox, Switch, Dialog, EmptyState, Skeleton, Badge, Alert, Spinner } from '@talos/ui';

    import { persona } from '../stores';
    import { Save, RefreshCw, AlertCircle, CheckCircle } from 'lucide-svelte';
    import { onMount } from 'svelte';

    let state = $derived($persona);
    
    let policy = {
        ChannelID: "",
        DefaultIdentityID: "",
        DefaultPersonaID: "",
        SelectionMode: "fixed",
        AllowedIdentityIDs: [],
        AllowedPersonaIDs: [],
        /** @type {Record<string, string>} */
        IntentPersonaMap: {},
        /** @type {Record<string, string>} */
        ModePersonaMap: {},
        ManualOverrideIdentityID: "",
        ManualOverridePersonaID: ""
    };

    let loading = false;
    let saving = false;
    /** @type {{type: string, text: string} | null} */
    let message = null;

    async function fetchPolicy() {
        loading = true;
        try {
            const res = await fetch('/api/persona-policy');
            if (res.ok) {
                const data = await res.json();
                if (data && data.ChannelID) {
                    policy = data;
                    if (!policy.IntentPersonaMap) policy.IntentPersonaMap = {};
                    if (!policy.ModePersonaMap) policy.ModePersonaMap = {};
                }
            }
        } catch (err) {
            console.error("Failed to fetch policy", err);
        } finally {
            loading = false;
        }
    }

    async function savePolicy() {
        saving = true;
        message = null;
        try {
            const res = await fetch('/api/persona-policy', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(policy)
            });
            if (res.ok) {
                message = { type: 'success', text: 'Política salva com sucesso!' };
                setTimeout(() => message = null, 3000);
            } else {
                message = { type: 'error', text: 'Falha ao salvar política.' };
            }
        } catch (err) {
            message = { type: 'error', text: 'Erro de conexão.' };
        } finally {
            saving = false;
        }
    }

    onMount(fetchPolicy);

    // Derived lists from state
    let identities = $derived(state?.Identities || []);
    let overlays = $derived(state?.Overlays || []);
</script>

<div class="policy-editor">
    <div class="editor-header">
        <div class="title-desc">
            <h3>Política de Persona por Canal</h3>
            <p>Define qual bot e qual estilo são usados em cada situação.</p>
        </div>
        <Button class="btn-save" onclick={savePolicy} disabled={saving}>
            {#if saving}
                <RefreshCw size={14} class="spin" />
            {:else}
                <Save size={14} />
            {/if}
            Salvar
        </Button>
    </div>

    {#if loading}
        <PageLoading type="text" lines={3} />
    {:else}
        <div class="grid">
            <div class="card settings">
                <h4>Configurações Base</h4>
                <div class="field">
                    <label for="selection-mode">Modo de Seleção</label>
                    <Select  id="selection-mode" bind:value={policy.SelectionMode}>
                        <option value="fixed">Fixo (Sempre o mesmo)</option>
                        <option value="intent_map">Por Intenção (Dinâmico)</option>
                        <option value="mode_map">Por Modo (Ambiental)</option>
                    </Select>
                </div>

                <div class="field">
                    <label for="default-identity">Identidade Padrão</label>
                    <Select  id="default-identity" bind:value={policy.DefaultIdentityID}>
                        <option value="">Selecione uma identidade...</option>
                        {#each identities as id}
                            <option value={id.id}>{id.name}</option>
                        {/each}
                    </Select>
                </div>

                <div class="field">
                    <label for="default-style">Estilo (Overlay) Padrão</label>
                    <Select  id="default-style" bind:value={policy.DefaultPersonaID}>
                        <option value="">Selecione um estilo...</option>
                        {#each overlays as ov}
                            <option value={ov.id}>{ov.name}</option>
                        {/each}
                    </Select>
                </div>
            </div>

            <div class="card override">
                <h4>Override Manual</h4>
                <p class="hint">Força o bot a assumir uma forma específica temporariamente.</p>
                
                <div class="field">
                    <label for="override-identity">Identidade Forçada</label>
                    <Select  id="override-identity" bind:value={policy.ManualOverrideIdentityID}>
                        <option value="">Nenhum (Usar padrão)</option>
                        {#each identities as id}
                            <option value={id.id}>{id.name}</option>
                        {/each}
                    </Select>
                </div>

                <div class="field">
                    <label for="override-style">Estilo Forçado</label>
                    <Select  id="override-style" bind:value={policy.ManualOverridePersonaID}>
                        <option value="">Nenhum (Usar padrão)</option>
                        {#each overlays as ov}
                            <option value={ov.id}>{ov.name}</option>
                        {/each}
                    </Select>
                </div>
            </div>
        </div>

        {#if policy.SelectionMode === 'intent_map'}
            <div class="card map">
                <h4>Mapeamento por Intenção</h4>
                <div class="map-grid">
                    <div class="map-row header">
                        <span>Intenção (System 1/3)</span>
                        <span>Estilo (Overlay) Correspondente</span>
                    </div>
                    {#each ['respond', 'clarify', 'synthesize', 'moderate', 'deescalate'] as intent}
                        <div class="map-row">
                            <label for="intent-{intent}" class="intent-tag">{intent}</label>
                            <Select  id="intent-{intent}" bind:value={policy.IntentPersonaMap[intent]}>
                                <option value="">(Usar padrão)</option>
                                {#each overlays as ov}
                                    <option value={ov.id}>{ov.name}</option>
                                {/each}
                            </Select>
                        </div>
                    {/each}
                </div>
            </div>
        {/if}

        {#if policy.SelectionMode === 'mode_map'}
            <div class="card map">
                <h4>Mapeamento por Modo</h4>
                <div class="map-grid">
                    <div class="map-row header">
                        <span>Modo (Ambiental)</span>
                        <span>Estilo (Overlay) Correspondente</span>
                    </div>
                    {#each ['question', 'philosophy', 'conflict'] as mode}
                        <div class="map-row">
                            <label for="mode-{mode}" class="intent-tag">{mode}</label>
                            <Select  id="mode-{mode}" bind:value={policy.ModePersonaMap[mode]}>
                                <option value="">(Usar padrão)</option>
                                {#each overlays as ov}
                                    <option value={ov.id}>{ov.name}</option>
                                {/each}
                            </Select>
                        </div>
                    {/each}
                </div>
            </div>
        {/if}
    {/if}

    {#if message}
        <div class="toast" class:success={message.type === 'success'} class:error={message.type === 'error'}>
            {#if message.type === 'success'}
                <CheckCircle size={16} />
            {:else}
                <AlertCircle size={16} />
            {/if}
            {message.text}
        </div>
    {/if}
</div>

<style>
    .policy-editor {
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
        position: relative;
    }

    .editor-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .title-desc h3 {
        margin: 0;
        font-size: 1.1rem;
    }

    .title-desc p {
        margin: 0.25rem 0 0;
        font-size: 0.85rem;
        color: var(--muted-foreground);
    }

    .btn-save {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        background: var(--primary);
        color: var(--background);
        border: none;
        padding: 0.6rem 1.2rem;
        border-radius: 8px;
        font-weight: 700;
        cursor: pointer;
        transition: transform 0.2s;
    }

    .btn-save:hover:not(:disabled) {
        transform: translateY(-2px);
    }

    .btn-save:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1.5rem;
    }

    .card {
        background: rgba(255,255,255,0.02);
        border: 1px solid var(--border);
        border-radius: 12px;
        padding: 1.25rem;
    }

    h4 {
        margin: 0 0 1rem;
        font-size: 0.9rem;
        text-transform: uppercase;
        color: var(--muted-foreground);
        letter-spacing: 0.05em;
    }

    .hint {
        font-size: 0.8rem;
        color: var(--muted-foreground);
        margin: -0.5rem 0 1.25rem;
    }

    .field {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        margin-bottom: 1rem;
    }

    label {
        font-size: 0.85rem;
        font-weight: 600;
        color: var(--foreground);
    }

    select {
        background: var(--card);
        border: 1px solid var(--border);
        color: var(--foreground);
        padding: 0.6rem;
        border-radius: 6px;
        font-size: 0.9rem;
        outline: none;
    }

    select:focus {
        border-color: var(--primary);
    }

    .map-grid {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .map-row {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1rem;
        padding: 0.5rem;
        align-items: center;
        border-radius: 6px;
    }

    .map-row.header {
        font-size: 0.75rem;
        font-weight: 700;
        color: var(--muted-foreground);
        text-transform: uppercase;
    }

    .intent-tag {
        font-family: var(--font-mono);
        background: rgba(255,255,255,0.05);
        padding: 4px 8px;
        border-radius: 4px;
        width: fit-content;
    }

    .loading {
        text-align: center;
        padding: 3rem;
        color: var(--muted-foreground);
    }

    .toast {
        position: fixed;
        bottom: 2rem;
        right: 2rem;
        padding: 0.75rem 1.5rem;
        border-radius: 8px;
        display: flex;
        align-items: center;
        gap: 0.75rem;
        font-weight: 600;
        box-shadow: 0 4px 20px rgba(0,0,0,0.4);
        z-index: 1000;
        animation: slideIn 0.3s ease-out;
    }

    .toast.success { background: var(--success); color: white; }
    .toast.error { background: var(--destructive); color: white; }

    @keyframes slideIn {
        from { transform: translateX(100%); opacity: 0; }
        to { transform: translateX(0); opacity: 1; }
    }

    :global(.spin) {
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
    }
</style>
