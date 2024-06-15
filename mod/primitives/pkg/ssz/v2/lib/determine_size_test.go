package ssz_test

import (
	"reflect"
	"testing"

	sszv2 "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/v2/lib"
)

func TestGetArrayDimensionality(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected int
	}{
		// TODO: add 2D and 3D empty array and slice
		{"1D empty array", [0]int32{}, 1},
		{"1D empty slice", []int32{}, 1},
		{"1D non-empty array", [3]int32{1, 2, 3}, 1},
		{"1D non-empty slice", []int32{1, 2, 3}, 1},
		{"2D empty array", [0][0]int32{}, 2},
		{"2D empty slice", [][]int32{}, 2},
		{"2D non-empty array", [2][3]int32{{1, 2, 3}, {4, 5, 6}}, 2},
		{"2D non-empty slice", [][]int32{{1, 2, 3}, {4, 5, 6}}, 2},
		{"3D empty array", [0][0][0]int32{}, 3},
		{"3D empty slice", [][][]int32{}, 3},
		{"3D non-empty array", [2][2][2]int32{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}}, 3},
		{"3D non-empty slice", [][][]int32{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}}, 3},
		{"1D empty byte array", [0]byte{}, 1},
		{"1D empty byte slice", []byte{}, 1},
		{"Byte array", [3]byte{1, 2, 3}, 1},
		{"Byte slice", []byte{1, 2, 3}, 1},
		{"2D empty bytes", [][]byte{}, 2},
		{"2D non-empty bytes", [][]byte{{1, 2, 3}, {4, 5, 6}}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := reflect.ValueOf(tt.input)
			result := sszv2.GetArrayDimensionality(val)
			if result != tt.expected {
				t.Errorf("Expected dimensionality %d, but got %d", tt.expected, result)
			}
		})
	}
}
