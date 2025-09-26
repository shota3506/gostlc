package environment

type Environment[T any] struct {
	name   string
	value  T
	parent *Environment[T]
}

func NewEnvironment[T any]() *Environment[T] {
	return &Environment[T]{}
}

func (r *Environment[T]) Lookup(name string) (T, bool) {
	var zero T

	if name == "" {
		return zero, false
	}
	if r.name == name {
		return r.value, true
	}
	if r.parent != nil {
		return r.parent.Lookup(name)
	}
	return zero, false
}

func (r *Environment[T]) Bind(name string, value T) *Environment[T] {
	return &Environment[T]{
		name:   name,
		value:  value,
		parent: r,
	}
}
