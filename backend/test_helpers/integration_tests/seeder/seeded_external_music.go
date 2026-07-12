package seeder

import (
	"github.com/XsedoX/RoomPlay/application/dtos/music_data_response_dto"
	"github.com/XsedoX/RoomPlay/application/dtos/page_meta_dto"
	"github.com/brianvoe/gofakeit/v7"
)

var ExternalSongData = music_data_response_dto.MusicDataResponseDto{
	Songs: []music_data_response_dto.SongDataResponseDto{
		{
			VideoId:       gofakeit.ID(),
			Title:         gofakeit.SongName(),
			Author:        gofakeit.SongArtist(),
			AlbumCoverUrl: gofakeit.URL(),
		},
		{
			VideoId:       gofakeit.ID(),
			Title:         gofakeit.SongName(),
			Author:        gofakeit.SongArtist(),
			AlbumCoverUrl: gofakeit.URL(),
		},
		{
			VideoId:       gofakeit.ID(),
			Title:         gofakeit.SongName(),
			Author:        gofakeit.SongArtist(),
			AlbumCoverUrl: gofakeit.URL(),
		},
	},
	PageMetaDto: page_meta_dto.PageMetaDto{
		NextPageToken:     new(gofakeit.ID()),
		PreviousPageToken: new(gofakeit.ID()),
		HasNextPage:       true,
		PageSize:          10,
	},
}
