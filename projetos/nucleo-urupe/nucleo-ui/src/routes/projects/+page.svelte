<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Badge, Button } from 'bindrunes';

  const apiBase = import.meta.env.VITE_MAZE_API_URL || 'http://localhost:9393';

  interface Project {
    id: number;
    name: string;
    slug: string;
    description: string;
    icon: string;
    url: string;
    github_url: string;
    category: string;
    status: string;
    is_featured: boolean;
    tags: string[];
  }

  let projects = $state<Project[]>([]);
  let loading = $state(true);

  const categoryLabels: Record<string, string> = {
    product: 'Produto',
    library: 'Biblioteca',
    framework: 'Framework',
    tool: 'Ferramenta',
  };

  type BadgeVariant = 'default' | 'primary' | 'secondary' | 'success' | 'warning' | 'destructive' | 'info' | 'outline';

  function categoryVariant(category: string): BadgeVariant {
    const map: Record<string, BadgeVariant> = {
      product: 'primary',
      library: 'info',
      framework: 'secondary',
      tool: 'warning',
    };
    return map[category] ?? 'secondary';
  }

  onMount(async () => {
    try {
      const res = await fetch(`${apiBase}/api/projects`);
      const data = await res.json();
      projects = data.projects ?? [];
    } catch (e) {
      console.error('Failed to load projects:', e);
    } finally {
      loading = false;
    }
  });

  let featured = $derived(projects.filter((p) => p.is_featured));
  let others = $derived(projects.filter((p) => !p.is_featured));
</script>

<svelte:head>
  <title>Projetos — Labirinto de Dédalo</title>
  <meta name="description" content="Conheça os projetos do ecossistema Labirinto de Dédalo." />
</svelte:head>

<div class="mx-auto max-w-5xl px-6 py-16">
  <div class="mb-12">
    <h1 class="text-4xl font-bold mb-3">Projetos</h1>
    <p class="text-muted-foreground text-lg">
      Tudo que construímos — automação, prática mística, IA e infraestrutura.
    </p>
  </div>

  {#if loading}
    <div class="text-center py-20 text-muted-foreground">Carregando projetos...</div>
  {:else}
    <!-- Featured -->
    {#if featured.length > 0}
      <section class="mb-16">
        <h2 class="text-xl font-semibold mb-6 text-muted-foreground">Em destaque</h2>
        <div class="grid gap-6 md:grid-cols-3">
          {#each featured as p}
            <Card variant="glass" class="!p-6 flex flex-col">
              <div class="flex items-center gap-3 mb-3">
                <span class="text-3xl">{p.icon}</span>
                <div>
                  <h3 class="font-semibold">{p.name}</h3>
                  <Badge variant={categoryVariant(p.category)} class="mt-1">
                    {categoryLabels[p.category] || p.category}
                  </Badge>
                </div>
              </div>
              <p class="text-sm text-muted-foreground flex-1 mb-4">{p.description}</p>
              <div class="flex gap-2">
                {#if p.github_url}
                  <Button onclick={() => window.open(p.github_url, '_blank')} variant="outline" size="sm">GitHub</Button>
                {/if}
                {#if p.url}
                  <Button onclick={() => window.open(p.url, '_blank')} variant="primary" size="sm">Site</Button>
                {/if}
              </div>
            </Card>
          {/each}
        </div>
      </section>
    {/if}

    <!-- Others -->
    {#if others.length > 0}
      <section>
        <h2 class="text-xl font-semibold mb-6 text-muted-foreground">Todos os projetos</h2>
        <div class="grid gap-4 md:grid-cols-2">
          {#each others as p}
            <Card variant="glass" class="!p-5 flex items-start gap-4">
              <span class="text-2xl">{p.icon}</span>
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-1">
                  <h3 class="font-semibold">{p.name}</h3>
                  <Badge variant={categoryVariant(p.category)}>
                    {categoryLabels[p.category] || p.category}
                  </Badge>
                </div>
                <p class="text-sm text-muted-foreground">{p.description}</p>
              </div>
              {#if p.github_url}
                <Button onclick={() => window.open(p.github_url, '_blank')} variant="ghost" size="sm">GitHub</Button>
              {/if}
            </Card>
          {/each}
        </div>
      </section>
    {/if}
  {/if}
</div>
