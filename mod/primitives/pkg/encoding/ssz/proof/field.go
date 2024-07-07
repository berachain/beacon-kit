package proof

type Field struct {
	SSZType
	name string
}

// GetName returns the name of the field.
func (f Field) GetName() string {
	return f.name
}
