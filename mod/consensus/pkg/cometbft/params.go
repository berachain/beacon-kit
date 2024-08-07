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

package cometbft

import (
	"encoding/json"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	v1 "github.com/cometbft/cometbft/api/cometbft/types/v1"
	cmttypes "github.com/cometbft/cometbft/types"
)

type ConsensusParams struct {
	*cmttypes.ConsensusParams
}

func DefaultConsensusParams() *ConsensusParams {
	params := cmttypes.DefaultConsensusParams()
	params.Validator.PubKeyTypes = []string{crypto.CometBLSType}
	return &ConsensusParams{
		ConsensusParams: params,
	}
}

func (cp *ConsensusParams) Default() *ConsensusParams {
	return DefaultConsensusParams()
}

func (cp *ConsensusParams) MarshalJSON() ([]byte, error) {
	return json.Marshal(cp.ConsensusParams)
}

type ConsensusParamsStore struct {
	chainSpec common.ChainSpec
}

func NewConsensusParamsStore(chainSpec common.ChainSpec) *ConsensusParamsStore {
	return &ConsensusParamsStore{
		chainSpec: chainSpec,
	}
}

func (cps *ConsensusParamsStore) Get(slot uint64) (*v1.ConsensusParams, error) {
	params, ok := cps.chainSpec.GetCometBFTConfigForSlot(math.U64(slot)).(*ConsensusParams)
	if !ok {
		return nil, errors.Newf(
			"failed to get consensus params for slot %d",
			slot,
		)
	}
	p := params.ConsensusParams.ToProto()
	return &p, nil
}
