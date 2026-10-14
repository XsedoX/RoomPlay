package shared

import (
	"time"
)

type (
	EventName    string
	IDomainEvent interface {
		EventName() EventName
		OccurredAtUtc() time.Time
	}
)
