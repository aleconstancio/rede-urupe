<script lang="ts">
  import { goto } from '$app/navigation';
  import { Card, Button, Badge } from 'bindrunes';

  const pillars = [
    { icon: '🌐', title: 'Mapa do Conhecimento', desc: 'Obsidian notes publicadas como um site interativo com grafo, busca e backlinks.', href: '/knowledge', cta: 'Explorar' },
    { icon: '🤖', title: 'Bot Discord Cognitivo', desc: 'Pipeline multi-LLM com memória episódica, persona dinâmica e triagem social.', href: '/dashboard/maze', cta: 'Ver dashboard' },
    { icon: '📦', title: 'Projetos', desc: 'Vico, Orb, Thoth, Bindrunes — conheça tudo que construímos.', href: '/projects', cta: 'Ver projetos' },
  ];

  const stats = [
    { value: '6', label: 'Projetos' },
    { value: '25+', label: 'Tabelas SQLite' },
    { value: '3', label: 'Providers LLM' },
    { value: '100%', label: 'Open Source' },
  ];
</script>

<svelte:head>
  <title>Labirinto de Dédalo — Hub de Produtos Inteligentes</title>
  <meta name="description" content="Hub central conectando automação jurídica, prática mística, bots cognitivos e mapa de conhecimento." />
</svelte:head>

<!-- Hero -->
<div class="relative overflow-hidden border-b border-border">
  <div class="absolute inset-0 bg-gradient-to-br from-primary/10 via-transparent to-accent/10"></div>
  <div class="relative mx-auto max-w-5xl px-6 py-24 text-center">
    <Badge variant="info" class="mb-6">Hub Central</Badge>
    <h1 class="text-5xl font-bold tracking-tight mb-6">
      Labirinto de Dédalo
    </h1>
    <p class="text-xl text-muted-foreground max-w-2xl mx-auto mb-8">
      Um ecossistema de produtos inteligentes — automação jurídica, prática mística,
      bots cognitivos e um mapa de conhecimento vivo.
    </p>
    <div class="flex gap-4 justify-center">
      <Button onclick={() => goto('/projects')} variant="default" class="px-8">Ver Projetos</Button>
      <Button onclick={() => goto('/knowledge')} variant="outline" class="px-8">Mapa do Conhecimento</Button>
    </div>
  </div>
</div>

<!-- Stats -->
<section class="border-b border-border">
  <div class="mx-auto max-w-5xl px-6 py-12 grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
    {#each stats as s}
      <div>
        <div class="text-3xl font-bold text-primary">{s.value}</div>
        <div class="text-sm text-muted-foreground mt-1">{s.label}</div>
      </div>
    {/each}
  </div>
</section>

<!-- Pillars -->
<section class="mx-auto max-w-5xl px-6 py-20">
  <div class="text-center mb-12">
    <h2 class="text-3xl font-bold mb-3">O Hub</h2>
    <p class="text-muted-foreground">Três pilares conectados em uma única plataforma</p>
  </div>
  <div class="grid gap-6 md:grid-cols-3">
    {#each pillars as p}
      <Card variant="glass" class="!p-8 flex flex-col">
        <div class="text-4xl mb-4">{p.icon}</div>
        <h3 class="text-xl font-semibold mb-2">{p.title}</h3>
        <p class="text-sm text-muted-foreground flex-1 mb-6">{p.desc}</p>
        <Button onclick={() => goto(p.href)} variant="outline" fullWidth>{p.cta}</Button>
      </Card>
    {/each}
  </div>
</section>

<!-- CTA -->
<section class="mx-auto max-w-5xl px-6 py-20 text-center">
  <Card variant="solid" class="!p-12 bg-gradient-to-r from-primary/20 to-accent/20">
    <h2 class="text-3xl font-bold mb-4">Explore o labirinto</h2>
    <p class="text-muted-foreground mb-6">Knowledge map público. Dashboard com login. Tudo conectado.</p>
    <div class="flex gap-4 justify-center">
      <Button onclick={() => goto('/knowledge')} class="px-8">Conhecimento</Button>
      <Button onclick={() => goto('/login')} variant="outline" class="px-8">Dashboard</Button>
    </div>
  </Card>
</section>

<!-- Footer -->
<footer class="border-t border-border py-8">
  <div class="mx-auto max-w-5xl px-6 flex flex-col sm:flex-row items-center justify-between gap-4">
    <div class="flex items-center gap-2">
      <span class="text-xl">🌀</span>
      <span class="font-semibold">Labirinto de Dédalo</span>
    </div>
    <div class="flex gap-6 text-sm text-muted-foreground">
      <a href="/projects" class="hover:text-foreground">Projetos</a>
      <a href="/knowledge" class="hover:text-foreground">Conhecimento</a>
      <a href="https://github.com/aleconstancio" target="_blank" rel="noopener" class="hover:text-foreground">GitHub</a>
    </div>
    <p class="text-xs text-muted-foreground">© 2026 Labirinto de Dédalo</p>
  </div>
</footer>
