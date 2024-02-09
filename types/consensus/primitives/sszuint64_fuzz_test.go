// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package primitives_test

import (
	"testing"

	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

func FuzzSSZUint64Marshal(f *testing.F) {
	f.Fuzz(func(t *testing.T, data uint64) {
		slot := primitives.SSZUint64(data)

		// Attempt to marshal the slot
		marshalledData, err := slot.MarshalSSZ()
		if err != nil {
			t.Fatal("Failed to marshal slot")
		}

		// Initialize a new slot to unmarshal data into
		var newSSZUint64 primitives.SSZUint64
		err = newSSZUint64.UnmarshalSSZ(marshalledData)
		if err != nil {
			t.Fatal("Failed to unmarshal slot")
		}

		// Fuzz logic to check if the unmarshalled slot matches the original slot
		if newSSZUint64 != slot {
			t.Fatal("SSZUint64 mismatch after round-trip marshalling/unmarshalling")
		}
	})
}

func FuzzSSZUint64Unmarshal(f *testing.F) {
	f.Add([]byte{1, 2, 3, 4, 5, 6, 7, 8})
	f.Fuzz(func(t *testing.T, data []byte) {
		// Pad data if its length is less than 8
		if len(data) < 8 {
			data = append(data, make([]byte, 8-len(data))...)
		}

		slot := primitives.SSZUint64(0)
		err := slot.UnmarshalSSZ(data)
		if err != nil {
			t.Skip() // Skip the test if unmarshalling fails
		}

		// Fuzz logic to check if the unmarshalled slot can be
		// marshalled back and matches the original data
		marshalledData, err := slot.MarshalSSZ()
		if err != nil {
			t.Fatal("Failed to marshal slot")
		}
		if len(data) != len(marshalledData) {
			t.Fatal("Marshalled data length mismatch")
		}
		for i := range data {
			if data[i] != marshalledData[i] {
				t.Fatal("Data mismatch at index", i)
			}
		}
	})
}
