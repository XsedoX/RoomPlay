import { fileURLToPath, URL } from 'node:url';
import { mergeConfig, defineConfig, configDefaults } from 'vitest/config';
import viteConfig from './vite.config';

export default defineConfig(({ mode }) =>
  mergeConfig(
    viteConfig({ mode, command: 'serve' }),
    {
      test: {
        environment: 'jsdom',
        exclude: [...configDefaults.exclude, 'e2e/**', 'features/**'],
        root: fileURLToPath(new URL('./', import.meta.url)),
        globals: true,
        server:{
          deps:{
            inline:['vuetify']
          }
        }
      },
    }
  )
);
