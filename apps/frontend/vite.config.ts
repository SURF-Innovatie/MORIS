// apps/frontend/vite.config.ts
import { defineConfig } from "vite";
import path from "path";
import react from "@vitejs/plugin-react-swc";

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
      "@api": path.resolve(__dirname, "src/api/generated-orval"),
    },
  },
  server: {
    host: "127.0.0.1",
    port: 3000,
    proxy: {
      "/api": {
        target: "http://127.0.0.1:8080", // Go backend dev server running on host
        changeOrigin: true,
      },
      // Proxy OpenAPI spec access during development for orval to read it
      "/openapi.json": {
        target: "http://127.0.0.1:8080",
        changeOrigin: true,
      },
    },
  },
});
