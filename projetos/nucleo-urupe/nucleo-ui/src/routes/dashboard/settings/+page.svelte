<script lang="ts">
  import { Card, Button, ThemeToggle } from 'bindrunes';
  import { getAuth } from '$lib/api';
  import { supabase } from '$lib/api';

  const auth = getAuth();
  const user = $derived(auth.user);

  async function handleLogout() {
    await supabase.auth.signOut();
    auth.logout();
  }
</script>

<svelte:head>
  <title>Configurações — Labirinto de Dédalo</title>
</svelte:head>

<div class="space-y-6">
  <div>
    <h1 class="text-2xl font-bold">Configurações</h1>
    <p class="text-muted-foreground">Conta e preferências</p>
  </div>

  <!-- Profile -->
  <Card variant="glass" class="!p-6">
    <h2 class="text-lg font-semibold mb-4">Perfil</h2>
    <div class="space-y-3">
      <div>
        <div class="text-sm text-muted-foreground">Email</div>
        <div class="font-medium">{user?.email ?? '—'}</div>
      </div>
      <div>
        <div class="text-sm text-muted-foreground">ID</div>
        <div class="font-mono text-xs text-muted-foreground">{user?.id ?? '—'}</div>
      </div>
    </div>
  </Card>

  <!-- Appearance -->
  <Card variant="glass" class="!p-6">
    <h2 class="text-lg font-semibold mb-4">Aparência</h2>
    <div class="flex items-center justify-between">
      <div>
        <div class="font-medium">Tema</div>
        <div class="text-sm text-muted-foreground">Alternar entre claro e escuro</div>
      </div>
      <ThemeToggle variant="outline" />
    </div>
  </Card>

  <!-- Notifications -->
  <Card variant="glass" class="!p-6">
    <h2 class="text-lg font-semibold mb-4">Notificações</h2>
    <p class="text-sm text-muted-foreground">Preferências de notificação serão implementadas em breve.</p>
  </Card>

  <!-- Danger Zone -->
  <Card variant="glass" class="!p-6 border border-destructive/20">
    <h2 class="text-lg font-semibold mb-4 text-destructive">Sair</h2>
    <p class="text-sm text-muted-foreground mb-4">Encerrar sessão neste dispositivo.</p>
    <Button variant="destructive" onclick={handleLogout}>Sair da conta</Button>
  </Card>
</div>
