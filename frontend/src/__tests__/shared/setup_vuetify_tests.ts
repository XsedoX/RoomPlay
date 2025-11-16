import { createPinia, type Pinia } from 'pinia';
import { mount } from '@vue/test-utils';
import { type Component } from 'vue';
import createVuetify from '@/vuetify-setup.ts';

const testingVuetify = createVuetify
export function mountVuetify(componentToRender: Component,
                             options?: Parameters<typeof mount>[1],
                             storeSetup?: (pinia: Pinia) => void) {
  const pinia = createPinia();

  // Run store setup BEFORE mounting the component
  // This ensures the store is properly configured when the component initializes
  if (storeSetup) {
    storeSetup(pinia);
  }

  return mount(componentToRender, {
    ...options,
    attachTo: document.body,
    global: {
      ...options?.global,
      plugins: [
        ...(options?.global?.plugins ?? []),
        testingVuetify,
        pinia
      ],
    },
  });
}
