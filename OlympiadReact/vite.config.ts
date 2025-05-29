import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  // Added to prod. To local comment
  server: {
    host: '172.16.1.39', // Доступ с других устройств в сети
    port: 5173,      // Новый порт (например, 3000)
    strictPort: true, // Запретить автоматический выбор порта, если указанный занят
  }
})
