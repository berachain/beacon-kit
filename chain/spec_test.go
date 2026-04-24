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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
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

// Fork timestamps used by the fork-gated method tests. Values are chosen to be
// well-separated so boundary cases around each fork can be tested explicitly.
const (
	forkGatedGenesisTime  uint64 = 0
	forkGatedDeneb1Time   uint64 = 100
	forkGatedElectraTime  uint64 = 200
	forkGatedElectra1Time uint64 = 300
	forkGatedFuluTime     uint64 = 400
)

// Pre-Fulu hysteresis values. Distinct from the Fulu values so tests can detect
// when the wrong fork branch is taken.
const (
	preFuluHysteresisQuotient          uint64 = 4
	preFuluHysteresisUpwardMultiplier  uint64 = 5
	fuluHysteresisQuotient             uint64 = 40
	fuluHysteresisUpwardMultiplierFulu uint64 = 50
)

// EVM inflation per block values, distinct per fork region.
const (
	evmInflationPerBlockGenesis uint64 = 10
	evmInflationPerBlockDeneb1  uint64 = 20
	evmInflationPerBlockFulu    uint64 = 30
)

var (
	evmInflationAddrGenesis = common.MustNewExecutionAddressFromHex(
		"0x1111111111111111111111111111111111111111",
	)
	evmInflationAddrDeneb1 = common.MustNewExecutionAddressFromHex(
		"0x2222222222222222222222222222222222222222",
	)
	evmInflationAddrFulu = common.MustNewExecutionAddressFromHex(
		"0x3333333333333333333333333333333333333333",
	)
)

// buildForkGatedSpec builds a Spec with distinct fork-gated values so tests can
// assert which fork branch is taken for a given timestamp.
func buildForkGatedSpec(t *testing.T) chain.Spec {
	t.Helper()
	data := baseSpecData()
	data.GenesisTime = forkGatedGenesisTime
	data.Deneb1ForkTime = forkGatedDeneb1Time
	data.ElectraForkTime = forkGatedElectraTime
	data.Electra1ForkTime = forkGatedElectra1Time
	data.FuluForkTime = forkGatedFuluTime

	data.HysteresisQuotient = preFuluHysteresisQuotient
	data.HysteresisUpwardMultiplier = preFuluHysteresisUpwardMultiplier
	data.HysteresisQuotientFulu = fuluHysteresisQuotient
	data.HysteresisUpwardMultiplierFulu = fuluHysteresisUpwardMultiplierFulu

	data.EVMInflationAddressGenesis = evmInflationAddrGenesis
	data.EVMInflationPerBlockGenesis = evmInflationPerBlockGenesis
	data.EVMInflationAddressDeneb1 = evmInflationAddrDeneb1
	data.EVMInflationPerBlockDeneb1 = evmInflationPerBlockDeneb1
	data.EVMInflationAddressFulu = evmInflationAddrFulu
	data.EVMInflationPerBlockFulu = evmInflationPerBlockFulu

	s, err := chain.NewSpec(data)
	require.NoError(t, err)
	return s
}

// TestHysteresisQuotient_ForkBoundary asserts HysteresisQuotient returns the
// pre-Fulu value for all fork versions up to and including Electra1, and the
// Fulu value at or after FuluForkTime.
func TestHysteresisQuotient_ForkBoundary(t *testing.T) {
	t.Parallel()
	s := buildForkGatedSpec(t)

	tests := []struct {
		name      string
		timestamp uint64
		expected  math.U64
	}{
		{name: "At genesis (Deneb)", timestamp: forkGatedGenesisTime, expected: math.U64(preFuluHysteresisQuotient)},
		{name: "At Deneb1 fork", timestamp: forkGatedDeneb1Time, expected: math.U64(preFuluHysteresisQuotient)},
		{name: "At Electra fork", timestamp: forkGatedElectraTime, expected: math.U64(preFuluHysteresisQuotient)},
		{name: "At Electra1 fork", timestamp: forkGatedElectra1Time, expected: math.U64(preFuluHysteresisQuotient)},
		{name: "Just before Fulu", timestamp: forkGatedFuluTime - 1, expected: math.U64(preFuluHysteresisQuotient)},
		{name: "At Fulu fork", timestamp: forkGatedFuluTime, expected: math.U64(fuluHysteresisQuotient)},
		{name: "Just after Fulu", timestamp: forkGatedFuluTime + 1, expected: math.U64(fuluHysteresisQuotient)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := s.HysteresisQuotient(math.U64(tt.timestamp))
			require.Equal(t, tt.expected, got)
		})
	}
}

// TestHysteresisUpwardMultiplier_ForkBoundary asserts HysteresisUpwardMultiplier
// returns the pre-Fulu value before FuluForkTime and the Fulu value at or
// after FuluForkTime.
func TestHysteresisUpwardMultiplier_ForkBoundary(t *testing.T) {
	t.Parallel()
	s := buildForkGatedSpec(t)

	tests := []struct {
		name      string
		timestamp uint64
		expected  math.U64
	}{
		{name: "At genesis (Deneb)", timestamp: forkGatedGenesisTime, expected: math.U64(preFuluHysteresisUpwardMultiplier)},
		{name: "At Deneb1 fork", timestamp: forkGatedDeneb1Time, expected: math.U64(preFuluHysteresisUpwardMultiplier)},
		{name: "At Electra fork", timestamp: forkGatedElectraTime, expected: math.U64(preFuluHysteresisUpwardMultiplier)},
		{name: "At Electra1 fork", timestamp: forkGatedElectra1Time, expected: math.U64(preFuluHysteresisUpwardMultiplier)},
		{name: "Just before Fulu", timestamp: forkGatedFuluTime - 1, expected: math.U64(preFuluHysteresisUpwardMultiplier)},
		{name: "At Fulu fork", timestamp: forkGatedFuluTime, expected: math.U64(fuluHysteresisUpwardMultiplierFulu)},
		{name: "Just after Fulu", timestamp: forkGatedFuluTime + 1, expected: math.U64(fuluHysteresisUpwardMultiplierFulu)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := s.HysteresisUpwardMultiplier(math.U64(tt.timestamp))
			require.Equal(t, tt.expected, got)
		})
	}
}

// TestEVMInflationAddress_ForkBoundary asserts EVMInflationAddress returns the
// correct address for each of the three fork regions: genesis (Deneb),
// Deneb1/Electra/Electra1, and Fulu.
func TestEVMInflationAddress_ForkBoundary(t *testing.T) {
	t.Parallel()
	s := buildForkGatedSpec(t)

	tests := []struct {
		name      string
		timestamp uint64
		expected  common.ExecutionAddress
	}{
		{name: "At genesis (Deneb)", timestamp: forkGatedGenesisTime, expected: evmInflationAddrGenesis},
		{name: "Between genesis and Deneb1", timestamp: forkGatedDeneb1Time - 1, expected: evmInflationAddrGenesis},
		{name: "At Deneb1 fork", timestamp: forkGatedDeneb1Time, expected: evmInflationAddrDeneb1},
		{name: "At Electra fork", timestamp: forkGatedElectraTime, expected: evmInflationAddrDeneb1},
		{name: "At Electra1 fork", timestamp: forkGatedElectra1Time, expected: evmInflationAddrDeneb1},
		{name: "Just before Fulu", timestamp: forkGatedFuluTime - 1, expected: evmInflationAddrDeneb1},
		{name: "At Fulu fork", timestamp: forkGatedFuluTime, expected: evmInflationAddrFulu},
		{name: "Just after Fulu", timestamp: forkGatedFuluTime + 1, expected: evmInflationAddrFulu},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := s.EVMInflationAddress(math.U64(tt.timestamp))
			require.Equal(t, tt.expected, got)
		})
	}
}

// TestEVMInflationPerBlock_ForkBoundary asserts EVMInflationPerBlock returns
// the correct amount for each of the three fork regions: genesis (Deneb),
// Deneb1/Electra/Electra1, and Fulu.
func TestEVMInflationPerBlock_ForkBoundary(t *testing.T) {
	t.Parallel()
	s := buildForkGatedSpec(t)

	tests := []struct {
		name      string
		timestamp uint64
		expected  math.Gwei
	}{
		{name: "At genesis (Deneb)", timestamp: forkGatedGenesisTime, expected: math.Gwei(evmInflationPerBlockGenesis)},
		{name: "Between genesis and Deneb1", timestamp: forkGatedDeneb1Time - 1, expected: math.Gwei(evmInflationPerBlockGenesis)},
		{name: "At Deneb1 fork", timestamp: forkGatedDeneb1Time, expected: math.Gwei(evmInflationPerBlockDeneb1)},
		{name: "At Electra fork", timestamp: forkGatedElectraTime, expected: math.Gwei(evmInflationPerBlockDeneb1)},
		{name: "At Electra1 fork", timestamp: forkGatedElectra1Time, expected: math.Gwei(evmInflationPerBlockDeneb1)},
		{name: "Just before Fulu", timestamp: forkGatedFuluTime - 1, expected: math.Gwei(evmInflationPerBlockDeneb1)},
		{name: "At Fulu fork", timestamp: forkGatedFuluTime, expected: math.Gwei(evmInflationPerBlockFulu)},
		{name: "Just after Fulu", timestamp: forkGatedFuluTime + 1, expected: math.Gwei(evmInflationPerBlockFulu)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := s.EVMInflationPerBlock(math.U64(tt.timestamp))
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestValidate_ForkOrder_Success(t *testing.T) {
	t.Parallel()
	data := baseSpecData()
	data.GenesisTime = 10
	data.Deneb1ForkTime = 20
	data.ElectraForkTime = 30
	data.Electra1ForkTime = 40
	data.FuluForkTime = 50

	_, err := chain.NewSpec(data)
	require.NoError(t, err)
}

func TestValidate_ForkOrder_GenesisAfterDeneb(t *testing.T) {
	t.Parallel()
	data := baseSpecData()
	data.GenesisTime = 50
	data.Deneb1ForkTime = 20
	data.ElectraForkTime = 60
	data.Electra1ForkTime = 70
	data.FuluForkTime = 80

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
	data.Electra1ForkTime = 50
	data.FuluForkTime = 60

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
	data.Electra1ForkTime = 0
	data.FuluForkTime = 0

	_, err := chain.NewSpec(data)
	require.NoError(t, err)
}
