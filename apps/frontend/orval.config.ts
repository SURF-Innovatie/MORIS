// apps/frontend/orval.config.ts
import { defineConfig } from 'orval';

export default defineConfig({
  backendApi: {
    input: '../../apps/backend/api/openapi.json',
    output: {
      mode: 'split',
      target: './src/api/generated-orval/',
      schemas: './src/api/generated-orval/model',
      client: 'react-query',
      override: {
        mutator: {
          path: './src/api/custom-axios.ts',
          name: 'customInstance',
        },
      },
    },
  },
});