<script setup lang="ts">
import PageTitle from '@/shared/page_title/PageTitle.vue';
import { computed, ref, shallowRef } from 'vue';
import SettingsListElement from '@/settings_page/settings_list_element/SettingsListElement.vue';
import SettingsSelect from '@/settings_page/settings_select/SettingsSelect.vue';
import SearchDefaultPlaylistPopup from '@/settings_page/search_default_playlist_popup/SearchDefaultPlaylistPopup.vue';
import type IMusicDataListElementDto from '@/shared/music_data_list_element/IMusicDataListElementDto.ts';
import ChooseHostDeviceList from '@/settings_page/choose_host_device_list/ChooseHostDeviceList.vue';
import UsersList from '@/settings_page/users_list/UsersList.vue';

const selectedCooldown = shallowRef('OFF');
const selectedLifespan = shallowRef<number>(24);
const roomCreationDateUTC = new Date().toISOString();

const chosenPlaylist = ref<IMusicDataListElementDto | null>(null);

const roomExpiryDateLocal = computed(() => {
  const expiry = new Date(roomCreationDateUTC);
  expiry.setHours(expiry.getHours() + selectedLifespan.value);
  const hours = String(expiry.getHours()).padStart(2, '0');
  const minutes = String(expiry.getMinutes()).padStart(2, '0');

  return `${hours}:${minutes} ${expiry.toLocaleDateString()}`;
});
function updateSelectedCooldown(value: string) {
  selectedCooldown.value = value;
}
function updateSelectedLifespan(value: string) {
  selectedLifespan.value = Number.parseInt(value);
}
function choosePlaylist(payload: IMusicDataListElementDto) {
  console.log('Playlist chosen:', payload);
  chosenPlaylist.value = payload;
}
</script>

<template>
  <v-container fluid class="surface d-flex flex-column pa-0 pb-1 ma-0 ga-1">
    <PageTitle title="Settings">
      <template v-slot:top-left-corner>
        <v-btn icon="arrow_back" :to="`/room/${$route.params['id']}`" variant="text"></v-btn>
      </template>
      <template v-slot:top-right-corner>
        <v-btn icon="arrow_back" disabled class="invisible" variant="text"></v-btn>
      </template>
    </PageTitle>
    <v-row justify="center" no-gutters>
      <v-col cols="12" sm="10" md="8">
        <v-list bg-color="transparent">
          <settings-list-element header="BOOST OPTIONS" sub-header="Force Play Cooldown">
            <settings-select
              :value="selectedCooldown"
              :on-update="updateSelectedCooldown"
              :items="['OFF', ...Array.from({ length: 60 }, (_, i) => (i + 1).toString())]"
            >
              <span
                v-if="selectedCooldown !== 'OFF' && selectedCooldown === '1'"
                class="text-primary text-subtitle-1 text-medium-emphasis"
              >
                min
              </span>
              <span
                v-else-if="selectedCooldown !== 'OFF'"
                class="text-primary text-subtitle-1 text-medium-emphasis"
              >
                mins
              </span>
            </settings-select>
          </settings-list-element>
          <settings-list-element
            :header="'ROOM LIFESPAN - ' + roomExpiryDateLocal"
            hint="At most 48 hours since creation."
            sub-header="Auto-delete the room after"
          >
            <settings-select
              :value="selectedLifespan.toString()"
              :on-update="updateSelectedLifespan"
              :items="['4', '8', '16', '24', '48']"
            >
              <span class="text-primary text-subtitle-1 text-medium-emphasis"> hours </span>
            </settings-select>
          </settings-list-element>
          <settings-list-element
            header="DEFAULT PLAYLIST"
            :hint="chosenPlaylist?.subtitle"
            :sub-header="chosenPlaylist?.title ?? 'Choose a playlist.'"
          >
            <template v-slot:image>
              <v-img v-if="chosenPlaylist"
                cover
                max-width="50px"
                min-width="50px"
                rounded="lg"
                aspect-ratio="1/1"
                :src="chosenPlaylist.imageUrl"
              ></v-img>
            </template>
            <template v-slot:default>
              <v-btn variant="text"
                     color="primary"
                     rounded="xl">
                Change
                <search-default-playlist-popup
                  @on-playlist-choice="choosePlaylist"
                ></search-default-playlist-popup>
              </v-btn>
            </template>
          </settings-list-element>
          <settings-list-element header="HOST DEVICE">
            <template v-slot:content>
              <choose-host-device-list></choose-host-device-list>
            </template>
          </settings-list-element>
          <settings-list-element header="USERS">
            <template v-slot:content>
              <users-list></users-list>
            </template>
          </settings-list-element>
        </v-list>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
@import '@/assets/shared.css';
</style>
