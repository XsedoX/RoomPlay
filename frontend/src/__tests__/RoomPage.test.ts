import { expect, describe, it, vi } from 'vitest';
import { flushPromises, mount } from '@vue/test-utils';
import type { StoreGeneric } from 'pinia';
import { mountVuetify } from '@/__tests__/shared/setup_vuetify_tests.ts';
import RoomPage from '@/pages/room_page/RoomPage.vue';
import { useRoomStore } from '@/stores/room_store.ts'
import { faker } from '@faker-js/faker/locale/ar'
import { TUserRole } from '@/infrastructure/user/TUserRole.ts'
import { useUserStore } from '@/stores/user_store.ts'
import { useRouter } from 'vue-router'

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

const prepareStores = () => {
  const roomStore = useRoomStore();
  const userStore = useUserStore();
  userStore.user = { name: faker.person.firstName(), surname: faker.person.lastName() };
  roomStore.room = {
    name: faker.word.sample({ length: { min: 5, max: 30 } }),
    boostData:null,
    qrCode: faker.string.uuid(),
    userRole: TUserRole.host,
  }
  roomStore.songs = []
  roomStore.playingSong = null;
  vi.spyOn(roomStore, 'getRoom').mockResolvedValue(false);
  return {roomStore: roomStore, userStore: userStore};
}

describe('Room Page', () => {
  it('checks if the room renders correctly', async () => {
    const wrapper = factory();
    const {roomStore, userStore} = prepareStores();

    await flushPromises();

    expect(wrapper.get('[data-testid="page-title"]').text()).toBe(roomStore.room!.name);
    expect(wrapper.get('[data-testid="user-initials"]').text()).toBe(userStore.user!.name.charAt(0) + userStore.user!.surname.charAt(0));
    expect(wrapper.get('[data-testid="leave-room-btn"').text()).toBe('Leave')
    expect(wrapper.get('[data-testid="search-song-text-field"').isVisible()).toBeTruthy()
  });
  it("checks if the 'Leave Room' button calls leave function", async () => {
    const router = useRouter();
    const wrapper = factory();
    const {roomStore} = prepareStores();
    await flushPromises();
    const leaveRoomButton = wrapper.get('[data-testid="leave-room-btn"')
    const leaveRoomSpy = vi.spyOn(roomStore, 'leaveRoom').mockResolvedValue();

    await leaveRoomButton.trigger("click")
    await flushPromises();

    expect(leaveRoomSpy).toHaveBeenCalledOnce();
    expect(router.replace).not.toHaveBeenCalled();
  })
});
