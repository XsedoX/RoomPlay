<script setup lang="ts">
import PageTitle from '@/shared/page_title/PageTitle.vue';
import { shallowRef } from 'vue';
import SettingsListElement from '@/settings_page/settings_list_element/SettingsListElement.vue';
import SettingsSelect from '@/settings_page/settings_select/SettingsSelect.vue';

const selectedCooldown = shallowRef('OFF');
const selectedLifespan = shallowRef<number>(24);
function updateSelectedCooldown(value: string) {
  selectedCooldown.value = value;
}
function updateSelectedLifespan(value: number) {
  selectedLifespan.value = value;
}
</script>

<template>
  <v-container fluid
               class="surface d-flex flex-column pa-0 pb-1 ma-0 ga-1">
    <PageTitle title="Settings">
      <template v-slot:top-left-corner>
        <v-btn icon="arrow_back"
               :to="`/room/${$route.params['id']}`"
               variant="text"></v-btn>
      </template>
      <template v-slot:top-right-corner>
        <v-btn icon="arrow_back"
               disabled
               class="invisible"
               variant="text"></v-btn>
      </template>
    </PageTitle>
    <v-row justify="center">
      <v-col cols="12"
             sm="10"
             md="8">
        <v-list bg-color="transparent">
          <settings-list-element header="BOOST OPTIONS"
                                 sub-header="Force Play Cooldown">
            <settings-select :value="selectedCooldown"
                             :on-update="updateSelectedCooldown"
                             :items="['OFF', ...Array.from({ length: 60 }, (_, i) => (i + 1).toString())]">
              <span v-if="selectedCooldown !== 'OFF' && selectedCooldown === '1'"
                    class="text-primary text-subtitle-1 text-medium-emphasis">
                  min
              </span>
              <span v-else-if="selectedCooldown !== 'OFF'"
                    class="text-primary text-subtitle-1 text-medium-emphasis">
                  mins
              </span>
            </settings-select>
          </settings-list-element>
          <settings-list-element header="ROOM LIFESPAN"
                                 hint="At most 48 hours since creation."
                                 sub-header="Auto-delete the room after">
            <settings-select :value="selectedLifespan"
                             :on-update="updateSelectedLifespan"
                             :items="[4, 8, 16, 24, 48]">
              <span class="text-primary text-subtitle-1 text-medium-emphasis">
                  hours
              </span>
            </settings-select>
          </settings-list-element>
        </v-list>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
@import '@/assets/shared.css';
</style>
