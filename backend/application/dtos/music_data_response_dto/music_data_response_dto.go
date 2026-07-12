package music_data_response_dto

import "github.com/XsedoX/RoomPlay/application/dtos/page_meta_dto"

type SongDataResponseDto struct {
	VideoId       string
	Title         string
	Author        string
	AlbumCoverUrl string
}
type MusicDataResponseDto struct {
	Songs       []SongDataResponseDto
	PageMetaDto page_meta_dto.PageMetaDto
}
