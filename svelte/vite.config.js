import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import preprocess from "svelte-preprocess";

export default defineConfig({
  server: {/* 使用していない */ },
  build: {
    outDir: "dist",
    /* 本番ビルド時はfalseにする */
    sourcemap: "inline",
    emptyOutDir: true
  },
  plugins: [
    svelte({
      preprocess: [preprocess()],
    }),
  ],
});
