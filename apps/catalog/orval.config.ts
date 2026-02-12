import { defineConfig } from "orval";

export default defineConfig({
  moris: {
    input: {
      target: "../../apps/backend/api/swag-docs/swagger.json",
    },
    output: {
      mode: "tags-split",
      target: "src/api/generated-orval/moris.ts",
      schemas: "src/api/generated-orval/model",
      client: "react-query",
      mock: false,
      baseUrl: "/api",
    },
  },
});
