package domain_errors

func NewScheduledSongScheduledInPastError() error {
	return &DomainError{
		"ScheduledSong.ScheduledAtUtc.TimeBeforeNow",
		"Song cannot be scheduled in the past",
	}
}
