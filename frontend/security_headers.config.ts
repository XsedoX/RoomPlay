export const securityHeaders: Record<string, string> = {
  'X-XSS-Protection': '0',
  'X-Content-Type-Options': 'nosniff',
  'Referrer-Policy': 'strict-origin-when-cross-origin',
  'Strict-Transport-Security': 'max-age=31536000; includeSubDomains; preload',
  'Cross-Origin-Opener-Policy': 'same-origin',
  'Cross-Origin-Embedder-Policy': 'require-corp',
  'Cross-Origin-Resource-Policy': 'same-site',
  Server: 'ThisIsASecret',
  'X-Powered-By': 'Nope',
  'Permissions-Policy':
    'accelerometer=(), ambient-light-sensor=(), bluetooth=(), camera=(), capture-surface-control=(), compute-pressure=(), display-capture=(), gamepad=(), geolocation=(), gyroscope=(), hid=(), magnetometer=(), microphone=(), midi=(), on-device-speech-recognition=(), payment=(), serial=(), speaker-selection=(), storage-access=(), usb=(), xr-spatial-tracking=()',
};
