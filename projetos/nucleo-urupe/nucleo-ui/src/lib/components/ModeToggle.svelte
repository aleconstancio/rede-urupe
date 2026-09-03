<script>
  import { Button, Input, Select, Checkbox, Switch, Dialog, EmptyState, Skeleton, Badge, Alert, Spinner } from '@talos/ui';

    import { persona } from '../stores';
    import { Zap, ZapOff } from 'lucide-svelte';

    async function toggleMode() {
        if (!$persona || !$persona.ActiveProfile) return;
        
        const nextPassive = !$persona.ActiveProfile.PassiveMode;
        
        try {
            const res = await fetch('/api/mode', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ passive: nextPassive })
            });
            
            if (res.ok) {
                // The SSE 'persona' event will trigger a store update, 
                // but we can update locally for immediate feedback
                $persona.ActiveProfile.PassiveMode = nextPassive;
            }
        } catch (err) {
            console.error("Failed to toggle mode:", err);
        }
    }

    $: isPassive = $persona?.ActiveProfile?.PassiveMode ?? true;
</script>

<div class="mode-toggle-wrapper">
    <Button 
        class="mode-btn" 
        class:active={!isPassive} 
        onclick={toggleMode}
        title={isPassive ? "Ativar Modo Ativo (Bot responderá sozinho)" : "Desativar Modo Ativo (Bot apenas observará)"}
    >
        {#if isPassive}
            <div class="icon-stack">
                <ZapOff size={18} />
            </div>
            <span class="label">PASSIVO</span>
        {:else}
            <div class="icon-stack active-zap">
                <Zap size={18} fill="currentColor" />
            </div>
            <span class="label">MODO ATIVO</span>
        {/if}
    </Button
</div>

<style>
    .mode-toggle-wrapper {
        display: flex;
        align-items: center;
    }

    .mode-btn {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        background: var(--card);
        border: 1px solid var(--border);
        padding: 0.5rem 1rem;
        border-radius: 999px;
        cursor: pointer;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        color: var(--muted-foreground);
    }

    .mode-btn:hover {
        border-color: var(--border);
        background: var(--card);
        transform: translateY(-1px);
    }

    .mode-btn.active {
        background: oklch(from var(--warning) l c h / 0.1);
        border-color: var(--primary);
        color: var(--primary);
        box-shadow: 0 0 20px oklch(from var(--warning) l c h / 0.15);
    }

    .icon-stack {
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .active-zap {
        filter: drop-shadow(0 0 8px var(--primary));
        animation: pulse 2s infinite;
    }

    @keyframes pulse {
        0% { opacity: 1; transform: scale(1); }
        50% { opacity: 0.8; transform: scale(1.1); }
        100% { opacity: 1; transform: scale(1); }
    }

    .label {
        font-size: 0.7rem;
        font-weight: 800;
        letter-spacing: 0.12em;
        text-transform: uppercase;
    }
</style>
