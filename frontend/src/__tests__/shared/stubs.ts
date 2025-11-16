import type { Component, ComponentOptions } from 'vue';

export function stub(tag: string, opts?: Partial<ComponentOptions>, template?: string): Component {
  const contents = template || `Stubbed ${tag}`;
  return {
    name: tag,
    template: `<div class="${tag}-stub">${contents}</div>`,
    ...(opts || {}),
  };
}

export const sharedStubs = {
  vDialog: stub(
    'v-dialog',
    {
      props: ['modelValue'],
    },
    '<slot />',
  ),
};
