package search_song

type SearchSongQueryResponse struct {
	Url           string `json:"url"`
	Author        string `json:"author"`
	AlbumCoverUrl string `json:"albumCoverUrl"`
	Title         string `json:"songTitle"`
	LengthSeconds uint16 `json:"lengthSeconds"`
}
