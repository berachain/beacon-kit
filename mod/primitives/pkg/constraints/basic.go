package constraints

// ForkTyped is a constraint that requires a type to have an Empty method.
type ForkTyped[SelfT any] interface {
	EmptyWithVersion[SelfT]
	Versionable
	Nillable
}

// EmptyWithForkVersion is a constraint that requires a type to have an Empty method.
type EmptyWithVersion[SelfT any] interface {
	Empty(uint32) SelfT
}

// IsNil is a constraint that requires a type to have an IsNil method.
type Nillable interface {
	IsNil() bool
}

// Versionable is a constraint that requires a type to have a Version method.ßß
type Versionable interface {
	Version() uint32
}
