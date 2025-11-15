import { testLogoWithTitleText } from '@/__tests__/shared/SharedTests.ts';
import { mountVuetify } from '@/vuetify-setup.ts';
import MainMenuPage from '@/pages/main_menu_page/MainMenuPage.vue';
import { describe, expect, it } from 'vitest';
import AvatarWithFullName from '@/pages/main_menu_page/AvatarWithFullName.vue';

describe('Main Menu', () => {
  it('checks if a name of a user is visible', async () => {
    const page = mountVuetify(MainMenuPage);
    const avatarWithFullNameComponent = page.findComponent(AvatarWithFullName);
    const fullNameString = avatarWithFullNameComponent.get('span.text-body-1');
    const avatar = avatarWithFullNameComponent.get('.v-avatar');

    expect(fullNameString.text()).toBe('Full Name');
    expect(avatar.text()).toBe('FN');
  });

  testLogoWithTitleText(() => mountVuetify(MainMenuPage));

  it("checks if the 'Join a Room' button is visible", () => {
    const page = mountVuetify(MainMenuPage);
    const joinRoomButton = page.get('[data-testid="join-room-btn"]');

    expect(joinRoomButton.isVisible()).toBe(true);
    expect(joinRoomButton.text()).toBe('Join a Room');
    expect(joinRoomButton.classes()).toContain('bg-primary');
  });
  it("checks if the 'Create a Room' button is visible", () => {
    const page = mountVuetify(MainMenuPage);
    const createRoomButton = page.get('[data-testid="create-room-btn"]');

    expect(createRoomButton.isVisible()).toBe(true);
    expect(createRoomButton.text()).toBe('Create a Room');
    expect(createRoomButton.classes()).toContain('v-btn--variant-outlined');
  });
  it("checks if the 'Logout' button is visible", () => {
    const page = mountVuetify(MainMenuPage);
    const logoutButton = page.get('[data-testid="logout-btn"]');

    expect(logoutButton.isVisible()).toBe(true);
    expect(logoutButton.text()).toBe('Logout');
  });
});
