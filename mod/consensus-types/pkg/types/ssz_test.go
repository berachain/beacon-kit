package types

import (
	"testing"
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"
	"time"
)

func TestUint64SSZSerialize(t *testing.T) {
	testCases := []uint64{
		0,
		1,
		math.MaxUint64,
		rand.Uint64(),
	}

	for _, tc := range testCases {
		serialized, err := SSZSerialize(tc)
		if err != nil {
			t.Errorf("Failed to serialize uint64 %d: %v", tc, err)
			continue
		}

		expected := make([]byte, 8)
		binary.LittleEndian.PutUint64(expected, tc)

		if !bytes.Equal(serialized, expected) {
			t.Errorf("Serialization mismatch for uint64 %d: got %v, want %v", tc, serialized, expected)
		}
	}
}

func TestUint64SSZDeserialize(t *testing.T) {
	testCases := []uint64{
		0,
		1,
		math.MaxUint64,
		rand.Uint64(),
	}

	for _, tc := range testCases {
		serialized := make([]byte, 8)
		binary.LittleEndian.PutUint64(serialized, tc)

		var deserialized uint64
		err := SSZDeserialize(serialized, &deserialized)
		if err != nil {
			t.Errorf("Failed to deserialize uint64 %d: %v", tc, err)
			continue
		}

		if deserialized != tc {
			t.Errorf("Deserialization mismatch for uint64: got %d, want %d", deserialized, tc)
		}
	}
}

func BenchmarkUint64SSZSerialize(b *testing.B) {
	value := rand.Uint64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := SSZSerialize(value)
		if err != nil {
			b.Fatalf("Failed to serialize uint64: %v", err)
		}
	}
}

func BenchmarkUint64SSZDeserialize(b *testing.B) {
	value := rand.Uint64()
	serialized, _ := SSZSerialize(value)
	var deserialized uint64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := SSZDeserialize(serialized, &deserialized)
		if err != nil {
			b.Fatalf("Failed to deserialize uint64: %v", err)
		}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}