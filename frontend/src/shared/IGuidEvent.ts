import type { IGuid } from '@/utils/Guid.ts';

export default interface IGuidEvent {
  event: Event;
  id: IGuid;
}
