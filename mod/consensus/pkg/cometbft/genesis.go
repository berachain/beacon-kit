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
	"time"

	cmttypes "github.com/cometbft/cometbft/types"
)

type Genesis struct {
	*cmttypes.GenesisDoc `mapstructure:",squash"`
}

func NewGenesis(
	chainID string,
	appState []byte,
	consensusParams *cmttypes.ConsensusParams,
) *Genesis {
	return &Genesis{
		GenesisDoc: &cmttypes.GenesisDoc{
			GenesisTime:     time.Now(),
			ChainID:         chainID,
			InitialHeight:   1,
			ConsensusParams: consensusParams,
			AppState:        appState,
		},
	}
}

// func (g *Genesis) UnmarshalJSON(data []byte) error {
// 	return cmtjson.Unmarshal(data, g.GenesisDoc)
// }

func (g *Genesis) Export(path string) error {
	if err := g.ValidateAndComplete(); err != nil {
		return err
	}
	return g.SaveAs(path)
}

func (g *Genesis) ToGenesisDoc() (*cmttypes.GenesisDoc, error) {
	return g.GenesisDoc, nil
}
