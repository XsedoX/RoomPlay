import { expect, describe, it, vi } from 'vitest';
import { flushPromises, mount } from '@vue/test-utils';
import type { StoreGeneric } from 'pinia';
import { mountVuetify } from '@/__tests__/shared/setup_vuetify_tests.ts';
import RoomPage from '@/pages/room_page/RoomPage.vue';
import { useRoomStore } from '@/stores/room_store.ts';
import { faker } from '@faker-js/faker/locale/ar';
import { TUserRole } from '@/infrastructure/user/TUserRole.ts';
import { useUserStore } from '@/stores/user_store.ts';
import { useRouter } from 'vue-router';
import { Guid } from '@/shared/Guid';
import { TSongState } from '@/infrastructure/room/TSongState';
import { TVoteStatus } from '@/infrastructure/room/TVoteStatus';
import { sharedStubs } from '@/__tests__/shared/stubs.ts';

const factory = (
  options?: Parameters<typeof mount>[1],
  piniaStubs?: boolean | string[] | ((actionName: string, store: StoreGeneric) => boolean),
) => mountVuetify(RoomPage, options, piniaStubs);

vi.mock('vue-router', async (importOriginal) => {
  const actual = await importOriginal();
  return Object.assign({}, actual, {
    useRouter: vi.fn(() => ({
      replace: vi.fn(),
      push: vi.fn(),
    })),
  });
});
vi.mock('@vueuse/core', async (importOriginal) => {
  const actual = await importOriginal();
  return Object.assign({}, actual, {
    useIntervalFn: vi.fn(() => ({
      pause: vi.fn(),
      resume: vi.fn(),
      isActive: { value: false },
    })),
  });
});
vi.mock('@vueuse/integrations/useQRCode', () => ({
  useQRCode: vi.fn(() => ({ value: 'mock-qr-code-data' })),
}));
const prepareStores = () => {
  const roomStore = useRoomStore();
  const userStore = useUserStore();
  userStore.user = { name: faker.person.firstName(), surname: faker.person.lastName() };
  roomStore.room = {
    name: faker.word.sample({ length: { min: 5, max: 30 } }),
    boostData: null,
    qrCode: faker.string.uuid(),
    userRole: TUserRole.host,
  };
  roomStore.songs = [];
  roomStore.playingSong = null;
  vi.spyOn(roomStore, 'getRoom').mockResolvedValue(false);
  return { roomStore: roomStore, userStore: userStore };
};

describe('Room Page', () => {
  it('checks if the room renders correctly', async () => {
    const wrapper = factory();
    const { roomStore, userStore } = prepareStores();
    roomStore.playingSong = {
      title: faker.word.sample({ length: { min: 5, max: 30 } }),
      author: faker.person.firstName(),
      lengthSeconds: 120,
      startedAtUtc: new Date(),
    };
    roomStore.songs = [
      {
        id: Guid.generate(),
        title: faker.music.songName(),
        author: faker.music.artist(),
        addedBy: faker.person.fullName() + ' ' + faker.person.lastName(),
        votes: 0,
        albumCoverUrl: faker.image.url(),
        state: TSongState.enqueued,
        voteStatus: TVoteStatus.notVoted,
      },
      {
        id: Guid.generate(),
        title: faker.music.songName(),
        author: faker.music.artist(),
        addedBy: faker.person.fullName() + ' ' + faker.person.lastName(),
        votes: 2,
        albumCoverUrl: faker.image.url(),
        state: TSongState.enqueued,
        voteStatus: TVoteStatus.upvoted,
      },
      {
        id: Guid.generate(),
        title: faker.music.songName(),
        author: faker.music.artist(),
        addedBy: faker.person.fullName() + ' ' + faker.person.lastName(),
        votes: 2,
        albumCoverUrl: faker.image.url(),
        state: TSongState.enqueued,
        voteStatus: TVoteStatus.downvoted,
      },
    ];

    await flushPromises();

    expect(wrapper.get('[data-testid="page-title"]').text()).toBe(roomStore.room!.name);
    expect(wrapper.get('[data-testid="user-initials"]').text()).toBe(
      userStore.user!.name.charAt(0) + userStore.user!.surname.charAt(0),
    );
    expect(wrapper.get('[data-testid="leave-room-btn"').text()).toBe('Leave');
    expect(wrapper.get('[data-testid="search-song-text-field"').isVisible()).toBeTruthy();
    const playingSong = wrapper.get('[data-testid="playing-song"]');
    expect(playingSong.isVisible()).toBeTruthy();
    expect(playingSong.text()).toContain(roomStore.playingSong!.title);
    expect(playingSong.text()).toContain(roomStore.playingSong!.author);
    const songs = wrapper.findAll('[data-testid="song-list-element"]');
    expect(songs.length).toBe(roomStore.songs!.length);
    expect(songs[0]!.text()).toContain(roomStore.songs[0]!.title);
    expect(songs[0]!.text()).toContain(roomStore.songs[0]!.author);
    expect(songs[0]!.text()).toContain(roomStore.songs[0]!.votes.toString());
    expect(songs[1]!.text()).toContain(roomStore.songs[1]!.title);
    expect(songs[1]!.text()).toContain(roomStore.songs[1]!.author);
    expect(songs[1]!.text()).toContain(roomStore.songs[1]!.votes.toString());
    expect(songs[2]!.text()).toContain(roomStore.songs[2]!.title);
    expect(songs[2]!.text()).toContain(roomStore.songs[2]!.author);
    expect(songs[2]!.text()).toContain(roomStore.songs[2]!.votes.toString());
  });
  it("checks if the 'Leave Room' button calls leave function", async () => {
    const router = useRouter();
    const wrapper = factory();
    const { roomStore } = prepareStores();
    await flushPromises();
    const leaveRoomButton = wrapper.get('[data-testid="leave-room-btn"');
    const leaveRoomSpy = vi.spyOn(roomStore, 'leaveRoom').mockResolvedValue();

    await leaveRoomButton.trigger('click');
    await flushPromises();

    expect(leaveRoomSpy).toHaveBeenCalledOnce();
    expect(router.replace).not.toHaveBeenCalled();
  });
  it("checks if 'logout' menu item calls logout function", async () => {
    const wrapper = factory({
      global: {
        stubs: {
          VMenu: sharedStubs.vMenu,
        },
      },
    });
    const { userStore } = prepareStores();
    const logoutSpy = vi.spyOn(userStore, 'logout').mockResolvedValue();
    const settingsMenuButton = wrapper.get('[data-testid="settings-menu-btn"]');
    await settingsMenuButton.trigger('click');
    await flushPromises();

    const logoutMenuItem = wrapper.get('[data-testid="logout-menu-item"]');
    await logoutMenuItem.trigger('click');
    await flushPromises();

    expect(logoutSpy).toHaveBeenCalledOnce();
  });
  it("checks if 'qr code' menu item makes qr code popup visible", async () => {
    const wrapper = factory({
      global: {
        stubs: {
          VMenu: sharedStubs.vMenu,
          VDialog: sharedStubs.vDialog,
          VImg: true,
        },
      },
    });
    prepareStores();
    await flushPromises();
    const settingsMenuButton = wrapper.get('[data-testid="settings-menu-btn"]');
    await settingsMenuButton.trigger('click');
    await flushPromises();

    const qrCodeMenuItem = wrapper.get('[data-testid="qr-code-menu-item"]');
    await qrCodeMenuItem.trigger('click');
    await flushPromises();

    expect(wrapper.get('[data-testid="qr-code-img"]').isVisible()).toBeTruthy();
  });
});
