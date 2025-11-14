/// <reference types="vite/client" />
interface ImportMetaEnv {
  readonly VITE_APP_HOST?: string;
  readonly VITE_APP_PORT?: string;
  // add other VITE_ variables here as needed
}
interface ImportMeta {
  readonly env: ImportMetaEnv;
}
