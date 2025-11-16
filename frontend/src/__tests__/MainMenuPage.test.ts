import { testLogoWithTitleText } from '@/__tests__/shared/SharedTests.ts';
import MainMenuPage from '@/pages/main_menu_page/MainMenuPage.vue';
import { describe, expect, it, vi } from 'vitest';
import AvatarWithFullName from '@/pages/main_menu_page/AvatarWithFullName.vue';
import { mountVuetify } from '@/__tests__/shared/setup_vuetify_tests.ts';
import { useUserStore } from '@/stores/user_store.ts';
import CreateRoomPopup from '@/pages/main_menu_page/CreateRoomPopup.vue';
import { VTextField } from 'vuetify/components';
import { flushPromises, mount } from '@vue/test-utils';
import { sharedStubs } from '@/__tests__/shared/stubs.ts';
import type { StoreGeneric } from 'pinia';
import { useRoomStore } from '@/stores/room_store.ts';
import JoinRoomPopup from '@/pages/main_menu_page/JoinRoomPopup.vue';
import { PlatformDiscoverer } from '@/infrastructure/utils/platform_discoverer.ts';
import { THostDevice } from '@/pages/settings_page/choose_host_device_list/THostDevice.ts';

const factory = (
  options?: Parameters<typeof mount>[1],
  piniaStubs?: boolean | string[] | ((actionName: string, store: StoreGeneric) => boolean),
) => mountVuetify(MainMenuPage, options, piniaStubs);

async function openRoomPopup(wrapper: ReturnType<typeof factory>) {
  const createRoomButton = wrapper.get('[data-testid="create-room-btn"]');
  await createRoomButton.trigger('click');
  const createRoomPopup = wrapper.getComponent(CreateRoomPopup);
  expect(createRoomPopup.isVisible()).toBe(true);
  const nameField = createRoomPopup.get('[data-testid="create-room-popup-name-input"] input');
  const passwordField = createRoomPopup.get(
    '[data-testid="create-room-popup-password-input"] input',
  );
  const repeatPasswordField = createRoomPopup.get(
    '[data-testid="create-room-popup-repeat-password-input"] input',
  );
  return {
    createRoomPopup: createRoomPopup,
    nameField: nameField,
    passwordField: passwordField,
    repeatPasswordField: repeatPasswordField,
  };
}

vi.mock('@/infrastructure/utils/platform_discoverer.ts', () => ({
  PlatformDiscoverer: {
    getDeviceType: vi.fn(),
  },
}));
describe('Main Menu', () => {
  it('checks if a name of a user is visible', async () => {
    const wrapper = factory();
    const userStore = useUserStore();
    userStore.user = { name: 'Full', surname: 'Name' };
    expect(userStore.user!.name).toBe('Full');

    const avatarWithFullNameComponent = wrapper.findComponent(AvatarWithFullName);
    const fullNameString = avatarWithFullNameComponent.get('span.text-body-1');
    const avatar = avatarWithFullNameComponent.get(
      '[data-testid="avatar-with-full-name-initials"]',
    );
    await flushPromises();

    expect(fullNameString.text()).toBe('Full Name');
    expect(avatar.text()).toBe('FN');
  });

  testLogoWithTitleText(() => {
    return mountVuetify(MainMenuPage);
  });

  it("checks if the 'Join a Room' button is visible", async () => {
    const wrapper = factory();
    const joinRoomButton = wrapper.get('[data-testid="join-room-btn"]');

    expect(joinRoomButton.isVisible()).toBe(true);
    // The button contains child components (popups) with text too
    // We use toContain to check if our expected text is present
    expect(joinRoomButton.text()).toContain('Join a Room');
  });

  it("checks if the 'Create a Room' button is visible", async () => {
    const wrapper = factory();
    const createRoomButton = wrapper.get('[data-testid="create-room-btn"]');

    expect(createRoomButton.isVisible()).toBe(true);
    // The button contains child components (popups) with text too
    // We use toContain to check if our expected text is present
    expect(createRoomButton.text()).toContain('Create a Room');
  });

  it("checks if the 'Logout' button is visible", async () => {
    const wrapper = factory();
    const logoutButton = wrapper.get('[data-testid="logout-btn"]');

    expect(logoutButton.isVisible()).toBe(true);
    expect(logoutButton.text()).toBe('Logout');
  });
  it('checks if create room popup renders correctly', async () => {
    const wrapper = factory({
      global: {
        stubs: {
          VDialog: sharedStubs.vDialog,
        },
      },
    });
    const createRoomButton = wrapper.get('[data-testid="create-room-btn"]');
    await createRoomButton.trigger('click');

    const createRoomPopup = wrapper.getComponent(CreateRoomPopup);
    expect(createRoomPopup.isVisible()).toBe(true);

    const fields = createRoomPopup.findAllComponents(VTextField);
    expect(fields.length).toBe(3);
  });
  it('checks if the CreateRoomPopup calls the createRoom method with the correct parameters', async () => {
    vi.useFakeTimers();
    const wrapper = factory({
      global: {
        stubs: {
          VDialog: sharedStubs.vDialog,
        },
      },
    });
    const roomStore = useRoomStore();
    vi.spyOn(roomStore, 'createRoom').mockResolvedValue(null);

    const { createRoomPopup, nameField, passwordField, repeatPasswordField } =
      await openRoomPopup(wrapper);
    await nameField.setValue('Test Room');
    await passwordField.setValue('TestPassword');
    await repeatPasswordField.setValue('TestPassword');
    await createRoomPopup.get('[data-testid="create-room-popup-btn"]').trigger('click');

    vi.runAllTimers();
    await flushPromises();

    expect(roomStore.createRoom).toHaveBeenCalledExactlyOnceWith({
      roomName: 'Test Room',
      roomPassword: 'TestPassword',
      repeatRoomPassword: 'TestPassword',
    });
    vi.useRealTimers();
  });
  it('checks if the room creation validation works correctly', async () => {
    vi.useFakeTimers();
    const wrapper = factory({
      global: {
        stubs: {
          VDialog: sharedStubs.vDialog,
        },
      },
    });
    const { createRoomPopup, nameField, passwordField, repeatPasswordField } =
      await openRoomPopup(wrapper);

    await nameField.setValue('T');
    await nameField.trigger('blur');
    await passwordField.setValue('TestP');
    await repeatPasswordField.setValue('Test');

    const defaultAlerts = createRoomPopup.findAll('div[role="alert"]');
    expect(defaultAlerts.length).toBe(3);
    expect(defaultAlerts[0]!.text()).toContain('This cannot be changed later');
    expect(defaultAlerts[1]!.text()).toContain('This cannot be changed later');
    expect(defaultAlerts[2]!.text()).toBe('');

    vi.runAllTimers();
    await flushPromises();

    const alerts = createRoomPopup.findAll('div[role="alert"]');
    expect(alerts.length).toBe(3);
    expect(alerts[0]!.text()).toContain('Room name has to have at least 5 characters');
    expect(alerts[1]!.text()).toContain('Password has to have at least 10 characters');
    expect(alerts[2]!.text()).toContain(`Passwords don't match`);

    vi.useRealTimers();
  });
  it("checks if 'Join a Room' popup renders correctly on mobile", async () => {
    vi.mocked(PlatformDiscoverer.getDeviceType).mockReturnValue(THostDevice.Mobile);
    const wrapper = factory({
      global: {
        stubs: {
          VDialog: sharedStubs.vDialog,
        },
      },
    });

    const joinRoomButton = wrapper.get('[data-testid="join-room-btn"]');
    await joinRoomButton.trigger('click');
    const joinRoomPopup = wrapper.getComponent(JoinRoomPopup);

    expect(joinRoomPopup.isVisible()).toBeTruthy();
    expect(joinRoomPopup.get('span[class="text-h5"]').text()).toContain('Join a Room');
    const scanBtnDom = joinRoomPopup
      .findAll('button')
      .find((b) => b.text().includes('Scan a QR Code'));
    expect(scanBtnDom?.isVisible()).toBeTruthy();
    expect(
      joinRoomPopup.get('[data-testid="join-room-popup-name-input"]').isVisible(),
    ).toBeTruthy();
    expect(
      joinRoomPopup.get('[data-testid="join-room-popup-password-input"]').isVisible(),
    ).toBeTruthy();
  });
  it("checks if 'Join a Room' popup renders correctly on desktop", async () => {
    vi.mocked(PlatformDiscoverer.getDeviceType).mockReturnValue(THostDevice.Desktop);
    const wrapper = factory({
      global: {
        stubs: {
          VDialog: sharedStubs.vDialog,
        },
      },
    });

    const joinRoomButton = wrapper.get('[data-testid="join-room-btn"]');
    await joinRoomButton.trigger('click');
    const joinRoomPopup = wrapper.getComponent(JoinRoomPopup);

    expect(joinRoomPopup.isVisible()).toBeTruthy();
    expect(joinRoomPopup.get('span[class="text-h5"]').text()).toContain('Join a Room');
    expect(
      joinRoomPopup.get('[data-testid="join-room-popup-name-input"]').isVisible(),
    ).toBeTruthy();
    expect(
      joinRoomPopup.get('[data-testid="join-room-popup-password-input"]').isVisible(),
    ).toBeTruthy();
  });
});
