import { testLogoWithTitleText } from '@/__tests__/shared/SharedTests.ts';
import MainMenuPage from '@/pages/main_menu_page/MainMenuPage.vue';
import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import AvatarWithFullName from '@/pages/main_menu_page/AvatarWithFullName.vue';
import { mountVuetify } from '@/__tests__/shared/setup_vuetify_tests.ts';
import { useUserStore } from '@/stores/user_store.ts';
import { RoomService } from '@/infrastructure/room/room_service.ts';
import { setActivePinia } from 'pinia';
import CreateRoomPopup from '@/pages/main_menu_page/CreateRoomPopup.vue';
import { VTextField } from 'vuetify/components';
import { mount } from '@vue/test-utils'
import { sharedStubs } from '@/__tests__/shared/stubs.ts';

// Mock the RoomService module at the top level
// This prevents real API calls during tests
vi.mock('@/infrastructure/room/room_service.ts', () => ({
  RoomService: {
    getUserRoomMembership: vi.fn(),
  },
}));
vi.stubGlobal('visualViewport', {
  offsetLeft: 0,
  offsetTop: 0,
  pageLeft: 0,
  pageTop: 0,
  clientWidth: 1024,
  clientHeight: 768,
  scale: 1,
  width: 1024,
  height: 768,
  addEventListener: vi.fn(),
  removeEventListener: vi.fn(),
});

const factory = (options?: Parameters<typeof mount>[1]) => mountVuetify(MainMenuPage, options, (pinia) => {
  // setActivePinia is needed to use stores outside of components
  // This makes Pinia available for use in this callback
  setActivePinia(pinia);

  const userStore = useUserStore();
  // Set user data before component mounts
  userStore.user = { name: 'Full', surname: 'Name' };
});
describe('Main Menu', () => {
  // Reset all mocks before each test to ensure clean state
  beforeEach(() => {
    vi.mocked(RoomService.getUserRoomMembership).mockResolvedValue(false);
  });
  afterEach(()=>{
    vi.clearAllMocks();
  })

  it('checks if a name of a user is visible', async () => {

    const wrapper = factory();

    const avatarWithFullNameComponent = wrapper.findComponent(AvatarWithFullName);
    const fullNameString = avatarWithFullNameComponent.get('span.text-body-1');
    const avatar = avatarWithFullNameComponent.get('[data-testid="avatar-with-full-name-initials"]');

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
          VDialog: sharedStubs.vDialog
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
});
