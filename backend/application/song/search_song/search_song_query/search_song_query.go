package search_song_query

type SearchSongQuery struct {
	Query string `json:"query" validate:"required,song_query_validation|gte=2,lte=50" fname:"Song Query"`
}
