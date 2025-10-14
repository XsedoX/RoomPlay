import { VueWrapper } from '@vue/test-utils';
import { it, expect } from 'vitest';

export function testLogoWithTitleText(mountFn: () => VueWrapper){
  it('checks if logo is visible', () => {
    const page = mountFn();
    const logo = page.get('[data-testid="logo"]');

    expect(logo.isVisible()).toBe(true);
  });
  it("checks if application name is visible", ()=>{
    const page = mountFn();
    const textOnPage = page.get('.text-h3').text();

    expect(textOnPage).toBe('RoomPlay2');
  });
  it("checks if subtitle is visible", ()=>{
    const page = mountFn();
    const subtitle = page.get('.text-subtitle-1').text();

    expect(subtitle).toBe('The playlist that creates itself.');
  });
}
