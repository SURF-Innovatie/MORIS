// apps/frontend/orval.config.ts
import { defineConfig } from "orval";

export default defineConfig({
  backendApi: {
    input: "./.orval/openapi.json",
    output: {
      mode: "split",
      target: "./src/api/generated-orval/",
      schemas: "./src/api/generated-orval/model",
      client: "react-query",
      httpClient: "axios",
      override: {
        mutator: {
          path: "./src/api/custom-axios.ts",
          name: "customInstance",
        },
      },
    },
  },
});
