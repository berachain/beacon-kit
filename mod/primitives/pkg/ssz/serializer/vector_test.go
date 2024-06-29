// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package serializer_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/serializer"
)

// Mock type implementing the required interface for testing
type MockType struct {
	data []byte
}

func (m MockType) NewFromSSZ(data []byte) (MockType, error) {
	return MockType{data: data}, nil
}

func (m MockType) SizeSSZ() int {
	return 10
}

func BenchmarkUnmarshalVectorFixed(b *testing.B) {
	// Prepare a buffer with mock data
	elementSize := 10
	numElements := 1000
	buf := make([]byte, elementSize*numElements)
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i % 256)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := serializer.UnmarshalVectorFixed[MockType](buf)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
