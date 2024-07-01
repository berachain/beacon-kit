package ssz

import (
	"fmt"
	"reflect"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/types"
)

var _ types.SSZEnumerable[*Container] = (*Container)(nil)

// Container is a container of SSZ types.
type Container struct {
	elements []types.BaseSSZType
}

// NewContainer creates a new Container from any struct, using reflection to get all the fields
// and put them into the elements list.
func NewContainer(v interface{}) (*Container, error) {
	val := reflect.ValueOf(v)

	// If v is a pointer, get the value it points to
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ensure v is a struct
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or pointer to struct")
	}

	// TODO: check struct tags to exclude fields.
	elements := make([]types.BaseSSZType, 0, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Check if the field implements SSZType
		if sszType, ok := field.Interface().(types.BaseSSZType); ok {
			elements = append(elements, sszType)
		} else {
			return nil, fmt.Errorf("field %s does not implement SSZType", val.Type().Field(i).Name)
		}
	}

	return &Container{elements: elements}, nil
}

func (c *Container) N() uint64 {
	return uint64(len(c.elements))
}

func (c *Container) Elements() []types.BaseSSZType {
	return c.elements
}

// MarshalSSZ marshals the container into SSZ format.
func (c *Container) MarshalSSZ() ([]byte, error) {
	bytes := make([]byte, 0)

	for _, element := range c.elements {
		elementBytes, err := element.MarshalSSZ()
		if err != nil {
			return nil, err
		}

		bytes = append(bytes, elementBytes...)
	}

	return bytes, nil
}

// SizeSSZ returns the size of the container in bytes.
func (c *Container) SizeSSZ() int {
	size := 0

	for _, element := range c.elements {
		size += element.SizeSSZ()
	}

	return size
}

// IsFixed returns true if the container is fixed size.
func (c *Container) IsFixed() bool {
	for _, element := range c.elements {
		if !element.IsFixed() {
			return false
		}
	}
	return true
}

// HashTreeRoot returns the hash tree root of the container.
func (c *Container) HashTreeRoot() ([32]byte, error) {
	return [32]byte{}, nil
}

// Type Definitions
func (c *Container) Type() types.Type {
	return types.Composite
}

// NewFromSSZ creates a new Container from SSZ format.
func (c *Container) NewFromSSZ(buf []byte) (*Container, error) {
	return nil, fmt.Errorf("not implemented")
}
