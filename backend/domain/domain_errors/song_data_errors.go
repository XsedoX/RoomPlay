package domain_errors

func NewSongDataSongLengthZeroError() error {
	return &DomainError{
		Code:        "SongData.LengthSeconds.Zero",
		Description: "Song length in seconds cannot be zero",
	}
}

func NewSongDataUrlEmptyError() error {
	return &DomainError{
		Code:        "SongData.Url.Empty",
		Description: "Song url cannot be empty",
	}
}

func NewSongDataTitleEmptyError() error {
	return &DomainError{
		Code:        "SongData.Title.Empty",
		Description: "Song title cannot be empty",
	}
}

func NewSongDataAuthorEmptyError() error {
	return &DomainError{
		Code:        "SongData.Author.Empty",
		Description: "Song author cannot be empty",
	}
}

func NewSongDataAlbumCoverUrlEmptyError() error {
	return &DomainError{
		Code:        "SongData.AlbumCoverUrl.Empty",
		Description: "Song album cover url cannot be empty",
	}
}

func NewSongDataIsrcIncorrectFormatError() error {
	return &DomainError{
		Code:        "SongData.Isrc.IncorrectFormat",
		Description: "ISRC must match the format: CC-XXX-YY-NNNNN without dashes.",
	}
}
