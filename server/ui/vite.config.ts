import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";

export default defineConfig({
  plugins: [react()],
  base: "/",
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  build: {
    outDir: "dist",
    emptyOutDir: true,
  },
  server: {
    proxy: {
      "/api": "http://localhost:8080",
    },
  },
  test: {
    globals: true,
    environment: "jsdom",
    setupFiles: ["./src/test-setup.ts"],
    // jsdom >=26 returns an opaque origin (and a SecurityError on
    // localStorage access) when no URL is set. The persisted zustand
    // store hits localStorage on construction, so tests crash on
    // import unless we hand jsdom a real origin.
    environmentOptions: {
      jsdom: { url: "http://localhost/" },
    },
  },
});
