package merkle_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
)

// Test getting a slice of the internal buffer and modifying.
func TestGet(t *testing.T) {
	buffer := merkle.NewBuffer[[32]byte]()

	testCases := []struct {
		size     int
		expected int
	}{
		{size: 0, expected: 0},
		{size: 1, expected: 1},
		{size: 5, expected: 5},
		{size: 16, expected: 16},
		{size: 33, expected: 33},
		{size: 17, expected: 17},
		{size: 100, expected: 100},
	}

	for i, tc := range testCases {
		result := buffer.Get(tc.size)

		if len(result) != tc.expected {
			t.Errorf(
				"Expected result size to be %d, got %d",
				tc.expected, len(result),
			)
		}

		// Ensure modifications to the underlying buffer persist.
		if i >= 1 {
			if result[0][i-1] != byte(i-1) {
				t.Errorf(
					"Expected result[0][%d] to be %b, got %d",
					i-1, byte(i-1), result[0][i-1],
				)
			}

			result[0] = [32]byte{}
			result[0][i] = byte(i)
		}
	}
}

// Benchmark for the Get method
//
// goos: darwin
// goarch: arm64
// pkg: github.com/berachain/beacon-kit/mod/primitives/pkg/merkle
// BenchmarkGet-12     173148679     6.917 ns/op     0 B/op     0 allocs/op
func BenchmarkGet(b *testing.B) {
	buffer := merkle.NewBuffer[[32]byte]()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ResetTimer()
	for range b.N {
		size := r.Intn(100) + 1
		result := buffer.Get(size)

		// Peform some operation on the result to avoid compiler optimizations.
		result[0] = [32]byte{}
		index := r.Intn(32)
		result[0][index] = byte(index)
	}
}
