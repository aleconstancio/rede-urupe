<script>
    import { onMount } from 'svelte';
    import { initSSE, status as statusStore } from './lib/stores';
    import { AppProvider, ThemeToggle, ErrorBoundary } from '@talos/ui';
    import { Home, Activity, Shield, BrainCircuit, Sparkles, Globe, FileText, Send, Compass, Cpu, LogOut } from 'lucide-svelte';
    
    import PublicLandingView from './lib/views/PublicLandingView.svelte';
    import DashboardView from './lib/views/DashboardView.svelte';
    import CMSArticleEditor from './lib/views/CMSArticleEditor.svelte';
    import MiceliumStudio from './lib/views/MiceliumStudio.svelte';
    import SporeOpsView from './lib/views/SporeOpsView.svelte';
    import GuaraGeoView from './lib/views/GuaraGeoView.svelte';
    import JataiOpsView from './lib/views/JataiOpsView.svelte';
    import MemoryView from './lib/views/MemoryView.svelte';
    import AnalyticsView from './lib/views/AnalyticsView.svelte';

    let currentMode = $state('public'); // 'public' | 'admin'
    let activeView = $state('dashboard');

    const viewLabels = {
        dashboard: 'Centro de Comando',
        cms: 'CMS & Publicações Públicas',
        micelium: 'Estúdio Micélium 🍄 (IA & Manifesto)',
        spore: 'Spore Ops 🍄 (Guerrilha & Agitprop)',
        guara: 'Guará 🪶 (Inteligência Geoespacial)',
        jatai: 'Jataí Ops 🐝 (Frota B2B & Transparência)',
        analytics: 'Analítica & Métricas',
        memory: 'Arquivos Cognitivos'
    };

    const navItems = [
        { id: 'dashboard', label: 'Centro de Comando', icon: Home, description: 'Painel principal e métricas' },
        { id: 'cms', label: 'CMS Publicações', icon: FileText, description: 'Artigos e site público' },
        { id: 'micelium', label: 'Micélium 🍄', icon: Sparkles, description: 'Manifesto, IA & Memórias' },
        { id: 'spore', label: 'Spore Ops 🍄', icon: Send, description: 'Guerrilha e agitprop' },
        { id: 'guara', label: 'Guará Geo 🪶', icon: Compass, description: 'Satélites e meio ambiente' },
        { id: 'jatai', label: 'Jataí Ops 🐝', icon: Cpu, description: 'Frota B2B e faturamento' },
        { id: 'analytics', label: 'Analítica', icon: Activity, description: 'Gráficos e tendências' }
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
        {#if currentMode === 'public'}
            <!-- Public Site Top Nav Overlay -->
            <div class="fixed top-4 right-4 z-50 flex items-center gap-3 bg-[var(--card)] p-2 rounded-full border border-[var(--border)] shadow-lg">
                <ThemeToggle />
                <button class="flex items-center gap-2 px-4 py-2 bg-[var(--primary)] text-[var(--background)] font-semibold rounded-full text-xs hover:brightness-110 transition-all cursor-pointer"
                    onclick={() => currentMode = 'admin'}>
                    <Shield size={14} />
                    <span>Acessar Console Admin</span>
                </button>
            </div>
            <PublicLandingView />
        {:else}
            <!-- Admin Dashboard Layout -->
            <div class="flex min-h-screen">
                <aside class="w-64 border-r flex flex-col justify-between p-4"
                    style="background: var(--card); border-color: var(--border);">
                    <div>
                        <div class="flex items-center justify-between px-2 pb-4 mb-4 border-b" style="border-color: var(--border)">
                            <div class="flex items-center gap-2">
                                <Sparkles size={22} style="color: var(--primary)" />
                                <span class="font-serif text-lg font-bold" style="color: var(--foreground)">Núcleo Urupê</span>
                            </div>
                            <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-[var(--muted)] text-[var(--primary)]">v1.2</span>
                        </div>
                        <nav class="flex flex-col gap-1">
                            {#each navItems as item}
                                {@const Icon = item.icon}
                                <button class="flex items-center gap-3 w-full px-3 py-2 rounded-[--radius] transition-all duration-[--duration-snappy] text-left cursor-pointer"
                                    style="background: {activeView === item.id ? 'var(--muted)' : 'transparent'}; color: {activeView === item.id ? 'var(--foreground)' : 'var(--muted-foreground)'};"
                                    onclick={() => activeView = item.id}>
                                    <Icon size={16} />
                                    <div><span class="text-xs font-semibold">{item.label}</span></div>
                                </button>
                            {/each}
                        </nav>
                    </div>
                    <div class="pt-4 border-t flex flex-col gap-3" style="border-color: var(--border)">
                        <button class="flex items-center gap-2 w-full px-3 py-2 rounded-[--radius] border border-[var(--border)] text-xs text-[var(--muted-foreground)] hover:text-[var(--foreground)] cursor-pointer"
                            onclick={() => currentMode = 'public'}>
                            <Globe size={14} />
                            <span>Ver Site Público</span>
                        </button>
                        <ThemeToggle />
                    </div>
                </aside>

                <div class="flex-1 flex flex-col min-h-screen">
                    <header class="h-[60px] flex items-center justify-between px-8 sticky top-0 z-10"
                        style="background: var(--glass-surface); backdrop-filter: blur(10px); border-bottom: 1px solid var(--border)">
                        <div class="flex items-center gap-4">
                            <div class="flex items-baseline gap-4">
                                <h2 class="font-serif text-lg font-semibold" style="color: var(--foreground)">{viewLabels[activeView] || 'Console Núcleo Urupê'}</h2>
                                {#if statusNote}<span class="text-xs" style="color: var(--muted-foreground)">{statusNote}</span>{/if}
                            </div>
                        </div>
                        <div class="flex items-center gap-4">
                            {#if engineStatus}
                                <div class="flex items-center gap-2 rounded-full px-3 py-1 text-xs font-semibold"
                                    style="border: 1px solid {engineStatus?.Active === true ? 'oklch(from var(--success) l c h / 0.3)' : 'oklch(from var(--destructive) l c h / 0.3)'}; color: {engineStatus?.Active === true ? 'var(--success)' : 'var(--destructive)'}">
                                    <span class="w-2 h-2 rounded-full" style="background: {engineStatus?.Active === true ? 'var(--success)' : 'var(--destructive)'}" />
                                    {engineStatus?.Active === true ? 'Ativo' : 'Inativo'}
                                </div>
                            {/if}
                        </div>
                    </header>
                    <main class="flex-1 relative">
                        <div class="view-wrapper">
                            {#if activeView === 'dashboard'}<DashboardView />
                            {:else if activeView === 'cms'}<CMSArticleEditor />
                            {:else if activeView === 'micelium'}<MiceliumStudio />
                            {:else if activeView === 'spore'}<SporeOpsView />
                            {:else if activeView === 'guara'}<GuaraGeoView />
                            {:else if activeView === 'jatai'}<JataiOpsView />
                            {:else if activeView === 'analytics'}<AnalyticsView />
                            {:else if activeView === 'memory'}<MemoryView />
                            {/if}
                        </div>
                    </main>
                </div>
            </div>
        {/if}
    </ErrorBoundary>
</AppProvider>

<style>
    .view-wrapper { flex: 1; position: relative; }
</style>
