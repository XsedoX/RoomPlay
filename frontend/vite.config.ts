import { fileURLToPath, URL } from 'node:url';

import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';
import vueDevTools from 'vite-plugin-vue-devtools';
import csp from 'vite-plugin-csp-guard';
import { getCspConfig } from './csp.config';
import { securityHeaders } from './security_headers.config';

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');
  const isDev = mode === 'development';
  return {
    test: {
      globals: true,
      environment: 'jsdom',
      server: {
        deps: {
          inline: ['vuetify'],
        },
      },
    },
    plugins: [vue(), vueJsx(), vueDevTools(), csp(getCspConfig(isDev))],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    server: {
      proxy: {
        '/api/v1': {
          target: 'http://localhost:7654',
        },
      },
      headers: securityHeaders,
      host: env.VITE_APP_HOST,
      port: parseInt(env.VITE_APP_PORT),
    },
  };
});
