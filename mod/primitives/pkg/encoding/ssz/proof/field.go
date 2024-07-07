package proof

type Field struct {
	SSZType
	name string
}

// NewField creates a new field.
func NewField(name string, typ SSZType) *Field {
	return &Field{name: name, SSZType: typ}
}

// GetName returns the name of the field.
func (f Field) GetName() string {
	return f.name
}
