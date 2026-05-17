import { closeSync, mkdirSync, openSync, writeFileSync } from 'node:fs';
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
        writeFileSync(
          resolve(staticDir, 'placeholder.txt'),
          'Static UI build output is generated here by Vite.\n',
        );
      },
    },
  ],
});
