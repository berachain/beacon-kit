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
//

package cometbft_test

import (
	"testing"

	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config/spec"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/primitives/crypto"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"
)

func TestSBTEnableHeightAlignedWithChainSpecs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		getChainSpecs func(t *testing.T) chain.Spec
	}{
		{
			name: "mainnet",
			getChainSpecs: func(t *testing.T) chain.Spec {
				t.Helper()
				cs, err := spec.MainnetChainSpec()
				require.NoError(t, err)
				return cs
			},
		},
		{
			name: "testnet",
			getChainSpecs: func(t *testing.T) chain.Spec {
				t.Helper()
				cs, err := spec.TestnetChainSpec()
				require.NoError(t, err)
				return cs
			},
		},
		{
			name: "devnet",
			getChainSpecs: func(t *testing.T) chain.Spec {
				t.Helper()
				cs, err := spec.DevnetChainSpec()
				require.NoError(t, err)
				return cs
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cs := tt.getChainSpecs(t)

			var cp *cmttypes.ConsensusParams
			require.NotPanics(t, func() {
				cp = cometbft.DefaultConsensusParams(crypto.CometBLSType, cs)
			})

			// Make sure that the Enable Height in consensus parameters matches with
			// chain spec one. This is relevant to make sure that consensus parameters
			// are well set when chain spec is modified and genesis is regenerated.
			require.Equal(t, cp.Feature.SBTEnableHeight, cs.SbtConsensusEnableHeight())
		})
	}
}
