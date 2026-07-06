import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import path from 'node:path'

// Vitest runs in a jsdom environment so component + DOM-dependent logic (the
// persisted auth store, invite links) can exercise real browser globals. The
// production build (vite.config.ts) is intentionally kept separate and knows
// nothing about tests.
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  // React 19 uses the automatic JSX runtime, so components under test never
  // import React explicitly. Pin it here so esbuild's transform matches.
  esbuild: {
    jsx: 'automatic',
    jsxImportSource: 'react',
  },
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./src/test/setup.ts'],
    css: false,
    include: ['src/**/*.{test,spec}.{ts,tsx}'],
  },
})
