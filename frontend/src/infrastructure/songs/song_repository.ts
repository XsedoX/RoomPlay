import api_client from '../utils/api_client';
import type ISearchSongRequest from './ISearchSongRequest';
import type { IRepositoryResponse } from '../utils/IRepositoryResponse';
import type { IApiSuccessResponse } from '../utils/IApiResponse';
import type { ISearchSongResponse } from './ISearchSongResponse';

const URLS = {
  searchSongs: '/song/search',
};

export const SongRepository = {
  searchSongs: async (
    searchQuery: ISearchSongRequest,
  ): Promise<IRepositoryResponse<ISearchSongResponse[]>> => {
    setTimeout(() => {}, 10000);
    const params = new URLSearchParams();
    params.set('query', searchQuery.query);
    params.set('pageSize', searchQuery.pageSize.toString());
    if (searchQuery.nextPageToken) {
      params.set('nextPageToken', searchQuery.nextPageToken);
    }
    const url = `${URLS.searchSongs}?${params.toString()}`;
    return await api_client
      .get<IApiSuccessResponse<ISearchSongResponse[]>>(url)
      .then((response) => ({
        isSuccess: true,
        data: response.data.data,
        meta: response.data.meta,
      }))
      .catch(() => {
        return {
          isSuccess: false,
        };
      });
  },
};
