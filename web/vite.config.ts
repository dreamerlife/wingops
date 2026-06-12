import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

const apiProxyTarget = process.env.VITE_API_PROXY_TARGET ?? 'http://127.0.0.1:8080'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/api': {
        target: apiProxyTarget,
        changeOrigin: true
      }
    }
  },
  test: {
    environment: 'jsdom'
  }
})
