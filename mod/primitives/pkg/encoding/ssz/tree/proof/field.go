package proof

type Field[T any] struct {
	name  string
	value T
}

// NewField creates a new field.
func NewField[T any](name string, value T) *Field[T] {
	return &Field[T]{name: name, value: value}
}

// GetName returns the name of the field.
func (f Field[_]) GetName() string {
	return f.name
}

// GetValue returns the value of the field.
func (f Field[T]) GetValue() T {
	return f.value
}
