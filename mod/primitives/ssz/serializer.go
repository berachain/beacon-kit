package ssz

import (
	"fmt"
	"reflect"
)

// SerializeError represents serialization errors.
type SerializeError struct {
	Err error
}

func (e *SerializeError) Error() string {
	return e.Err.Error()
}

// NewSerializeErrorMaximumLengthReached creates a new SerializeError for maximum length reached.
func NewSerializeErrorMaximumLengthReached(size int) *SerializeError {
	return &SerializeError{Err: fmt.Errorf("the encoded length is %d which meets or exceeds the maximum length %d", size, MaximumLength)}
}

// NewSerializeErrorInvalidInstance creates a new SerializeError for invalid instances.
func NewSerializeErrorInvalidInstance(err error) *SerializeError {
	return &SerializeError{Err: fmt.Errorf("invalid instance: %v", err)}
}

// NewSerializeErrorInvalidType creates a new SerializeError for invalid types.
func NewSerializeErrorInvalidType(err error) *SerializeError {
	return &SerializeError{Err: fmt.Errorf("invalid type: %v", err)}
}

// Serializable is an interface for data structures that can be serialized using SSZ.
type Serializable interface {
	Serialize(buffer *[]byte) (int, error)
}
type Serializer struct {
	ISerializer
}

type ISerializer interface {
	Elements(s SSZTypeGeneric) []SSZTypeGeneric
	MarshalSSZ(s SSZTypeGeneric) ([]byte, error)
}

func NewSerializer() Serializer {
	return *new(Serializer)
}

func RouteUint(val reflect.Value, typ reflect.Type) []byte {
	kind := typ.Kind()
	switch kind {
	case reflect.Uint8:
		return MarshalU8(val.Interface().(uint8))
	case reflect.Uint16:
		return MarshalU16(val.Interface().(uint16))
	case reflect.Uint32:
		return MarshalU32(val.Interface().(uint32))
	case reflect.Uint64:
		return MarshalU64(val.Interface().(uint64))
	// TODO(Chibera): Handle numbers over 64bit
	// case reflect.Uint128:
	// 	return MarshalU128(val.Interface().(uint128))
	// case reflect.Uint256:
	// 	return MarshalU256(val.Interface().(uint256))
	default:
		return make([]byte, 0)
	}
}

func IsUintLike(typ reflect.Type) bool {
	kind := typ.Kind()
	isUintLike := false

	switch kind {
	case reflect.Uint8:
		isUintLike = true
	case reflect.Uint16:
		isUintLike = true
	case reflect.Uint32:
		isUintLike = true
	case reflect.Uint64:
		fmt.Println("u64")
		isUintLike = true
	default:
		return isUintLike
	}

	fmt.Println(isUintLike)

	return isUintLike
}

func (s *Serializer) MarshalSSZ(c interface{}) ([]byte, error) {
	typ := reflect.TypeOf(c)
	k := typ.Kind()
	isUintLike := IsUintLike(typ)

	if isUintLike {
		return RouteUint(reflect.ValueOf(c), reflect.TypeOf(c)), nil
	}
	switch k {
	case reflect.Bool:
		return MarshalBool(reflect.ValueOf(c).Interface().(bool)), nil
	//TODO(Chibera): handle composite types. same algo for all 3
	// case KindList:
	// 	return true
	// case KindVector:
	// 	return IsVariableSize(t.(VectorType).ElemType())
	// case KindContainer:
	// 	for _, ft := range t.(ContainerType).FieldTypes() {
	// 		if IsVariableSize(ft) {
	// 			return true
	// 		}
	// 	}
	// 	return false
	default:
		return make([]byte, 0), nil
	}
}

// TODO - no longer needed?
func (s *Serializer) Elements(c SSZTypeGeneric) []SSZTypeGeneric {
	return make([]SSZTypeGeneric, 0)
}

// TODO

// func (b *structSSZ) Marshal(val reflect.Value, typ reflect.Type, buf []byte, startOffset uint64) (uint64, error) {
// 	if typ.Kind() == reflect.Ptr {
// 		if val.IsNil() {
// 			newVal := reflect.New(typ.Elem()).Elem()
// 			return b.Marshal(newVal, newVal.Type(), buf, startOffset)
// 		}
// 		return b.Marshal(val.Elem(), typ.Elem(), buf, startOffset)
// 	}
// 	fixedIndex := startOffset
// 	fixedLength := uint64(0)
// 	// For every field, we add up the total length of the items depending if they
// 	// are variable or fixed-size fields.
// 	for i := 0; i < typ.NumField(); i++ {
// 		// We skip protobuf related metadata fields.
// 		if strings.Contains(typ.Field(i).Name, "XXX_") {
// 			continue
// 		}
// 		fType, err := determineFieldType(typ.Field(i))
// 		if err != nil {
// 			return 0, err
// 		}
// 		if isVariableSizeType(fType) {
// 			fixedLength += BytesPerLengthOffset
// 		} else {
// 			if val.Type().Kind() == reflect.Ptr && val.IsNil() {
// 				elem := reflect.New(val.Type().Elem()).Elem()
// 				fixedLength += determineFixedSize(elem, fType)
// 			} else {
// 				fixedLength += determineFixedSize(val.Field(i), fType)
// 			}
// 		}
// 	}
// 	currentOffsetIndex := startOffset + fixedLength
// 	nextOffsetIndex := currentOffsetIndex
// 	for i := 0; i < typ.NumField(); i++ {
// 		// We skip protobuf related metadata fields.
// 		if strings.Contains(typ.Field(i).Name, "XXX_") {
// 			continue
// 		}
// 		fType, err := determineFieldType(typ.Field(i))
// 		if err != nil {
// 			return 0, err
// 		}
// 		factory, err := SSZFactory(val.Field(i), fType)
// 		if err != nil {
// 			return 0, err
// 		}
// 		if !isVariableSizeType(fType) {
// 			fixedIndex, err = factory.Marshal(val.Field(i), fType, buf, fixedIndex)
// 			if err != nil {
// 				return 0, err
// 			}
// 		} else {
// 			nextOffsetIndex, err = factory.Marshal(val.Field(i), fType, buf, currentOffsetIndex)
// 			if err != nil {
// 				return 0, err
// 			}
// 			// Write the offset.
// 			offsetBuf := make([]byte, BytesPerLengthOffset)
// 			binary.LittleEndian.PutUint32(offsetBuf, uint32(currentOffsetIndex-startOffset))
// 			copy(buf[fixedIndex:fixedIndex+BytesPerLengthOffset], offsetBuf)

// 			// We increase the offset indices accordingly.
// 			currentOffsetIndex = nextOffsetIndex
// 			fixedIndex += BytesPerLengthOffset
// 		}
// 	}
// 	return currentOffsetIndex, nil
// }

// // Part represents either a fixed sized part of the serialization
// // or an offset pointing to a variably sized part of the serialization.
// type Part struct {
// 	Fixed   []byte
// 	Offset  int
// 	IsFixed bool // Added to distinguish between Fixed and Offset without using Go interfaces
// }

// // Serializer is responsible for serializing data structures.
// type Serializer struct {
// 	Parts              []Part
// 	Variable           []byte
// 	FixedLengthsSum    int
// 	VariableLengthsSum int
// }

// // Serialize serializes the data structure into the buffer.
// func (s *Serializer) Serialize(buffer *[]byte) (int, error) {
// 	totalSize := s.FixedLengthsSum + s.VariableLengthsSum
// 	if totalSize >= MaximumLength {
// 		return 0, NewSerializeErrorMaximumLengthReached(totalSize)
// 	}

// 	var runningLength uint32 = uint32(s.FixedLengthsSum)
// 	for _, part := range s.Parts {
// 		if part.IsFixed {
// 			*buffer = append(*buffer, part.Fixed...)
// 		} else {
// 			bytesWritten, err := SerializeU32(runningLength, buffer)
// 			if err != nil {
// 				return 0, err
// 			}
// 			if bytesWritten != BytesPerLengthOffset {
// 				return 0, fmt.Errorf("unexpected number of bytes written: %d", bytesWritten)
// 			}
// 			runningLength += uint32(part.Offset)
// 		}
// 	}

// 	*buffer = append(*buffer, s.Variable...)

// 	return totalSize, nil
// }

// // WithElement adds an element to the serializer.
// func (s *Serializer) WithElement(element Serializable) error {
// 	var elementBuffer []byte
// 	bytesWritten, err := element.Serialize(&elementBuffer)
// 	if err != nil {
// 		return err
// 	}

// 	if IsVariableSize(element) {
// 		s.Parts = append(s.Parts, Part{Offset: bytesWritten, IsFixed: false})
// 		s.Variable = append(s.Variable, elementBuffer...)
// 		s.FixedLengthsSum += BytesPerLengthOffset
// 		s.VariableLengthsSum += bytesWritten
// 	} else {
// 		s.Parts = append(s.Parts, Part{Fixed: elementBuffer, IsFixed: true})
// 		s.FixedLengthsSum += bytesWritten
// 	}
// 	return nil
// }
