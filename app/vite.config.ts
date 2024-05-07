import { defineConfig } from "vite";
import preact from "@preact/preset-vite";
import generouted from "@generouted/react-router/plugin";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    generouted(),
    preact({
      devToolsEnabled: true,
      reactAliasesEnabled: true,
    }),
  ],
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
