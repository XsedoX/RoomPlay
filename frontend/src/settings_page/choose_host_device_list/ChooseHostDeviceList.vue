<script setup lang="ts">
import { ref } from 'vue';
import type IHostDeviceDto from '@/settings_page/choose_host_device_list/IHostDeviceDto.ts';
import { THostDeviceState } from '@/settings_page/choose_host_device_list/THostDeviceState.ts';
import { THostDevice } from '@/settings_page/choose_host_device_list/THostDevice.ts';
import { Guid } from '@/shared/Guid.ts';

const devices = ref<IHostDeviceDto[]>([
  {
    id: Guid.generate(),
    isHost: true,
    hostDeviceType: THostDevice.Desktop,
    friendlyName: 'Living Room PC',
    state: THostDeviceState.Online,
    isCurrentDevice: false,
  },
  {
    id: Guid.generate(),
    isHost: false,
    hostDeviceType: THostDevice.Mobile,
    friendlyName: 'My Work Laptop',
    state: THostDeviceState.Offline,
    isCurrentDevice: false,
  },
  {
    id: Guid.generate(),
    isHost: false,
    hostDeviceType: THostDevice.Mobile,
    friendlyName: 'Personal Phone dwadwwwwwwwwwww',
    state: THostDeviceState.Online,
    isCurrentDevice: true,
  },
]);
const hostDevice = devices.value.find((d) => d.isHost)?.id;
const devicesRef = ref<string>(hostDevice!.toString());
function onDeviceSelected(id: string | null) {
  if (id === null) return;

  devicesRef.value = id;

  devices.value = devices.value.map((device) => {
    if (device.id.toString() === id) {
      return { ...device, isHost: true };
    } else {
      return { ...device, isHost: false };
    }
  });
}
function minIconWidth(device: IHostDeviceDto) {
  if(device.isHost) return '50px';
  if(device.isCurrentDevice) return '70px';
  return '0px';
}
</script>

<template>
  <v-radio-group :model-value="devicesRef"
                 @update:model-value="onDeviceSelected"
                 :hide-details="true">
    <template v-for="(device, index) in devices" :key="device.id.toString()">
      <v-radio
        :value="device.id.toString()"
        :name="device.id.toString()"
        :disabled="device.state === THostDeviceState.Offline"
        color="primary"
      >
        <template v-slot:label>
          <div class="d-flex ga-2 align-center justify-start py-2 w-100">
            <v-sheet
              color="transparent"
              :min-width="minIconWidth(device)"
              class="d-flex align-center justify-start ma-0 pa-0"
            >
              <v-badge
                color="primary"
                v-if="device.isCurrentDevice || device.isHost"
                :content="device.isHost ? 'Host' : device.isCurrentDevice ? 'Current' : undefined"
              >
                <v-icon
                  size="large"
                  :icon="
                    device.hostDeviceType === THostDevice.Desktop ? 'computer' : 'smartphone'
                  "
                ></v-icon>
              </v-badge>
              <v-icon
                v-else
                size="large"
                :icon="
                  device.hostDeviceType === THostDevice.Desktop ? 'computer' : 'smartphone'
                "
              ></v-icon>
            </v-sheet>
            <div class="d-flex flex-column">
              <div class="h6 font-weight-bold">
                {{ device.friendlyName }}
              </div>
              <div class="text-medium-emphasis text-body-2">({{ device.state }})</div>
            </div>
          </div>
        </template>
      </v-radio>
      <v-divider v-if="index < devices.length - 1"></v-divider>
    </template>
  </v-radio-group>
</template>

<style scoped>
:deep(.v-label) {
  width: 100%;
}
</style>
