package music_data_enqueue_song_response_dto

import "github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"

type MusicDataEnqueueSongResponseDto struct {
	VideoId       string
	Title         string
	Author        string
	LengthSeconds uint16
	AlbumCoverUrl string
	MusicProvider music_provider.MusicProvider
	Isrc          *string
}
