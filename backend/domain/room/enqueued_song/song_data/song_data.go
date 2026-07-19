package song_data

import (
	"regexp"

	"github.com/XsedoX/RoomPlay/domain/domain_errors"
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
		return nil, domain_errors.NewSongDataSongLengthZeroError()
	}
	if url == "" {
		return nil, domain_errors.NewSongDataUrlEmptyError()
	}
	if title == "" {
		return nil, domain_errors.NewSongDataTitleEmptyError()
	}
	if author == "" {
		return nil, domain_errors.NewSongDataAuthorEmptyError()
	}
	if albumCoverUrl == "" {
		return nil, domain_errors.NewSongDataAlbumCoverUrlEmptyError()
	}
	if isrc != nil {
		isrcConcrete := *isrc
		if len(isrcConcrete) > 12 || len(isrcConcrete) < 12 {
			return nil, domain_errors.NewSongDataIsrcIncorrectFormatError()
		}
		if !validIsrc.MatchString(isrcConcrete) {
			return nil, domain_errors.NewSongDataIsrcIncorrectFormatError()
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

func (s SongData) Equal(o SongData) bool {
	if s.url != o.url ||
		s.title != o.title ||
		s.author != o.author ||
		s.lengthSeconds != o.lengthSeconds ||
		s.albumCoverUrl != o.albumCoverUrl ||
		s.musicProvider != o.musicProvider {
		return false
	}
	if (s.isrc == nil) != (o.isrc == nil) {
		return false
	}
	if s.isrc != nil && *s.isrc != *o.isrc {
		return false
	}
	return true
}
