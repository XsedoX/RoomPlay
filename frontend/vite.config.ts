import { fileURLToPath, URL } from 'node:url';

import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';
import vueDevTools from 'vite-plugin-vue-devtools';
import csp from "vite-plugin-csp-guard";

// https://vite.dev/config/
export default defineConfig(({mode})=>{
  const env = loadEnv(mode, process.cwd(), '');
  return {
    plugins: [vue(),
      vueJsx(),
      vueDevTools(),
      csp({
        dev: {
          run: true,
          outlierSupport: ["vue"]
        },
        policy: {
          "default-src": ["'self'"],
          "script-src": ["'self'"],
          "style-src": ["'self'"],
          "style-src-elem": ["'unsafe-inline'", "https://fonts.googleapis.com"],
          "img-src": ["'self'"],
          "font-src": ["'self'", "https://fonts.googleapis.com", "https://fonts.gstatic.com"],
          "connect-src": ["'self'","http://localhost:7654"],
          "base-uri": ["'self'"],
          "form-action": ["'self'"],
          "frame-ancestors": ["'none'"],
          "upgrade-insecure-requests": []
        },
        build:{
          sri: true
        }
      })],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      }
    },
    server:{
      proxy:{
        '/api/v1':{
          target:'http://localhost:7654',
        },
      },
      headers:{
        "X-XSS-Protection": "0",
        "X-Content-Type-Options": "nosniff",
        "Referrer-Policy": "strict-origin-when-cross-origin",
        "Strict-Transport-Security": "max-age=31536000; includeSubDomains; preload",
        "Cross-Origin-Opener-Policy": "same-origin",
        "Cross-Origin-Embedder-Policy": "require-corp",
        "Cross-Origin-Resource-Policy": "same-site",
        "Server": "ThisIsASecret",
        "X-Powered-By":"Nope",
        "Permissions-Policy": "accelerometer=(), ambient-light-sensor=(), bluetooth=(), camera=(), capture-surface-control=(), compute-pressure=(), display-capture=(), gamepad=(), geolocation=(), gyroscope=(), hid=(), magnetometer=(), microphone=(), midi=(), on-device-speech-recognition=(), payment=(), serial=(), speaker-selection=(), storage-access=(), usb=(), xr-spatial-tracking=()"
      },
      host: env.VITE_HOST,
      port: parseInt(env.VITE_PORT),
    }
  }
});
