package entities

type Entity[T comparable] struct {
	id T
}

func (e *Entity[T]) SetId(id T) {
	e.id = id
}
func (e *Entity[T]) GetId() T {
	return e.id
}
