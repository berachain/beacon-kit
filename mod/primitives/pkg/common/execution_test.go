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
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
	"github.com/stretchr/testify/require"
)

func TestExecutionAddressMarshalling(t *testing.T) {
	var (
		v   common.ExecutionAddress
		err error
	)

	// No panic with hex string too short
	require.NotPanics(t, func() {
		err = json.Unmarshal([]byte("\"0xab\""), &v)
	})
	require.ErrorIs(t, err, hex.ErrInvalidHexStringLength)

	// No panic with hex string missing 0x prefix
	require.NotPanics(t, func() {
		err = json.Unmarshal([]byte("\"abc\""), &v)
	})
	require.ErrorIs(t, err, hex.ErrMissingPrefix)

	// Err upon trunctation on hex string too long
	err = json.Unmarshal(
		[]byte("\"0x000102030405060708090a0b0c0d0e0f101112131415161718\""),
		&v,
	)
	require.ErrorIs(t, err, hex.ErrInvalidHexStringLength)
}
