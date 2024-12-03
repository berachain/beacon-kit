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

package hex_test

import (
	"testing"

	"github.com/berachain/beacon-kit/primitives/pkg/encoding/hex"
	"github.com/stretchr/testify/require"
)

func TestIsValidHex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "Valid hex string",
			input:   "0x48656c6c6f",
			wantErr: nil,
		},
		{
			name:    "Empty string",
			input:   "",
			wantErr: hex.ErrEmptyString,
		},
		{
			name:    "No 0x prefix",
			input:   "48656c6c6f",
			wantErr: hex.ErrMissingPrefix,
		},
		{
			name:    "Valid single hex character",
			input:   "0x0",
			wantErr: nil,
		},
		{
			name:    "Empty hex string",
			input:   "0x",
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := hex.IsValidHex(test.input)
			if test.wantErr != nil {
				require.ErrorIs(t, test.wantErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
