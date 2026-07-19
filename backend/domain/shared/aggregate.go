package shared

type AggregateRoot[T comparable] struct {
	Entity[T]
	domainEvents []IDomainEvent
}

func (ar *AggregateRoot[T]) RaiseDomainEvent(event IDomainEvent) {
	ar.domainEvents = append(ar.domainEvents, event)
}

func (ar *AggregateRoot[T]) ConsumeDomainEvents() []IDomainEvent {
	result := ar.domainEvents
	ar.domainEvents = []IDomainEvent{}
	return result
}
