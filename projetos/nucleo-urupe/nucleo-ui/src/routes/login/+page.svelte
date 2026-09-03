<script lang="ts">
  import { goto } from '$app/navigation';
  import { Card, Input, Button } from 'bindrunes';
  import { supabase, getAuth } from '$lib/api';

  let email = $state('');
  let password = $state('');
  let error = $state('');
  let loading = $state(false);
  let isSignUp = $state(false);

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    error = '';
    loading = true;

    try {
      const result = isSignUp
        ? await supabase.auth.signUp({ email, password })
        : await supabase.auth.signInWithPassword({ email, password });

      if (result.error) {
        error = result.error.message;
      } else if (result.data.session) {
        const auth = getAuth();
        auth.login(result.data.session.access_token, {
          id: result.data.session.user.id,
          email: result.data.session.user.email ?? '',
          name: result.data.session.user.user_metadata?.name,
          roles: [],
          permissions: []
        });
        goto('/dashboard');
      }
    } catch (err) {
      error = 'Erro ao conectar ao servidor';
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Login — Labirinto de Dédalo</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center">
  <Card variant="glass" class="!p-10 max-w-[420px] w-full">
    <div class="text-center mb-6">
      <h1 class="text-2xl font-bold">Labirinto de Dédalo</h1>
      <p class="text-sm text-muted-foreground mt-1">
        {isSignUp ? 'Crie sua conta' : 'Entre na sua conta'}
      </p>
    </div>

    <form onsubmit={handleSubmit}>
      <Input
        name="email"
        type="email"
        label="E-mail"
        bind:value={email}
        placeholder="seu@email.com"
      />
      <Input
        name="password"
        type="password"
        label="Senha"
        bind:value={password}
        placeholder="••••••••"
        {error}
      />

      <Button type="submit" fullWidth class="mt-2" loading={loading}>
        {isSignUp ? 'Criar conta' : 'Entrar'}
      </Button>
    </form>

    <div class="text-center text-sm mt-4">
      <button
        onclick={() => isSignUp = !isSignUp}
        class="text-primary hover:underline"
      >
        {isSignUp ? 'Já tem conta? Entrar' : 'Não tem conta? Criar'}
      </button>
    </div>
  </Card>
</div>
