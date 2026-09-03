<script>
  import { Button, Input, Select, Checkbox, Switch, Dialog, EmptyState, Skeleton, Badge, Alert, Spinner, Tabs, TabsList, TabsTrigger, TabsContent } from '@talos/ui';

    import { onMount } from 'svelte';
    import { Shield, Users, AlertTriangle, Activity, Settings, Search } from 'lucide-svelte';

    let activeTab = $state('members');
    let members = [];
    let auditLog = [];
    let warnings = [];
    let welcomeConfig = { enabled: false, channel_id: '', welcome_message: '', goodbye_message: '' };
    let searchQuery = '';
    let loading = false;

    onMount(() => { loadMembers(); loadAuditLog(); loadWelcomeConfig(); });

    $effect(() => {
      if (activeTab === 'modlog') loadWarnings();
      if (activeTab === 'audit') loadAuditLog();
      if (activeTab === 'members') loadMembers();
    });

    async function loadMembers() {
        loading = true;
        try {
            const res = await fetch('/api/admin/members?limit=50');
            if (res.ok) {
                const data = await res.json();
                members = data.members || [];
            }
        } catch (err) {
            console.error('Failed to load members:', err);
        }
        loading = false;
    }

    async function loadAuditLog() {
        try {
            const res = await fetch('/api/admin/audit?limit=30');
            if (res.ok) {
                const data = await res.json();
                auditLog = data.events || [];
            }
        } catch (err) {
            console.error('Failed to load audit log:', err);
        }
    }

    async function loadWarnings() {
        try {
            const res = await fetch('/api/admin/modlog?limit=30');
            if (res.ok) {
                const data = await res.json();
                warnings = data.events || [];
            }
        } catch (err) {
            console.error('Failed to load warnings:', err);
        }
    }

    async function loadWelcomeConfig() {
        try {
            const res = await fetch('/api/admin/welcome');
            if (res.ok) {
                welcomeConfig = await res.json();
            }
        } catch (err) {
            console.error('Failed to load welcome config:', err);
        }
    }

    async function saveWelcomeConfig() {
        try {
            const res = await fetch('/api/admin/welcome', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(welcomeConfig)
            });
            if (res.ok) {
                alert('Configuracao salva!');
            }
        } catch (err) {
            console.error('Failed to save welcome config:', err);
        }
    }

    function getActionIcon(action) {
        const icons = {
            warn: '⚠️', kick: '👢', ban: '🔨', delete: '🗑️',
            timeout: '⏰', join: '➡️', leave: '⬅️', clear_warnings: '✅'
        };
        return icons[action] || '📋';
    }

    let filteredMembers = $derived(searchQuery
        ? members.filter(m => m.name.toLowerCase().includes(searchQuery.toLowerCase()) || m.discord_id.includes(searchQuery))
        : members);
</script>

<div class="admin-view">
    <div class="admin-header">
        <Shield size={24} />
        <h2>Painel de Administracao</h2>
    </div>

  <Tabs bind:value={activeTab}>
    <TabsList>
      <TabsTrigger value="members"><Users size={16} /> Membros</TabsTrigger>
      <TabsTrigger value="modlog"><AlertTriangle size={16} /> Moderacao</TabsTrigger>
      <TabsTrigger value="audit"><Shield size={16} /> Auditoria</TabsTrigger>
      <TabsTrigger value="welcome"><Settings size={16} /> Boas-vindas</TabsTrigger>
    </TabsList>

    <TabsContent value="members">
      <div class="search-bar"><Search size={16} /><input type="text" placeholder="Buscar membros..." bind:value={searchQuery} /></div>
      <div class="members-grid">
        {#each filteredMembers as member}
          <div class="member-card">
            <div class="member-avatar">{member.name[0]}</div>
            <div class="member-info">
              <span class="member-name">{member.name}</span>
              <span class="member-id">{member.discord_id}</span>
              {#if member.notes}<span class="member-notes">{member.notes}</span>{/if}
            </div>
            <div class="member-meta"><span class="member-updated">{new Date(member.updated_at).toLocaleDateString('pt-BR')}</span></div>
          </div>
        {/each}
        {#if filteredMembers.length === 0}<div class="empty-state">Nenhum membro encontrado.</div>{/if}
      </div>
    </TabsContent>

    <TabsContent value="modlog">
      <div class="events-list">
        {#each warnings as event}
          <div class="event-card">
            <span class="event-icon">{getActionIcon(event.action)}</span>
            <div class="event-info">
              <span class="event-action">{event.action}</span>
              <span class="event-target">{event.target_name || event.target_id}</span>
              <span class="event-details">{event.details}</span>
            </div>
            <span class="event-time">{new Date(event.created_at).toLocaleString('pt-BR')}</span>
          </div>
        {/each}
        {#if warnings.length === 0}<div class="empty-state">Nenhum evento de moderacao.</div>{/if}
      </div>
    </TabsContent>

    <TabsContent value="audit">
      <div class="events-list">
        {#each auditLog as event}
          <div class="event-card">
            <span class="event-icon">📋</span>
            <div class="event-info">
              <span class="event-action">{event.action}</span>
              <span class="event-target">{event.target_name || event.target_id}</span>
              <span class="event-details">{event.details}</span>
            </div>
            <span class="event-time">{new Date(event.created_at).toLocaleString('pt-BR')}</span>
          </div>
        {/each}
        {#if auditLog.length === 0}<div class="empty-state">Nenhum evento de auditoria.</div>{/if}
      </div>
    </TabsContent>

    <TabsContent value="welcome">
      <div class="welcome-config">
        <label><Checkbox bind:checked={welcomeConfig.enabled} label="Ativar Mensagens de Boas-vindas" /></label>
        <div class="field"><label>Canal de Texto ID</label><input type="text" bind:value={welcomeConfig.channel_id} /></div>
        <div class="field"><label>Mensagem de Boas-vindas</label><textarea bind:value={welcomeConfig.welcome_message}></textarea></div>
        <div class="field"><label>Mensagem de Saída</label><textarea bind:value={welcomeConfig.goodbye_message}></textarea></div>
        <Button onclick={saveWelcomeConfig}>Salvar Configuração</Button>
      </div>
    </TabsContent>
  </Tabs>
</div>

<style>
    .admin-view {
        padding: 1.5rem;
        max-width: 1200px;
        margin: 0 auto;
    }

    .admin-header {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        margin-bottom: 1.5rem;
        color: var(--foreground);
    }

    .admin-header h2 {
        margin: 0;
        font-size: 1.5rem;
    }
    
    .search-bar {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.5rem 1rem;
        background: var(--glass-surface);
        border: 1px solid var(--border);
        border-radius: 8px;
        margin-bottom: 1rem;
        color: var(--muted-foreground);
    }

    .search-bar input {
        flex: 1;
        border: none;
        background: transparent;
        color: var(--foreground);
        font-size: 0.9rem;
        outline: none;
    }

    .members-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: 1rem;
    }

    .member-card {
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 1rem;
        background: var(--glass-surface);
        border: 1px solid var(--border);
        border-radius: 12px;
        transition: all 0.2s;
    }

    .member-card:hover {
        border-color: var(--muted);
    }

    .member-avatar {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        background: var(--muted);
        color: var(--primary);
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: 600;
        font-size: 1.1rem;
    }

    .member-info {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .member-name {
        font-weight: 600;
        color: var(--foreground);
    }

    .member-id {
        font-size: 0.75rem;
        color: var(--muted-foreground);
        font-family: monospace;
    }

    .member-notes {
        font-size: 0.8rem;
        color: var(--muted-foreground);
    }

    .member-meta {
        text-align: right;
    }

    .member-updated {
        font-size: 0.75rem;
        color: var(--muted-foreground);
    }

    .events-list {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .event-card {
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 0.75rem 1rem;
        background: var(--glass-surface);
        border: 1px solid var(--border);
        border-radius: 8px;
    }

    .event-icon {
        font-size: 1.2rem;
    }

    .event-info {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .event-action {
        font-weight: 600;
        color: var(--foreground);
        text-transform: capitalize;
    }

    .event-target, .event-actor {
        font-size: 0.85rem;
        color: var(--muted-foreground);
    }

    .event-details {
        font-size: 0.8rem;
        color: var(--muted-foreground);
    }

    .event-time {
        font-size: 0.75rem;
        color: var(--muted-foreground);
        white-space: nowrap;
    }

    .welcome-config {
        max-width: 600px;
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .config-field {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .config-field label {
        font-weight: 600;
        color: var(--foreground);
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .config-field input[type="text"],
    .config-field textarea {
        padding: 0.5rem 1rem;
        background: var(--glass-surface);
        border: 1px solid var(--border);
        border-radius: 8px;
        color: var(--foreground);
        font-size: 0.9rem;
        resize: vertical;
    }

    .config-field input[type="text"]:focus,
    .config-field textarea:focus {
        outline: none;
        border-color: var(--primary);
    }

    .config-hint {
        font-size: 0.75rem;
        color: var(--muted-foreground);
    }

    .save-btn {
        padding: 0.75rem 1.5rem;
        background: var(--primary);
        color: var(--background);
        border: none;
        border-radius: 8px;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s;
        align-self: flex-start;
    }

    .save-btn:hover {
        opacity: 0.9;
        transform: translateY(-1px);
    }

    .empty-state {
        text-align: center;
        padding: 2rem;
        color: var(--muted-foreground);
    }
</style>
