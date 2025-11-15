<script setup lang="ts">
import type IUserListElementDto from '@/pages/settings_page/users_list/IUserListElementDto.ts';
import { Guid, type IGuid } from '@/shared/Guid.ts';
import { ref } from 'vue';

const users = ref<IUserListElementDto[]>([
  { id: Guid.generate(), name: 'John', surname: 'Doe' },
  { id: Guid.generate(), name: 'Jane', surname: 'Smith' },
  { id: Guid.generate(), name: 'Peter', surname: 'Jones' },
  { id: Guid.generate(), name: 'John2', surname: 'Doe2' },
  { id: Guid.generate(), name: 'Jane3', surname: 'Smith3' },
  { id: Guid.generate(), name: 'Peter2', surname: 'Jones2' },
]);
function deleteUser(id: IGuid) {
  const index = users.value.findIndex((user) => user.id === id);
  if (index === -1) throw new Error(`Could not find user with id ${id} in users: ${users.value}`);
  users.value.splice(index, 1);
}
function blockUser(id: IGuid) {
  const index = users.value.findIndex((user) => user.id === id);
  if (index === -1) throw new Error(`Could not find user with id ${id} in users: ${users.value}`);
  users.value.splice(index, 1);
}
</script>
<template>
  <div class="d-flex flex-column w-100">
    <div v-if="users.length === 0" class="h3 font-weight-bold text-center">
      Do you have friends? If so invite them to your room!
    </div>
    <div v-for="(user, index) in users" :key="user.id.toString()" class="d-flex flex-column">
      <div class="d-flex align-center">
        <v-avatar color="primary" size="small">
          {{ user.name[0] + user.surname[0]! }}
        </v-avatar>
        <div class="d-flex flex-column ml-2">
          <div class="h6 font-weight-bold text-break">
            {{ user.name }}
          </div>
          <div class="text-caption text-medium-emphasis text-break">
            {{ user.surname }}
          </div>
        </div>
        <div class="ml-auto d-flex align-center">
          <v-btn color="primary" @click="blockUser(user.id)" variant="text" icon="block"></v-btn>
          <v-btn variant="text" @click="deleteUser(user.id)" color="primary" icon="delete"></v-btn>
        </div>
      </div>
      <v-divider v-if="index !== users.length - 1"></v-divider>
    </div>
  </div>
</template>
