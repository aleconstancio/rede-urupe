<script lang="ts">
  import { Card, Button } from 'bindrunes';
</script>

<svelte:head>
  <title>Mapa do Conhecimento — Labirinto de Dédalo</title>
  <meta name="description" content="Mapa de conhecimento interativo construído com Obsidian e Quartz." />
</svelte:head>

<div class="mx-auto max-w-5xl px-6 py-16">
  <div class="text-center mb-12">
    <h1 class="text-4xl font-bold mb-3">Mapa do Conhecimento</h1>
    <p class="text-muted-foreground text-lg max-w-2xl mx-auto">
      Notas do Obsidian publicadas como um site interativo com grafo de conexões,
      busca full-text e backlinks.
    </p>
  </div>

  <Card variant="glass" class="!p-8 text-center">
    <div class="text-5xl mb-4">🗺️</div>
    <h2 class="text-xl font-semibold mb-2">Abrir o Mapa</h2>
    <p class="text-sm text-muted-foreground mb-6">
      Navegue pelo grafo de conhecimento, busque notas e descubra conexões.
    </p>
    <Button onclick={() => window.location.assign('/knowledge/')} class="px-8">
      Entrar no Mapa
    </Button>
  </Card>
</div>
