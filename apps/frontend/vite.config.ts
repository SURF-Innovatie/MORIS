// apps/frontend/vite.config.ts
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // Go backend dev server running on host
        changeOrigin: true,
      },
      // Proxy OpenAPI spec access during development for orval to read it
      '/openapi.json': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
});