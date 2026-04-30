package youtube_music_data_provider

import "github.com/XsedoX/RoomPlay/application/dtos/music_data_response_dto"

type YoutubeMusicDataProvider struct{}

func NewYoutubeMusicDataProvider() *YoutubeMusicDataProvider {
	return &YoutubeMusicDataProvider{}
}

func (musicDataProvider *YoutubeMusicDataProvider) SearchSongsByQuery(accessToken, query string) (*[]music_data_response_dto.MusicDataResponseDto, error) {
}
