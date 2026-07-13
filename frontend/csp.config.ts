import csp from 'vite-plugin-csp-guard';

type CspOptions = Parameters<typeof csp>[0];

export function getCspConfig(isDev: boolean): CspOptions {
  return {
    dev: {
      run: true,
      outlierSupport: ['vue'],
    },
    policy: {
      'default-src': ["'self'"],
      'script-src': isDev ? ["'self'", "'unsafe-inline'"] : ["'self'"],
      'style-src': ["'self'", 'https://fonts.googleapis.com'],
      'style-src-elem': ["'self'", "'unsafe-inline'", 'https://fonts.googleapis.com'],
      'img-src': ["'self'", 'data:', 'https://i.ytimg.com', 'https://yt3.ggpht.com'],
      'font-src': ["'self'", 'https://fonts.gstatic.com'],
      'connect-src': ["'self'", 'http://localhost:7654'],
      'base-uri': ["'self'"],
      'form-action': ["'self'"],
      'frame-ancestors': ["'none'"],
      ...(isDev ? {} : { 'upgrade-insecure-requests': [] }),
    },
    build: {
      sri: true,
    },
  };
}
