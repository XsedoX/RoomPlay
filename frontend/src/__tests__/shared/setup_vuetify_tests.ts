import { mount } from '@vue/test-utils';
import { type Component} from 'vue';
import createVuetify from '@/vuetify-setup.ts';
import { createTestingPinia, type TestingOptions } from '@pinia/testing';
import type { StoreGeneric } from 'pinia';


const testingVuetify = createVuetify;
export function mountVuetify(
  componentToRender: Component,
  options?: Parameters<typeof mount>[1],
  piniaStubs?: boolean | string[] | ((actionName: string, store: StoreGeneric) => boolean),
) {
  return mount(componentToRender, {
    ...options,
    attachTo: document.body,
    global: {
      ...options?.global,
      plugins: [...(options?.global?.plugins ?? []), testingVuetify, createTestingPinia({ stubActions: piniaStubs } as unknown as TestingOptions)],
    },
  });
}
