<script setup lang="ts">
import LogoWithTitleText from '@/shared/LogoWithTitleText.vue';
import AvatarWithFullName from '@/pages/main_menu_page/AvatarWithFullName.vue';
import JoinRoomPopup from '@/pages/main_menu_page/JoinRoomPopup.vue';
import CreateRoomPopup from '@/pages/main_menu_page/CreateRoomPopup.vue';
import { useUserStore } from '@/stores/user_store.ts';
import { onMounted } from 'vue';
import { useRoomStore } from '@/stores/room_store.ts';
import { useRouter } from 'vue-router';

const userStore = useUserStore();
const roomStore = useRoomStore();
const router = useRouter();

onMounted(async ()=>{
  if(await roomStore.getUserRoomMembership()){
    await router.replace({ name: 'RoomPage' })
  }
})
async function logout() {
  await userStore.logout()
  await router.replace({ name: 'LoginPage' });
}
</script>

<template>
  <v-container fluid class="fill-height align-start justify-center">
    <v-row justify="end"
           no-gutters
           class="w-100">
      <v-col cols="auto"
             class="text-center">
        <AvatarWithFullName :full-name="`${userStore.user?.name} ${userStore.user?.surname}`" :avatar-abbreviation="userStore.usersInitials"></AvatarWithFullName>
      </v-col>
    </v-row>
    <v-container class="pa-0 ma-0 justify-center align-center d-flex flex-column">
      <v-row class="w-100"
             justify="center">
        <v-col cols="8"
               sm="6"
               md="3"
               class="text-center pa-2">
          <LogoWithTitleText></LogoWithTitleText>
        </v-col>
      </v-row>
      <v-row class="w-100"
             justify="center"
             align="start">
        <v-col cols="8"
               sm="6"
               md="3"
               class="text-center">
          <v-btn data-testid="join-room-btn"
                 rounded="xl"
                 color="primary">
            Join a Room
            <JoinRoomPopup/>
          </v-btn>
        </v-col>
      </v-row>
      <v-row class="w-100"
             justify="center">
        <v-col cols="8"
               sm="6"
               md="3"
               class="text-center">
          <v-btn data-testid = "create-room-btn"
                 rounded="xl"
                 variant="outlined"
                 color="primary">
            Create a Room
            <CreateRoomPopup/>
          </v-btn>
        </v-col>
      </v-row>
      <v-row class="w-100"
             justify="center">
        <v-col cols="8"
               sm="6"
               md="3"
               class="text-center">
          <v-btn data-testid = "logout-btn"
                 variant="plain"
                 rounded="xl"
                 :ripple="false"
                 @click="logout()"
                 size="small"
                 color="error">
            Logout
          </v-btn>
        </v-col>
      </v-row>
    </v-container>
  </v-container>
</template>

