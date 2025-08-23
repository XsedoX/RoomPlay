import { describe, it, expect } from 'vitest';

import { mount } from '@vue/test-utils';
import LoginPage from '@/loginPage/LoginPage.vue';
import vuetify from '@/vuetify.ts'

describe('Login menu', () => {
  it('checks if logo is visible', () => {
    const wrapper = mount(LoginPage, {
      props:{},
      global: {
        components: { LoginPage },
        plugins: [vuetify],
      },
    });

    const logo = wrapper.find("img[alt='RoomPlay2 logo']");

    expect(logo.exists()).toBe(true);
  });
});
