// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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
	"time"

	"github.com/itsdevbear/bolaris/primitives"
)

func TestSlot_Casting(t *testing.T) {
	slot := primitives.Slot(42)

	t.Run("time.Duration", func(t *testing.T) {
		if uint64(time.Duration(slot)) != uint64(slot) {
			t.Error("Slot should produce the same result with time.Duration")
		}
	})

	t.Run("floats", func(t *testing.T) {
		var x1 float32 = 42.2
		if primitives.Slot(x1) != slot {
			t.Errorf("Unequal: %v = %v", primitives.Slot(x1), slot)
		}

		var x2 = 42.2
		if primitives.Slot(x2) != slot {
			t.Errorf("Unequal: %v = %v", primitives.Slot(x2), slot)
		}
	})

	t.Run("int", func(t *testing.T) {
		var x = 42
		if primitives.Slot(x) != slot {
			t.Errorf("Unequal: %v = %v", primitives.Slot(x), slot)
		}
	})
}

func TestSlot_Marshalling(t *testing.T) {
	slot := primitives.Slot(42)

	t.Run("MarshalSSZ", func(t *testing.T) {
		bz, err := slot.MarshalSSZ()
		if err != nil {
			t.Error("Error during marshalling:", err)
		}
		if len(bz) != 8 {
			t.Error("Marshalled length should be 8")
		}
	})

	t.Run("MarshalSSZTo", func(t *testing.T) {
		dst := make([]byte, 0)
		bz, err := slot.MarshalSSZTo(dst)
		if err != nil {
			t.Error("Error during marshalling:", err)
		}
		if len(bz) != 8 {
			t.Error("Marshalled length should be 8")
		}
	})

	t.Run("UnmarshalSSZ", func(t *testing.T) {
		// big endian
		bz := []byte{42, 0, 0, 0, 0, 0, 0, 0}
		var s primitives.Slot
		err := s.UnmarshalSSZ(bz)
		if err != nil {
			t.Error("Error during unmarshalling:", err)
		}
		if s != slot {
			t.Errorf("Unequal: %v = %v", s, slot)
		}
	})

	t.Run("SizeSSZ", func(t *testing.T) {
		if slot.SizeSSZ() != 8 {
			t.Error("Size should be 8")
		}
	})
}
