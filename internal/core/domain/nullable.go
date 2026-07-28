package domain

type Nullable[T any] struct {
	Value *T
	Set   bool
}

func NewNullable[T any](value T) Nullable[T] {
	return Nullable[T]{
		Value: &value,
		Set:   true,
	}
}
