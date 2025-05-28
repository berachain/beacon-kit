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

package chain_test

import (
	"testing"

	"github.com/berachain/beacon-kit/chain"
	"github.com/stretchr/testify/require"
)

func baseSpecData() *chain.SpecData {
	return &chain.SpecData{
		// satisfy the pre-checks in validate()
		MaxWithdrawalsPerPayload: 2,
		ValidatorSetCap:          100,
		ValidatorRegistryLimit:   100,
	}
}

func TestValidate_ForkOrder_Success(t *testing.T) {
	t.Parallel()
	data := baseSpecData()
	data.GenesisTime = 10
	data.Deneb1ForkTime = 20
	data.ElectraForkTime = 30

	_, err := chain.NewSpec(data)
	require.NoError(t, err)
}

func TestValidate_ForkOrder_GenesisAfterDeneb(t *testing.T) {
	t.Parallel()
	data := baseSpecData()
	data.GenesisTime = 50
	data.Deneb1ForkTime = 20
	data.ElectraForkTime = 60

	_, err := chain.NewSpec(data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "timestamp at index 0 (50) > index 1 (20)")
}

func TestValidate_ForkOrder_DenebAfterElectra(t *testing.T) {
	t.Parallel()
	data := baseSpecData()
	data.GenesisTime = 10
	data.Deneb1ForkTime = 80
	data.ElectraForkTime = 40

	_, err := chain.NewSpec(data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "timestamp at index 1 (80) > index 2 (40)")
}

func TestValidate_ForkOrder_AllForksAtGenesis(t *testing.T) {
	t.Parallel()
	data := baseSpecData()
	data.GenesisTime = 0
	data.Deneb1ForkTime = 0
	data.ElectraForkTime = 0

	_, err := chain.NewSpec(data)
	require.NoError(t, err)
}
