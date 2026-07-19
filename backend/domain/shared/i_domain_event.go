package shared

import (
	"time"
)

type (
	EventName    string
	IDomainEvent interface {
		EventName() EventName
		OccurredAt() time.Time
	}
)
