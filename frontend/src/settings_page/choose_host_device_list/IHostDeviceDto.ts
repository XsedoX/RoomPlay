import type { HostDeviceTypes } from '@/settings_page/choose_host_device_list/HostDeviceTypes.ts';
import type { HostDeviceStateTypes } from '@/settings_page/choose_host_device_list/HostDeviceStateTypes.ts';
import type { IGuid } from '@/shared/Guid.ts';

export default interface IHostDeviceDto {
  isHost: boolean;
  hostDeviceType: HostDeviceTypes;
  friendlyName: string;
  state: HostDeviceStateTypes;
  id: IGuid;
  isCurrentDevice: boolean;
}
