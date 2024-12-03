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

package common_test

import (
	"encoding/json"
	"testing"

	"github.com/berachain/beacon-kit/primitives/pkg/common"
	"github.com/berachain/beacon-kit/primitives/pkg/encoding/hex"
	"github.com/stretchr/testify/require"
)

func TestExecutionAddressMarshalling(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectedErr error
	}{
		{
			name:        "address too short",
			input:       []byte("\"0xab\""),
			expectedErr: hex.ErrInvalidHexStringLength,
		},
		{
			name:        "address missing hex prefix",
			input:       []byte("\"abc\""),
			expectedErr: hex.ErrMissingPrefix,
		},
		{
			name: "address too long",
			input: []byte(
				"\"0x000102030405060708090a0b0c0d0e0f101112131415161718\"",
			),
			expectedErr: hex.ErrInvalidHexStringLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				v   common.ExecutionAddress
				err error
			)
			require.NotPanics(t, func() {
				err = json.Unmarshal(tt.input, &v)
			})
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
