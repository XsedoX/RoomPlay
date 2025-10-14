<script setup lang="ts">
import { shallowRef } from 'vue';
import MusicDataListElement from '@/shared/music_data_list_element/MusicDataListElement.vue';
import type ISearchMusicPopupBaseProps from '@/shared/search_popup_base/ISearchMusicPopupBaseProps.ts';
import type IGuidEvent from '@/shared/IGuidEvent.ts'
import type { TSearchMusicPopupBaseEmits } from '@/shared/search_popup_base/TSearchMusicPopupBaseEmits.ts';

const popup = shallowRef(false)
const props = defineProps<ISearchMusicPopupBaseProps>();
const emit = defineEmits<TSearchMusicPopupBaseEmits>()
function musicChosen(event: IGuidEvent) {
  emit('on-music-choice', event)
  popup.value = false;
}
</script>

<template>
  <v-dialog activator="parent"
            v-model="popup">
    <template v-slot:default>
      <v-row no-gutters justify="center">
        <v-col cols="12" sm="10" md="8" lg="6" xl="4">
          <v-card rounded="xl" color="surface">
            <v-text-field
              :label="props.searchBoxPlaceholder"
              prepend-inner-icon="search"
              append-inner-icon="close"
              hide-details
              @click:append-inner="popup = false"
              single-line>
            </v-text-field>
            <v-sheet color="surface-container" class="music-list-popup-height overflow-y-auto hide-scrollbar">
              <MusicDataListElement v-for="(song, index) in props.musicList"
                                    :key="song.id.toString()"
                                    :class="{
                                      'rounded-b-xl': index === props.musicList.length - 1,
                                    }"
                                    :musicDataListDto="song"
                                    @on-music-choice="musicChosen">
              </MusicDataListElement>
            </v-sheet>
          </v-card>
        </v-col>
      </v-row>
    </template>
  </v-dialog>
</template>

<style scoped>
@import '@/assets/shared.css';
</style>
