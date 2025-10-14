import type { IGuid } from '@/shared/Guid.ts';

export default interface IGuidEvent {
  event: Event;
  id: IGuid;
}
