package music_data_response_dto

import (
	"github.com/XsedoX/RoomPlay/application/dtos/page_meta_dto"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
)

type SongDataResponseDto struct {
	VideoId       string
	Title         string
	Author        string
	AlbumCoverUrl string
	LengthSeconds uint16
	MusicProvider music_provider.MusicProvider
	Isrc          *string
}
type MusicDataResponseDto struct {
	Songs       []SongDataResponseDto
	PageMetaDto page_meta_dto.PageMetaDto
}
