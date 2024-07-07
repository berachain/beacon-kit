package types

type Type uint8

const (
	Basic Type = iota
	Vector
	List
	Container
)

// IsBasic returns true if the type is a basic type.
func (t Type) IsBasic() bool {
	return t == Basic
}

// IsElements returns true if the type is an enumerable type.
func (t Type) IsElements() bool {
	return t == Vector || t == List
}

// IsComposite returns true if the type is a composite type.
func (t Type) IsComposite() bool {
	return t == Vector || t == List || t == Container
}

// IsContainer returns true if the type is a container type.
func (t Type) IsContainer() bool {
	return t == Container
}
