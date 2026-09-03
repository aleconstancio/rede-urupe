<script>
  import { Button, Input, Select, Checkbox, Switch, Dialog, EmptyState, Skeleton, Badge, Alert, Spinner } from '@talos/ui';

    import { status } from '../stores';
    import { toggleMode, mode } from 'mode-watcher';

    let { activeView = 'dashboard', viewLabel = '' } = $props();

    const viewTitles = {
        dashboard: 'Console do Bot',
        analytics: 'Analise',
        persona: 'Persona'
    };

    let title = $derived(viewLabel || viewTitles[activeView] || 'Console do Bot');
</script>

<header class="flex items-center justify-between px-8 h-[70px] sticky top-0 z-10"
    style="background: var(--glass-surface); backdrop-filter: blur(10px); border-bottom: 1px solid var(--border)">
    <div class="flex items-baseline gap-4">
        <h2 class="font-serif text-xl font-semibold" style="color: var(--foreground)">{title}</h2>
        {#if $status?.LastNote}
            <span class="text-sm" style="color: var(--muted-foreground)">{$status.LastNote}</span>
        {/if}
    </div>
    <div class="flex items-center gap-3">
        <Button onclick={toggleMode} class="w-9 h-9 flex items-center justify-center rounded-[--radius] border cursor-pointer"
            style="border-color: var(--border); color: var(--muted-foreground)">
            {#if $mode === 'dark'}
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>
            {:else}
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>
            {/if}
        </Button
    </div>
</header>
