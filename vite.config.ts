import { closeSync, mkdirSync, openSync } from 'node:fs';
import { resolve } from 'node:path';
import { defineConfig } from 'vite';

const staticDir = resolve(import.meta.dirname, 'internal/web/static');

export default defineConfig({
  root: 'internal/web/ui',
  build: {
    outDir: '../static',
    emptyOutDir: true,
    sourcemap: true,
  },
  plugins: [
    {
      name: 'webhix-static-gitkeep',
      closeBundle() {
        mkdirSync(staticDir, { recursive: true });
        closeSync(openSync(resolve(staticDir, '.gitkeep'), 'a'));
      },
    },
  ],
});
