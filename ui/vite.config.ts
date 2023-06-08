import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import dynamicImport from "vite-plugin-dynamic-import";
import { fileURLToPath, URL } from "url";

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      "@constants": fileURLToPath(
        new URL("./src/core/constants", import.meta.url)
      ),
      "@router": fileURLToPath(new URL("./src/core/routers", import.meta.url)),
      "@components": fileURLToPath(
        new URL("./src/shared/components", import.meta.url)
      ),
      "@guards": fileURLToPath(new URL("./src/core/guards", import.meta.url)),
      "@pages": fileURLToPath(new URL("./src/pages", import.meta.url)),
      "@helpers": fileURLToPath(new URL("./src/core/helpers", import.meta.url)),
      "@interfaces": fileURLToPath(
        new URL("./src/shared/interfaces", import.meta.url)
      ),
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  plugins: [vue(), dynamicImport()],
});
