package search_song_query_song_data_dto

type SearchSongQuerySongDataDto struct {
	VideoId       string `json:"videoId"`
	Author        string `json:"author"`
	AlbumCoverUrl string `json:"albumCoverUrl"`
	Title         string `json:"title"`
}
