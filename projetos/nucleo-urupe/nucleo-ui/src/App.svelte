<script>
    import { onMount } from 'svelte';
    import { initSSE, status as statusStore } from './lib/stores';
    import { AppProvider, ThemeToggle, ErrorBoundary } from '@talos/ui';
    import { Home, Activity, Shield, BrainCircuit, Sparkles } from 'lucide-svelte';
    
    import DashboardView from './lib/views/DashboardView.svelte';
    import MemoryView from './lib/views/MemoryView.svelte';
    import PersonaView from './lib/views/PersonaView.svelte';
    import AnalyticsView from './lib/views/AnalyticsView.svelte';
    import AdminView from './lib/views/AdminView.svelte';

    let activeView = $state('dashboard');
    const viewLabels = { dashboard: 'Centro de Comando', analytics: 'Analitica', admin: 'Admin', memory: 'Arquivos Cognitivos', persona: 'Estudio de Persona' };
    const navItems = [
        { id: 'dashboard', label: 'Centro de Comando', icon: Home, description: 'Painel principal e métricas' },
        { id: 'analytics', label: 'Analitica', icon: Activity, description: 'Gráficos e tendências' },
        { id: 'admin', label: 'Admin', icon: Shield, description: 'Configurações do servidor' },
        { id: 'memory', label: 'Arquivos Cognitivos', icon: BrainCircuit, description: 'Cápsulas de memória' },
        { id: 'persona', label: 'Estudio de Persona', icon: Sparkles, description: 'Edição de identidade' }
    ];

    let engineStatus = $state(null);
    let statusNote = $state('');

    onMount(() => {
        const cleanup = initSSE();
        statusStore.subscribe(v => { engineStatus = v; statusNote = v?.LastNote || ''; });
        return cleanup;
    });
</script>

<AppProvider>
    <ErrorBoundary>
        <div class="flex min-h-screen">
            <aside class="w-64 border-r flex flex-col justify-between p-4"
                style="background: var(--card); border-color: var(--border);">
                <div>
                    <div class="flex items-center gap-3 px-2 pb-4 mb-4 border-b" style="border-color: var(--border)">
                        <Sparkles size={24} style="color: var(--primary)" />
                        <span class="font-serif text-xl font-bold" style="color: var(--foreground)">Núcleo Urupê</span>
                    </div>
                    <nav class="flex flex-col gap-1">
                        {#each navItems as item}
                            {@const Icon = item.icon}
                            <button class="flex items-center gap-3 w-full px-3 py-2.5 rounded-[--radius] transition-all duration-[--duration-snappy] text-left cursor-pointer"
                                style="background: {activeView === item.id ? 'var(--muted)' : 'transparent'}; color: {activeView === item.id ? 'var(--foreground)' : 'var(--muted-foreground)'};"
                                onclick={() => activeView = item.id}>
                                <Icon size={18} />
                                <div><span class="text-sm font-semibold">{item.label}</span><p class="text-xs mt-0.5 opacity-70">{item.description}</p></div>
                            </button>
                        {/each}
                    </nav>
                </div>
                <div class="pt-4 border-t" style="border-color: var(--border)">
                    <ThemeToggle />
                </div>
            </aside>

            <div class="flex-1 flex flex-col min-h-screen">
                <header class="h-[70px] flex items-center justify-between px-8 sticky top-0 z-10"
                    style="background: var(--glass-surface); backdrop-filter: blur(10px); border-bottom: 1px solid var(--border)">
                    <div class="flex items-center gap-4">
                        <div class="flex items-baseline gap-4">
                            <h2 class="font-serif text-xl font-semibold" style="color: var(--foreground)">{viewLabels[activeView] || 'Console Núcleo Urupê'}</h2>
                            {#if statusNote}<span class="text-sm" style="color: var(--muted-foreground)">{statusNote}</span>{/if}
                        </div>
                    </div>
                    {#if engineStatus}
                        <div class="flex items-center gap-2 rounded-full px-3 py-1 text-xs font-semibold"
                            style="border: 1px solid {engineStatus?.Active === true ? 'oklch(from var(--success) l c h / 0.3)' : 'oklch(from var(--destructive) l c h / 0.3)'}; color: {engineStatus?.Active === true ? 'var(--success)' : 'var(--destructive)'}">
                            <span class="w-2 h-2 rounded-full" style="background: {engineStatus?.Active === true ? 'var(--success)' : 'var(--destructive)'}" />
                            {engineStatus?.Active === true ? 'Ativo' : 'Inativo'}
                        </div>
                    {/if}
                </header>
                <main class="flex-1 relative">
                    <div class="view-wrapper">
                        {#if activeView === 'dashboard'}<DashboardView />
                        {:else if activeView === 'analytics'}<AnalyticsView />
                        {:else if activeView === 'admin'}<AdminView />
                        {:else if activeView === 'memory'}<MemoryView />
                        {:else if activeView === 'persona'}<PersonaView />
                        {/if}
                    </div>
                </main>
            </div>
        </div>
    </ErrorBoundary>
</AppProvider>

<style>
    .view-wrapper { flex: 1; position: relative; }
</style>
