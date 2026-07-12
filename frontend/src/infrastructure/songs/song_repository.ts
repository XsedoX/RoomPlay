import api_client from '../utils/api_client';
import type ISearchSongRequest from './ISearchSongRequest';
import type { IRepositoryResponse } from '../utils/IRepositoryResponse';
import type { IApiProblemDetailsResponse, IApiSuccessResponse } from '../utils/IApiResponse';
import { useNotificationStore } from '@/stores/notification_store';
import { TSnackbarColor } from '../utils/TSnackbarColor';
import type { ISearchSongResponse } from './ISearchSongResponse';

const URLS = {
  searchSongs: '/song/search',
};

export const SongRepository = {
  searchSongs: async (
    searchQuery: ISearchSongRequest,
  ): Promise<IRepositoryResponse<ISearchSongResponse[]>> => {
    console.log('searchSongs', searchQuery);
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
      .catch((error) => {
        const notificationStore = useNotificationStore();
        const problemDetails = error.response.data as IApiProblemDetailsResponse;
        notificationStore.showSnackbar(problemDetails.description, TSnackbarColor.ERROR);
        return {
          isSuccess: false,
        };
      });
  },
};
