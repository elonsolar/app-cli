import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

import setupConfig from './config'


// https://vite.dev/config/
export default defineConfig({

  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
      __ROOT__: path.resolve(__dirname, '')
    },
    extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.vue'] // 自动匹配文件后缀名
  },
   server: {
      proxy: {
        '/api': {
          target: 'http://127.0.0.1:8080',  // 目标后端服务器
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, ''),
        }
      }
    }
})
