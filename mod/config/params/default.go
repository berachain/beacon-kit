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

package params

import "github.com/ethereum/go-ethereum/common"

func DefaultBeaconConfig() BeaconChainConfig {
	//nolint:gomnd // default settings.
	return BeaconChainConfig{
		// Gwei value constants.
		MinDepositAmount:          uint64(1e9),
		MaxEffectiveBalance:       uint64(32e9),
		EffectiveBalanceIncrement: uint64(1e9),
		// Time parameters constants.
		SlotsPerEpoch:          8,
		SlotsPerHistoricalRoot: 1,
		// Eth1-related values.
		DepositContractAddress: common.HexToAddress(
			"0x00000000219ab540356cbb839cbe05303d7705fa",
		),
		// Fork-related values.
		ElectraForkEpoch: 9999999999999999,
		// State list length constants.
		EpochsPerHistoricalVector: 8,
		EpochsPerSlashingsVector:  1,
		HistoricalRootsLimit:      1,
		ValidatorRegistryLimit:    1099511627776,
		// Max operations per block constants.
		MaxDepositsPerBlock: 16,
		// Slashing
		ProportionalSlashingMultiplier: 1,
		// Capella values.
		MaxWithdrawalsPerPayload:         16,
		MaxValidatorsPerWithdrawalsSweep: 1 << 14,
		// Deneb values.
		MinEpochsForBlobsSidecarsRequest: 4096,
		MaxBlobCommitmentsPerBlock:       16,
		MaxBlobsPerBlock:                 6,
		FieldElementsPerBlob:             4096,
		BytesPerBlob:                     131072,
		KZGInclusionProofDepth:           17,
	}
}
