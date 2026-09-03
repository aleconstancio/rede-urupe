<script lang="ts">
  import type { Snippet } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { AppProvider, ErrorBoundary, Button } from 'bindrunes';
  import '../app.css';

  let { children }: { children: Snippet } = $props();

  let isDashboard = $derived($page.url.pathname.startsWith('/dashboard'));
</script>

<AppProvider theme="labirinto" aesthetic="editorial" density="comfortable">
  <ErrorBoundary>
    {#if !isDashboard}
      <nav class="sticky top-0 z-50 border-b border-border bg-background/80 backdrop-blur-md">
        <div class="mx-auto max-w-6xl px-6 h-14 flex items-center justify-between">
          <a href="/" class="flex items-center gap-2 font-semibold">
            <span class="text-xl">🌀</span>
            <span>Labirinto</span>
          </a>
          <div class="flex items-center gap-6 text-sm">
            <a href="/projects" class="text-muted-foreground hover:text-foreground transition-colors">Projetos</a>
            <a href="/knowledge" class="text-muted-foreground hover:text-foreground transition-colors">Conhecimento</a>
            <a href="/dashboard" class="text-muted-foreground hover:text-foreground transition-colors">Dashboard</a>
            <Button onclick={() => goto('/login')} variant="default" size="sm">Entrar</Button>
          </div>
        </div>
      </nav>
    {/if}
    {@render children()}
  </ErrorBoundary>
</AppProvider>
