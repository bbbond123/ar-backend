import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // 代理 /api 开头的请求到nginx转发的后端
      "/api": {
        target: "https://www.ifoodme.com/",
        changeOrigin: true,
        secure: true,
      },
    },
  },
});
