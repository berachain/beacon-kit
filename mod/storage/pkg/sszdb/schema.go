package sszdb

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
)

func CreateSchema(obj any) (schema.SSZType, error) {
	typ := reflect.TypeOf(obj)
	return traverseMonolith(typ, nil)
}

func traverseMonolith(
	typ reflect.Type,
	field *reflect.StructField,
) (schema.SSZType, error) {
	kind := typ.Kind()

	switch kind {
	case reflect.Ptr:
		return traverseMonolith(typ.Elem(), field)
	case reflect.Bool:
		return schema.Bool(), nil
	case reflect.Uint8:
		return schema.U8(), nil
	case reflect.Uint16:
		return schema.U16(), nil
	case reflect.Uint32:
		return schema.U32(), nil
	case reflect.Uint64:
		return schema.U64(), nil
	case reflect.Slice:
		// hack: slices with an `ssz-size` tag to be treated as vectors.
		length, ok, err := getFastSSZTag(field, "ssz-size")
		if err != nil {
			return nil, err
		}
		var elemType schema.SSZType
		if ok {
			// vector
			elemType, err = traverseMonolith(typ.Elem(), nil)
			if err != nil {
				return nil, err
			}
			return schema.DefineVector(elemType, length), nil
		} else {
			// list
			length, ok, err = getFastSSZTag(field, "ssz-max")
			if !ok {
				return nil, err
			}
			elemType, err = traverseMonolith(typ.Elem(), nil)
			if err != nil {
				return nil, err
			}
			return schema.DefineList(elemType, length), nil
		}
	case reflect.Array:
		// vector
		elemType, err := traverseMonolith(typ.Elem(), nil)
		if err != nil {
			return nil, err
		}
		return schema.DefineVector(elemType, uint64(typ.Len())), nil
	case reflect.Struct:
		var fields []*schema.Field[schema.SSZType]
		for _, field := range flattenStructFields(typ) {
			sszType, err := traverseMonolith(field.Type, &field)
			if err != nil {
				return nil, err
			}
			fields = append(fields, schema.NewField(field.Name, sszType))
		}
		return schema.DefineContainer(fields...), nil
	default:
		return nil, fmt.Errorf("unsupported type: %v", kind)
	}
}

// getFastSSZTag returns the value of a struct field tag as a uint64.
// These tags are required by ferranbt/fastssz to generate SSZ serialization code
// and reused here for similar metadata.
func getFastSSZTag(
	field *reflect.StructField,
	tag string,
) (uint64, bool, error) {
	str := field.Tag.Get(tag)
	if str == "" {
		return 0, false, nil
	}
	multi := strings.Split(str, ",")
	if len(multi) > 1 {
		return 0, false, nil
	}
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf(
			"tag %s value %s not an integer: %w", tag, str, err)
	}
	return i, true, nil
}

func flattenStructFields(typ reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	for i := range typ.NumField() {
		field := typ.Field(i)
		if field.Anonymous {
			// flatten embedded struct fields
			embedded := flattenStructFields(field.Type)
			fields = append(fields, embedded...)
		} else {
			fields = append(fields, field)
		}
	}
	return fields
}
