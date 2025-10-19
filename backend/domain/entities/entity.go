package entities

type Entity[T comparable] struct {
	id T
}

func (e *Entity[T]) SetId(id T) {
	e.id = id
}
func (e *Entity[T]) Id() T {
	return e.id
}
