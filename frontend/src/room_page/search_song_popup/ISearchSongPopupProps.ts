import type IGuidEvent from '@/shared/IGuidEvent.ts';

export default interface ISearchSongPopupProps {
  chooseSong: (event: IGuidEvent) => void;
}
