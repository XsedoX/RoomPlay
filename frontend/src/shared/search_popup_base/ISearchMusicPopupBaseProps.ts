import type IMusicDataListElementDto from '@/shared/music_data_list_element/IMusicDataListElementDto.ts';

export default interface ISearchMusicPopupBaseProps {
  searchBoxPlaceholder: string;
  musicList: IMusicDataListElementDto[];
}
