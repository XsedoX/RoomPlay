<script setup lang="ts">
import { ref } from 'vue';
import type IHostDeviceDto from '@/settings_page/choose_host_device_list/IHostDeviceDto.ts'
import { HostDeviceStateTypes } from '@/settings_page/choose_host_device_list/HostDeviceStateTypes.ts'
import { HostDeviceTypes } from '@/settings_page/choose_host_device_list/HostDeviceTypes.ts'
import { Guid } from '@/shared/Guid.ts'

const devicesRef = ref<string>()
const devices = ref<IHostDeviceDto[]>([
  {
    id: Guid.generate(),
    isHost: true,
    hostDeviceType: HostDeviceTypes.Computer,
    friendlyName: 'Living Room PC',
    state: HostDeviceStateTypes.Online,
    isCurrentDevice: true,
  },
  {
    id: Guid.generate(),
    isHost: false,
    hostDeviceType: HostDeviceTypes.Mobile,
    friendlyName: 'My Work Laptop',
    state: HostDeviceStateTypes.Offline,
    isCurrentDevice: false,
  },
  {
    id: Guid.generate(),
    isHost: false,
    hostDeviceType: HostDeviceTypes.Mobile,
    friendlyName: 'Personal Phone',
    state: HostDeviceStateTypes.Online,
    isCurrentDevice: false,
  }
]);
</script>

<template>
  <v-radio-group v-model="devicesRef">
    <v-radio :value="device.id"
             v-for="device in devices"
             :key="device.id.toString()">
      <template v-slot:label>
        <div class="d-flex ga-2 align-center">
          <v-icon v-if="device.hostDeviceType === HostDeviceTypes.Computer"
                  icon="computer"></v-icon>
          <v-icon v-else-if="device.hostDeviceType === HostDeviceTypes.Mobile"
                  icon="smartphone"></v-icon>
          <div class="d-flex flex-column">
            <div class="h3 font-weight-bold">
              {{device.friendlyName}}
            </div>
            <div class="text-medium-emphasis text-body-2">
              ({{device.state}})
            </div>
          </div>
        </div>
      </template>
    </v-radio>
  </v-radio-group>
</template>
