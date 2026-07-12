package search_song_query

type SearchSongQuery struct {
	Query         string  `json:"query" schema:"query" validate:"required,gte=2,lte=50" fname:"Song Query"`
	NextPageToken *string `json:"nextPageToken" schema:"nextPageToken" fname:"Next Page Token"`
	PageSize      int     `json:"pageSize" schema:"pageSize" validate:"gte=1,lte=50" fname:"Page Size"`
}
