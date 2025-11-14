import { THostDevice } from '@/settings_page/choose_host_device_list/THostDevice.ts';

export const PlatformDiscoverer = {
  getDeviceType: (): THostDevice => {
    const userAgent = navigator.userAgent;
    if (/Mobi|Android|Tablet|iPad/i.test(userAgent)) {
      return THostDevice.Mobile;
    } else {
      return THostDevice.Desktop;
    }
  }
}
