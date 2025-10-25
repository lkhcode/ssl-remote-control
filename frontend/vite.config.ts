import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url';

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [vue()],
    resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
    },
    server: {
        proxy: {
            '/api': {
                target: 'http://localhost:8084',
                changeOrigin: true,
                secure: false,
                ws: true,
            }
        }
    },
    build: {
        chunkSizeWarningLimit: 1000,
        rollupOptions: {
        output: {
            manualChunks(id) {
            if (id.includes('node_modules')) {
                // Group dependencies into separate chunks
                return id.toString().split('node_modules/')[1].split('/')[0].toString();
            }
            },
        },
        },
    },
})
