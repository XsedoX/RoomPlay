import type { IGuid } from '@/shared/Guid.ts';

export default interface IMusicDataListElementDto {
  title: string;
  subtitle: string;
  imageUrl: string;
  id: IGuid;
}
