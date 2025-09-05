import { describe, it, expect } from 'vitest';

import LoginPage from '@/login_page/LoginPage.vue';
import { mountVuetify } from '@/vuetify-setup.ts'

describe('Login menu', () => {
  it('checks if logo is visible', () => {
    const page = mountVuetify(LoginPage)

    const logo = page.get('[data-testid="logo"]');

    expect(logo.isVisible()).toBe(true);
  });
  it('checks if login with google button is visible', () =>{
    const page = mountVuetify(LoginPage)

    const googleButton = page.get('[data-testid="login-with-google-btn"]');

    expect(googleButton.isVisible()).toBe(true);
  });
  it("checks if application name is visible", ()=>{
    const page = mountVuetify(LoginPage)

    const textOnPage = page.get('.text-h3').text();

    expect(textOnPage).toBe('RoomPlay2');
  });
  it("checks if subtitle is visible", ()=>{
    const page = mountVuetify(LoginPage)

    const subtitle = page.get('.text-subtitle-1').text();

    expect(subtitle).toBe('The playlist that creates itself.');
  });
});
