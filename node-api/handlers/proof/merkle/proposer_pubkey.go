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

import (
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

// ProveProposerPubkeyInBlock generates a proof for the proposer pubkey in the
// beacon block.
func ProveProposerPubkeyInBlock(
	bbh *ctypes.BeaconBlockHeader,
	bsm types.BeaconStateMarshallable,
) ([]common.Root, common.Root, error) {
	return ProveValidatorPubkeyInBlock(bbh.GetProposerIndex(), bbh, bsm)
}

// ProveProposerPubkeyInState generates a proof for the proposer pubkey
// in the beacon state.
func ProveProposerPubkeyInState(
	forkVersion common.Version,
	bsm types.BeaconStateMarshallable,
	proposerOffset math.U64,
) ([]common.Root, common.Root, error) {
	return ProveValidatorPubkeyInState(forkVersion, bsm, proposerOffset)
}
