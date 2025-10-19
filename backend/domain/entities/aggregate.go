package entities

import (
	"xsedox.com/main/domain/shared"
)

type AggregateRoot[T comparable] struct {
	Entity[T]
	DomainEvents []shared.IDomainEvent
}

func (ar *AggregateRoot[T]) AddDomainEvent(event shared.IDomainEvent) {
	ar.DomainEvents = append(ar.DomainEvents, event)
}
