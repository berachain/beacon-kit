package sszdb

import (
	"context"
	"fmt"
)

func DecodeListOfStaticElements[T any](
	ctx context.Context,
	db *SchemaDB,
	path ObjectPath,
	size int,
	dec func([]byte) (T, error),
) ([]T, error) {
	length, err := db.GetListLength(ctx, path)
	if err != nil {
		return nil, err
	}
	bz, err := db.GetPath(ctx, path)
	if err != nil {
		return nil, err
	}
	if len(bz)%size != 0 {
		return nil, fmt.Errorf(
			"expected multiple of %d bytes, got %d",
			size,
			len(bz),
		)
	}
	elements := make([]T, length)
	for i := range int(length) {
		elements[i], err = dec(bz[i*size : (i+1)*size])
		if err != nil {
			return nil, err
		}
	}
	return elements, nil
}
