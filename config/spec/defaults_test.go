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

package spec_test

import (
	"testing"

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/stretchr/testify/require"
)

func TestDomainTypeConversion(t *testing.T) {
	cs := spec.MainnetChainSpecData()
	require.Equal(t, bytes.B4([]byte{0x00, 0x00, 0x00, 0x00}), cs.DomainTypeProposer)
	require.Equal(t, bytes.B4([]byte{0x01, 0x00, 0x00, 0x00}), cs.DomainTypeAttester)
	require.Equal(t, bytes.B4([]byte{0x02, 0x00, 0x00, 0x00}), cs.DomainTypeRandao)
	require.Equal(t, bytes.B4([]byte{0x03, 0x00, 0x00, 0x00}), cs.DomainTypeDeposit)
	require.Equal(t, bytes.B4([]byte{0x04, 0x00, 0x00, 0x00}), cs.DomainTypeVoluntaryExit)
	require.Equal(t, bytes.B4([]byte{0x05, 0x00, 0x00, 0x00}), cs.DomainTypeSelectionProof)
	require.Equal(t, bytes.B4([]byte{0x06, 0x00, 0x00, 0x00}), cs.DomainTypeAggregateAndProof)
	require.Equal(t, bytes.B4([]byte{0x00, 0x00, 0x00, 0x01}), cs.DomainTypeApplicationMask)
}
