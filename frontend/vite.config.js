import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

export default defineConfig({
  plugins: [svelte()],
  server: {
    proxy: {
      '/ws': {
        target: 'https://localhost:8443',
        ws: true,
        secure: false
      },
      '/register': {
        target: 'https://localhost:8443',
        secure: false,
      },
      '/login': {
        target: 'https://localhost:8443',
        secure: false,
      }
    }
  }
})
