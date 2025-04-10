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

package merkle

const (
	// ProposerIndexGIndexBlock is the generalized index of the proposer
	// index in the beacon block in the  fork.
	ProposerIndexGIndexBlock = 9

	// StateGIndexBlock is the generalized index of the beacon state in
	// the beacon block in the  fork.
	StateGIndexBlock = 11

	// ZeroValidatorPubkeyGIndexState is the generalized index of the 0
	// validator's pubkey in the beacon state in the  fork. To get the
	// GIndex of the pubkey of validator at index n, the formula is:
	// GIndex = ZeroValidatorPubkeyGIndexState +
	//          (ValidatorPubkeyGIndexOffset * n)
	ZeroValidatorPubkeyGIndexState = 439804651110400

	// ZeroValidatorPubkeyGIndexBlock is the generalized index of the 0
	// validator's pubkey in the beacon block in the  fork. This is
	// calculated by concatenating the (ZeroValidatorPubkeyGIndexState,
	// StateGIndexBlock) GIndices. To get the GIndex of the pubkey of
	// validator at index n, the formula is:
	// GIndex = ZeroValidatorPubkeyGIndexBlock +
	//          (ValidatorPubkeyGIndexOffset * n)
	ZeroValidatorPubkeyGIndexBlock = 3254554418216960

	// ValidatorPubkeyGIndexOffset is the offset of a validator pubkey GIndex.
	ValidatorPubkeyGIndexOffset = 8

	// ExecutionNumberGIndexState is the generalized index of the latest
	// execution payload header in the beacon state in the  fork.
	ExecutionNumberGIndexState = 774

	// ExecutionNumberGIndexBlock is the generalized index of the number
	// in the latest execution payload header in the beacon block in the 
	// fork. This is calculated by concatenating the
	// (ExecutionNumberGIndexState, StateGIndexBlock) GIndices.
	ExecutionNumberGIndexBlock = 5894

	// ExecutionFeeRecipientGIndexState is the generalized index of the
	// fee recipient in the latest execution payload header in the beacon state
	// in the  fork.
	ExecutionFeeRecipientGIndexState = 769

	// ExecutionFeeRecipientGIndexBlock is the generalized index of the
	// fee recipient in the latest execution payload header in the beacon block
	// in the  fork. This is calculated by concatenating the
	// (ExecutionFeeRecipientGIndexState, StateGIndexBlock) GIndices.
	ExecutionFeeRecipientGIndexBlock = 5889
)
