// vite.config.ts
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import os from 'node:os'

// Получаем первый доступный локальный IPv4 (не localhost)
function getLocalExternalIP() {
  const interfaces = os.networkInterfaces()
  for (const name of Object.keys(interfaces)) {
    for (const iface of interfaces[name] || []) {
      if (iface.family === 'IPv4' && !iface.internal) {
        return iface.address
      }
    }
  }
  return 'localhost'
}

const localIP = getLocalExternalIP()

export default defineConfig({
  plugins: [react()],
  server: {
    host: localIP, // автоматом твой IP
    port: 5173,    // можно поменять, если нужно
  },
})
