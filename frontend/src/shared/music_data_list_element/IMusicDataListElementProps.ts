import type IGuidEvent from '@/shared/IGuidEvent.ts';
import type { IGuid } from '@/utils/Guid.ts';

export interface IMusicDataListElementProps {
  musicDataListElementViewModel: IMusicDataListElementViewModel;
  onElementClick: (event: IGuidEvent) => void;
}
export interface IMusicDataListElementViewModel{
  title: string;
  subtitle: string;
  imageUrl: string;
  id: IGuid;
}
