import { defineConfig, loadEnv } from "vite";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "@tailwindcss/vite";
import AutoImport from "unplugin-auto-import/vite";
import Component from "unplugin-vue-components/vite";
import Icons from "unplugin-icons/vite";
import { resolve } from "path";
import { viteMockServe } from "vite-plugin-mock";

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd());
  return {
    plugins: [
      vue(),
      tailwindcss(),
      AutoImport({
        imports: ["vue", "vue-router", "vue-i18n", "@vueuse/core"],
        dts: "src/auto-import.d.ts",
        eslintrc: {
          enabled: true,
        },
      }),
      Component({
        dts: "src/components.d.ts",
        dirs: ["src/components"],
        extensions: ["vue"],
      }),
      Icons({
        compiler: "vue3",
        autoInstall: true,
      }),
      viteMockServe({
        mockPath: "./src/mocks",
        enable: env.VITE_USE_MOCK === "true",
      }),
    ],
    resolve: {
      alias: {
        "@": resolve(__dirname, "src"),
      },
    },
    server: {
      proxy: {
        "/api": {
          target: "http://localhost:8080",
          changeOrigin: true,
        },
      },
    },
  };
});
