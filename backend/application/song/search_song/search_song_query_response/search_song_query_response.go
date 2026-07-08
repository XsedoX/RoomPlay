package search_song_query_response

type SearchSongQueryResponse struct {
	VideoId       string `json:"videoId"`
	Author        string `json:"author"`
	AlbumCoverUrl string `json:"albumCoverUrl"`
	Title         string `json:"songTitle"`
	NextPageToken string `json:"nextPageToken"`
}
