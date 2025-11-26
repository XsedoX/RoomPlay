import { describe, it, expect } from 'vitest';
import { testLogoWithTitleText } from '@/__tests__/shared/shared_tests.ts';
import LoginPage from '@/pages/login_page/LoginPage.vue';
import { mountVuetify } from '@/__tests__/shared/setup_vuetify_tests.ts';

describe('Login menu', () => {
  it('checks if login with google button is visible', () => {
    const page = mountVuetify(LoginPage);
    const googleButton = page.get('[data-testid="login-with-google-btn"]');

    expect(googleButton.isVisible()).toBe(true);
    expect(googleButton.text()).toBe('Continue with Google');
  });
  testLogoWithTitleText(() => mountVuetify(LoginPage));
});
