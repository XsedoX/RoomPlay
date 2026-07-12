package search_song_query_dto

import (
	"github.com/XsedoX/RoomPlay/application/dtos/page_meta_dto"
	"github.com/XsedoX/RoomPlay/application/song/search_song/search_song_query_song_data_dto"
)

type SearchSongQueryDto struct {
	PageMetaDto page_meta_dto.PageMetaDto
	Songs       []search_song_query_song_data_dto.SearchSongQuerySongDataDto
}
