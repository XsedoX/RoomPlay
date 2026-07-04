package song_data

import (
	"regexp"

	"github.com/XsedoX/RoomPlay/domain/domain_errors/empty_string_domain_error"
	"github.com/XsedoX/RoomPlay/domain/domain_errors/validation_domain_error"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
)

var validIsrc = regexp.MustCompile(`^[A-Z]{2}[A-Z0-9]{3}\d{2}\d{5}$`)

type SongData struct {
	url           string
	title         string
	author        string
	lengthSeconds uint16
	albumCoverUrl string
	musicProvider music_provider.MusicProvider
	isrc          *string
}

func (s SongData) Isrc() *string {
	return s.isrc
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

func (s SongData) MusicProvider() music_provider.MusicProvider {
	return s.musicProvider
}

func HydrateSongData(
	url,
	title,
	author,
	albumCoverUrl string,
	lengthSeconds uint16,
	musicProvider music_provider.MusicProvider,
	isrc *string,
) *SongData {
	return &SongData{
		url:           url,
		isrc:          isrc,
		albumCoverUrl: albumCoverUrl,
		title:         title,
		author:        author,
		lengthSeconds: lengthSeconds,
		musicProvider: musicProvider,
	}
}

func NewSongData(
	url,
	title,
	author,
	albumCoverUrl string,
	lengthSeconds uint16,
	musicProvider music_provider.MusicProvider,
	isrc *string,
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
	if isrc != nil {
		isrcConcrete := *isrc
		if len(isrcConcrete) > 12 || len(isrcConcrete) < 12 {
			return nil, validation_domain_error.NewValidationDomainError(
				"SongData.ISRC.Length",
				"ISRC must be exactly 12 characters long",
			)
		}
		if !validIsrc.MatchString(isrcConcrete) {
			return nil, validation_domain_error.NewValidationDomainError(
				"SongData.ISRC.Format",
				"ISRC must match the format: CC-XXX-YY-NNNNN",
			)
		}
	}

	return &SongData{
		url:           url,
		albumCoverUrl: albumCoverUrl,
		title:         title,
		author:        author,
		lengthSeconds: lengthSeconds,
		musicProvider: musicProvider,
		isrc:          isrc,
	}, nil
}
