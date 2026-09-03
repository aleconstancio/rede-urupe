<script>
  import { Button, Input, Select, Checkbox, Switch, Dialog, EmptyState, Skeleton, Badge, Alert, Spinner, Tabs, TabsList, TabsTrigger, TabsContent } from '@talos/ui';

    import { persona } from '../stores';
    import IdentityManager from './IdentityManager.svelte';
    import OverlayEditor from './OverlayEditor.svelte';
    import PolicyEditor from './PolicyEditor.svelte';
    import ProposalReviewer from './ProposalReviewer.svelte';

    let state = $derived($persona);
    let activeTab = $state('identities');
</script>

<div class="panel studio">
    <div class="panel-header">
        <div class="header-left">
            <h2>Persona Studio <span class="v-tag">v5.0</span></h2>
        </div>
    </div>

    <div class="studio-content">
        <Tabs bind:value={activeTab}>
          <TabsList>
            <TabsTrigger value="identities">Identidades</TabsTrigger>
            <TabsTrigger value="overlays">Estilos (Overlays)</TabsTrigger>
            <TabsTrigger value="policy">Política de Canal</TabsTrigger>
            <TabsTrigger value="proposals">Propostas</TabsTrigger>
            <TabsTrigger value="adaptive">Memória Adaptativa</TabsTrigger>
          </TabsList>
          <TabsContent value="identities"><IdentityManager /></TabsContent>
          <TabsContent value="overlays"><OverlayEditor /></TabsContent>
          <TabsContent value="policy"><PolicyEditor /></TabsContent>
          <TabsContent value="proposals"><ProposalReviewer channelId={state?.ChannelID || ""} /></TabsContent>
          <TabsContent value="adaptive">
            <div class="adaptive-view">
              <h3>Memórias Adaptativas Ativas</h3>
              {#if state && state.AdaptiveMemories && state.AdaptiveMemories.length > 0}
                <div class="mem-list">{#each state.AdaptiveMemories as mem}
                  <div class="mem-card">
                    <div class="mem-header">
                      <strong>{mem.IdentityID} / {mem.PersonaID}</strong>
                      <span class="conf">{(mem.Confidence * 100).toFixed(0)}% conf</span>
                    </div>
                    <pre>{JSON.stringify(mem.AdaptiveStyle, null, 2)}</pre>
                  </div>
                {/each}</div>
              {:else}
                <div class="empty">Nenhuma adaptação ativa para este canal.</div>
              {/if}
            </div>
          </TabsContent>
        </Tabs>
    </div>
</div>

<style>
    .studio {
        display: flex;
        flex-direction: column;
        height: 100%;
        overflow: hidden;
    }

    .panel-header {
        padding: 1.25rem 1.5rem;
        border-bottom: 1px solid var(--border);
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .header-left {
        display: flex;
        align-items: center;
        gap: 0.75rem;
    }

    h2 {
        margin: 0;
        font-size: 1.25rem;
        font-weight: 800;
        letter-spacing: -0.02em;
    }

    .v-tag {
        font-size: 0.7rem;
        background: rgba(255,255,255,0.1);
        padding: 2px 6px;
        border-radius: 4px;
        color: var(--muted-foreground);
        font-weight: 600;
    }

    .studio-content {
        padding: 1.5rem;
        flex: 1;
        overflow-y: auto;
    }

    .adaptive-view h3 {
        font-size: 0.9rem;
        text-transform: uppercase;
        color: var(--muted-foreground);
        margin-top: 0;
        margin-bottom: 1rem;
    }

    .mem-list {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .mem-card {
        background: rgba(255,255,255,0.02);
        border: 1px solid var(--border);
        border-radius: 8px;
        padding: 1rem;
    }

    .mem-header {
        display: flex;
        justify-content: space-between;
        margin-bottom: 0.5rem;
    }

    .conf {
        font-size: 0.8rem;
        color: var(--muted-foreground);
    }

    pre {
        background: rgba(0,0,0,0.3);
        padding: 0.5rem;
        border-radius: 6px;
        font-size: 0.8rem;
        margin: 0;
    }

    .empty {
        padding: 2rem;
        text-align: center;
        color: var(--muted-foreground);
        border: 1px dashed var(--border);
        border-radius: 12px;
    }
</style>
