package song_data

import (
	"github.com/XsedoX/RoomPlay/domain/domain_errors/empty_string_domain_error"
	"github.com/XsedoX/RoomPlay/domain/domain_errors/validation_domain_error"
)

type SongData struct {
	url           string
	title         string
	author        string
	lengthSeconds uint16
	albumCoverUrl string
}

func (s SongData) AlbumCoverUrl() string {
	return s.albumCoverUrl
}

func (s SongData) Url() string {
	return s.url
}

func (s SongData) Title() string {
	return s.title
}

func (s SongData) Author() string {
	return s.author
}

func (s SongData) LengthSeconds() uint16 {
	return s.lengthSeconds
}

func HydrateSongData(
	url,
	title,
	author,
	albumCoverUrl string,
	lengthSeconds uint16,
) *SongData {
	return &SongData{
		url:           url,
		albumCoverUrl: albumCoverUrl,
		title:         title,
		author:        author,
		lengthSeconds: lengthSeconds,
	}
}

func NewSongData(
	url,
	title,
	author,
	albumCoverUrl string,
	lengthSeconds uint16,
) (*SongData, error) {
	if lengthSeconds == 0 {
		return nil, validation_domain_error.NewValidationDomainError(
			"SongData.LengthSeconds.Zero",
			"Song length in seconds cannot be zero",
		)
	}
	if url == "" {
		return nil, empty_string_domain_error.NewEmptyStringDomainError(
			"SongData.URL",
			"url",
		)
	}
	if title == "" {
		return nil, empty_string_domain_error.NewEmptyStringDomainError(
			"SongData.Title",
			"title",
		)
	}
	if author == "" {
		return nil, empty_string_domain_error.NewEmptyStringDomainError(
			"SongData.Author",
			"author",
		)
	}
	if albumCoverUrl == "" {
		return nil, empty_string_domain_error.NewEmptyStringDomainError(
			"SongData.AlbumCoverURL",
			"album cover url",
		)
	}

	return &SongData{
		url:           url,
		albumCoverUrl: albumCoverUrl,
		title:         title,
		author:        author,
		lengthSeconds: lengthSeconds,
	}, nil
}
