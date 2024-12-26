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
//

package cometbft

import (
	"fmt"

	cmtcfg "github.com/cometbft/cometbft/config"
	cmttypes "github.com/cometbft/cometbft/types"
)

// DefaultConsensusParams returns the default consensus parameters
// shared by every node in the network. Consensus parameters are
// inscripted in genesis.
func DefaultConsensusParams(consensusKeyAlgo string) *cmttypes.ConsensusParams {
	res := cmttypes.DefaultConsensusParams()
	res.Validator.PubKeyTypes = []string{consensusKeyAlgo}

	if err := res.ValidateBasic(); err != nil {
		panic(fmt.Errorf("invalid default consensus parameters: %w", err))
	}

	return res
}

func extractConsensusParams(cmtCfg *cmtcfg.Config) (*cmttypes.ConsensusParams, error) {
	// Consensus parameters are immutable (do not change as slots go by).
	// So we reuse the parameters specified in genesis.
	// Todo: add validation for genesis params by chainID

	genFunc := GetGenDocProvider(cmtCfg)
	genDoc, err := genFunc()
	if err != nil {
		return nil, err
	}

	cmtConsensusParams := genDoc.GenesisDoc.ConsensusParams
	return cmtConsensusParams, nil
}
