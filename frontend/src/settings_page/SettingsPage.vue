<script setup lang="ts">
import PageTitle from '@/shared/page_title/PageTitle.vue';
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
    <v-list bg-color="transparent">
      <v-list-subheader>BOOST OPTIONS</v-list-subheader>
      <v-list-item>
        <v-sheet rounded="xl" color="surface-container pa-2 d-flex flex-nowrap ga-2 align-center">
          <div class="h6 font-weight-bold mr-auto">Force Play Cooldown</div>
          <v-select menu-icon="arrow_drop_down"
                    rounded="xl"
                    suffix="min"
                    no-data-text="OFF"
                    :hide-details="true"
                    max-width="200"
                    flat
                    item-color="primary"
                    :items="[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]"
                    variant="solo"
                    bg-color="transparent">
          </v-select>
        </v-sheet>
      </v-list-item>
    </v-list>
  </v-container>
</template>

<style scoped>
@import '@/assets/shared.css';
</style>
