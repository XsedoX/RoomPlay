package search_song_query

type SearchSongQuery struct {
	Query         string  `json:"query" validate:"required,song_query_validation|gte=2,lte=50" fname:"Song Query"`
	NextPageToken *string `json:"next_page_token" fname:"Next Page Token"`
	PageSize      uint8   `json:"page_size" validate:"gte=1,lte=50" fname:"Page Size"`
}
