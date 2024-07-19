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

package types

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconBlockValidatorPubkey represents a block in the beacon chain during the
// DenebPlus fork, with only the minimally required values to prove a
// validator's pubkey.
//
//go:generate go run github.com/ferranbt/fastssz/sszgen --path block_validator_pubkey.go -objs BeaconBlockValidatorPubkey -include ../../../../../primitives/pkg/crypto,../../../../../primitives/pkg/common,../../../../../primitives/pkg/bytes,../../../../../consensus-types/pkg/types,../../../../../engine-primitives/pkg/engine-primitives,../../../../../primitives/pkg/math,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil,../../../../../primitives/pkg/common/common.go -output block_validator_pubkey.ssz.go
type BeaconBlockValidatorPubkey struct {
	// Slot represents the position of the block in the chain.
	Slot math.Slot
	// ProposerIndex is the index of the validator who proposed the block.
	ProposerIndex math.ValidatorIndex
	// ParentBlockRoot is the hash of the parent block.
	ParentBlockRoot common.Root
	// StateRoot is the summary of the BeaconState type with only the required
	// raw values to prove a validator's pubkey.
	StateRoot *struct {
		GenesisValidatorsRoot common.Root
		Slot                  math.Slot
		// Fork is the hash tree root of the Fork.
		Fork common.Root
		// LatestBlockHeader is the hash tree root of the latest block header.
		LatestBlockHeader common.Root
		// BlockRoots is the hash tree root of the block headers buffer.
		BlockRoots common.Root
		// StateRoots is the hash tree root of the beacon states buffer.
		StateRoots common.Root
		// Eth1Data is the hash tree root of the eth1 data.
		Eth1Data         common.Root
		Eth1DepositIndex uint64
		// LatestExecutionPayloadHeader is the hash tree root of the latest
		// execution payload header.
		LatestExecutionPayloadHeader common.Root
		// Validators is the list of the summary of the Validator type with only
		// the required raw values to prove a validator's pubkey.
		Validators []*struct {
			// Pubkey is the validator's 48-byte BLS public key.
			Pubkey                     crypto.BLSPubkey `ssz-size:"48"`
			WithdrawalCredentials      types.WithdrawalCredentials
			EffectiveBalance           math.Gwei
			Slashed                    bool
			ActivationEligibilityEpoch math.Epoch
			ActivationEpoch            math.Epoch
			ExitEpoch                  math.Epoch
			WithdrawableEpoch          math.Epoch
		} `ssz-max:"1099511627776"`
		// Balances is the hash tree root of the validator balances.
		Balances common.Root
		// RandaoMixes is the hash tree root of the randao mixes.
		RandaoMixes                  common.Root
		NextWithdrawalIndex          uint64
		NextWithdrawalValidatorIndex math.ValidatorIndex
		// Slashings is the hash tree root of the slashings.
		Slashings     common.Root
		TotalSlashing math.Gwei
	}
	// BodyRoot is the root of the block body.
	BodyRoot common.Root
}
