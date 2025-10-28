package shared

type AggregateRoot[T comparable] struct {
	Entity[T]
	domainEvents []IDomainEvent
}

func (ar *AggregateRoot[T]) RaiseDomainEvent(event IDomainEvent) {
	ar.domainEvents = append(ar.domainEvents, event)
}
