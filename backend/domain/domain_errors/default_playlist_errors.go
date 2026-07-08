package domain_errors

func NewDefaultPlaylistExternalIdEmptyError() error {
	return &DomainError{
		Code:        "DefaultPlaylist.ExternalId.EmptyString",
		Description: "The field 'external id' cannot be an empty.",
	}
}

func NewDefaultPlaylistSongsAmountZeroError() error {
	return &DomainError{
		Code:        "DefaultPlaylist.SongsAmount.Zero",
		Description: "Playlist has to have at least one song.",
	}
}

func NewDefaultPlaylistTitleEmptyError() error {
	return &DomainError{
		Code:        "DefaultPlaylist.Title.EmptyString",
		Description: "The field 'title' cannot be an empty.",
	}
}
