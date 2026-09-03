import { createClient } from '@supabase/supabase-js';
import { createAuth } from 'bindrunes';

const supabaseUrl = import.meta.env.VITE_SUPABASE_URL;
const supabaseAnonKey = import.meta.env.VITE_SUPABASE_ANON_KEY;

export const supabase = createClient(supabaseUrl, supabaseAnonKey);

let boundAuth: ReturnType<typeof createAuth> | null = null;

export function getAuth() {
  if (!boundAuth) {
    boundAuth = createAuth();
  }
  return boundAuth;
}

export async function initSupabaseAuth() {
  const auth = getAuth();
  const { data: { session } } = await supabase.auth.getSession();

  if (session?.access_token && !auth.isAuthenticated) {
    auth.login(session.access_token, {
      id: session.user.id,
      email: session.user.email ?? '',
      name: session.user.user_metadata?.name,
      roles: [],
      permissions: []
    });
  }

  supabase.auth.onAuthStateChange((_event, session) => {
    if (session?.access_token) {
      auth.login(session.access_token, {
        id: session.user.id,
        email: session.user.email ?? '',
        name: session.user.user_metadata?.name,
        roles: [],
        permissions: []
      });
    } else {
      auth.logout();
    }
  });
}
