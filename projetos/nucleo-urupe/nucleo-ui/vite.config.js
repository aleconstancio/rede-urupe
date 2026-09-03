import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [tailwindcss(), svelte()],
  resolve: {
    alias: {
      '@talos/ui': 'bindrunes'
    }
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
      '/events': 'http://localhost:8080',
      '/fragments': 'http://localhost:8080',
      '/persona/save': 'http://localhost:8080',
      '/context/clear': 'http://localhost:8080',
    }
  },
  build: {
    outDir: '../internal/presentation/web/static_dist',
    emptyOutDir: true,
  }
})
