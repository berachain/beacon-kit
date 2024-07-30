package config

type Config[T any] interface {
	IsNil() bool
	Default() T
}
