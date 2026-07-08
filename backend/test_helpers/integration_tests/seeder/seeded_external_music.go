package seeder

import (
	"github.com/XsedoX/RoomPlay/application/dtos/music_data_response_dto"
	"github.com/brianvoe/gofakeit/v7"
)

var ExternalSongData = []music_data_response_dto.MusicDataResponseDto{
	{
		VideoId:       gofakeit.ID(),
		Title:         gofakeit.SongName(),
		Author:        gofakeit.SongArtist(),
		AlbumCoverUrl: gofakeit.URL(),
		NextPageToken: gofakeit.ID(),
	},
	{
		VideoId:       gofakeit.ID(),
		Title:         gofakeit.SongName(),
		Author:        gofakeit.SongArtist(),
		AlbumCoverUrl: gofakeit.URL(),
		NextPageToken: gofakeit.ID(),
	},
	{
		VideoId:       gofakeit.ID(),
		Title:         gofakeit.SongName(),
		Author:        gofakeit.SongArtist(),
		AlbumCoverUrl: gofakeit.URL(),
		NextPageToken: gofakeit.ID(),
	},
}
