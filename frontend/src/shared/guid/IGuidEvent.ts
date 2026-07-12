import type { IGuid } from './Guid';

export default interface IGuidEvent {
  event: Event;
  id: IGuid;
}
