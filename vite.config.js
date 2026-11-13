import react, { reactCompilerPreset } from '@vitejs/plugin-react'
import babel from '@rolldown/plugin-babel'
import { defineConfig } from 'vite'
import tailwindcss from "@tailwindcss/vite";
import electron from "vite-plugin-electron/simple";
import { fileURLToPath } from "node:url";
// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    babel({ presets: [reactCompilerPreset()] }),
      tailwindcss(),
    electron({
      main: { entry: "electron/main.js" },
      preload: { input: "electron/preload.js" },
    }),
  ],
   resolve: {
    alias: { "@": fileURLToPath(new URL("./src", import.meta.url)) },
  },
})
