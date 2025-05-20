import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // 代理 /api 开头的请求到后端
      "/api": {
        // target: "http://localhost:3000",
        target: "https://api.ifoodme.com/",
        changeOrigin: true,
        // 如果后端没有 /api 前缀，可以加上 rewrite
        // rewrite: (path) => path.replace(/^\/api/, ''),
      },
      // // 你也可以加上 /auth 代理
      // "/auth": {
      //   target: "http://localhost:3000",
      //   changeOrigin: true,
      // },
    },
  },
});
