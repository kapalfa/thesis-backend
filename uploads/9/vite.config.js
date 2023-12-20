import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import fs from 'fs'
// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  esbuild: {
    loader: 'jsx',
    include: [
      "src/*.jsx",
      "src/**/*.jsx",
    ],
    exclude: /node_modules/,
  },
 server: {
   https: {
     key : fs.readFileSync('/home/andromachi/Documents/go/backend/localhost-key.pem'),
      cert: fs.readFileSync('/home/andromachi/Documents/go/backend/localhost.pem'),
    }
  },
})
