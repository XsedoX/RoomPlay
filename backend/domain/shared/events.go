package shared

type IEventPublisher interface {
	Publish(event IDomainEvent) error
}
type IDomainEvent interface{}
