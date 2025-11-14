import type { THostDevice } from '@/pages/settings_page/choose_host_device_list/THostDevice.ts';
import type { THostDeviceState } from '@/pages/settings_page/choose_host_device_list/THostDeviceState.ts';
import type { IGuid } from '@/shared/Guid.ts';

export default interface IHostDeviceDto {
  isHost: boolean;
  hostDeviceType: THostDevice;
  friendlyName: string;
  state: THostDeviceState;
  id: IGuid;
  isCurrentDevice: boolean;
}
