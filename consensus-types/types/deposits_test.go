// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// TestDeposits_MarshalSSZ verifies that marshalling Deposits produces a non-nil,
// correctly sized byte slice.
func TestDeposits_MarshalSSZ(t *testing.T) {
	t.Parallel()

	// Create a Deposits slice with two valid deposits.
	deposits := types.Deposits{generateValidDeposit(), generateValidDeposit()}

	// Marshal the deposits into SSZ format.
	data, err := (&deposits).MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	// Check that the data length equals the expected size.
	expectedSize := ssz.Size(&deposits)
	require.Len(t, data, int(expectedSize))
}

// TestDeposits_NewFromSSZ verifies that unmarshalling SSZ data returns an equivalent
// Deposits object.
func TestDeposits_NewFromSSZ(t *testing.T) {
	t.Parallel()

	// Create original Deposits with two valid deposits.
	originalDeposits := types.Deposits{generateValidDeposit(), generateValidDeposit()}
	data, err := (&originalDeposits).MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	// Unmarshal the data into a new Deposits instance.
	var newDeposits *types.Deposits
	newDeposits, err = newDeposits.NewFromSSZ(data)
	require.NoError(t, err)

	// Ensure the unmarshalled Deposits is equal to the original.
	require.Equal(t, originalDeposits, *newDeposits)
}
