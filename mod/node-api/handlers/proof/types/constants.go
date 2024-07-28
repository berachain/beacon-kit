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

const (
	// ZeroValidatorPubkeyGIndexDenebState is the generalized index of the 0
	// validator's pubkey in the beacon state in the Deneb fork. To get the
	// GIndex of the pubkey of validator at index n, the formula is:
	// GIndex = ZeroValidatorPubkeyGIndexDenebState +
	//          (ValidatorPubkeyGIndexOffset * n)
	ZeroValidatorPubkeyGIndexDenebState = 439804651110400

	// StateGIndexDenebBlock is the generalized index of the beacon state in
	// the beacon block in the Deneb fork.
	StateGIndexDenebBlock = 11

	// ZeroValidatorPubkeyGIndexDenebBlock is the generalized index of the 0
	// validator's pubkey in the beacon block in the Deneb fork. This is
	// calculated by concatenating the (ZeroValidatorPubkeyGIndexDenebState,
	// StateGIndexDenebBlock) GIndices. To get the GIndex of the pubkey of
	// validator at index n, the formula is:
	// GIndex = ZeroValidatorPubkeyGIndexDenebBlock +
	//          (ValidatorPubkeyGIndexOffset * n)
	ZeroValidatorPubkeyGIndexDenebBlock = 3254554418216960

	// ValidatorPubkeyGIndexOffset is the offset of a validator pubkey GIndex.
	ValidatorPubkeyGIndexOffset = 8
)
