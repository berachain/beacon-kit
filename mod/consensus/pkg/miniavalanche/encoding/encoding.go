package encoding

import (
	"bytes"
	"encoding/gob"
)

func Encode[T any](obj T) ([]byte, error) {
	var (
		buf bytes.Buffer
		enc = gob.NewEncoder(&buf)
	)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode[T any](b []byte, obj *T) error {
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	return dec.Decode(obj)
}
