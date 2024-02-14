package bytes_test

import (
	"bytes"
	"testing"

	byteslib "github.com/itsdevbear/bolaris/lib/bytes"
)

func TestSafeCopy(t *testing.T) {
	tests := []struct {
		name     string
		original []byte
	}{
		{name: "Normal case", original: []byte{1, 2, 3, 4, 5}},
		{name: "Empty slice", original: []byte{}},
		{name: "Single element slice", original: []byte{9}},
		{name: "Large slice", original: make([]byte, 100)},
		{name: "Another normal case", original: []byte{6, 6, 6, 6, 6}},
		{name: "Another single element slice", original: []byte{5}},
		{name: "Another large slice", original: make([]byte, 200)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copied := byteslib.SafeCopy(tt.original)

			if !bytes.Equal(tt.original, copied) {
				t.Errorf("SafeCopy did not copy the slice correctly")
			}

			// Modifying the copied slice should not affect the original slice
			if len(copied) > 0 {
				copied[0] = 10
				if tt.original[0] == copied[0] {
					t.Errorf("Modifying the copied slice affected the original slice")
				}
			}
		})
	}
}
