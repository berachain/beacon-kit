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

package vm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/block"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/encoding"
)

//nolint:gochecknoglobals // seems clearer to me
var genesisTime = time.Date(2024, 8, 21, 10, 17, 0, 0, time.UTC)

// Genesis is the in-memory representation of genesis.
type Genesis struct {
	Validators []*Validator
	EthGenesis []byte
}

func DefaultEthGenesisBytes() ([]byte, error) {
	var (
		gen = make(map[string]json.RawMessage)
		err error
	)
	gen["beacon"], err = json.Marshal(types.DefaultGenesisDeneb())
	if err != nil {
		return nil, err
	}
	return json.Marshal(gen)
}

// process genesisBytes and from them build:
// genesis block, the first block in the chain
// genesis validators, the validators initially responsible for the chain
// ethGenesis bytes, to be passed to the middleware.
func ParseGenesis(
	genesisBytes []byte,
) (*block.StatelessBlock, []*Validator, []byte, error) {
	gen, err := parseInMemoryGenesis(genesisBytes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed parsing genesis: %w", err)
	}

	genBlk, err := buildGenesisBlock(genesisBytes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed building genesis block: %w", err)
	}

	return genBlk, gen.Validators, gen.EthGenesis, nil
}

// parseInMemoryGenesis is used in VM initialization
// to retrieve genesis from its bytes.
func parseInMemoryGenesis(genesisBytes []byte) (*Genesis, error) {
	inMemGen := &Genesis{}
	if err := encoding.Decode(genesisBytes, inMemGen); err != nil {
		return nil, fmt.Errorf("unable to gob decode genesis content: %w", err)
	}

	// make sure to calculate ID of every validator
	for i, v := range inMemGen.Validators {
		if err := v.initValID(); err != nil {
			return nil, fmt.Errorf("validator pos %d: %w", i, err)
		}
	}

	return inMemGen, nil
}

// build a block from genesis content and keep it as
// first block in the chain.
func buildGenesisBlock(genesisBytes []byte) (*block.StatelessBlock, error) {
	// Genesis block must be parsable as a block, but
	// genesis bytes do no encode a block. We create genesis block
	// by using genesis bytes as block content so that
	// genesis block ID depends on genesisBytes
	return block.NewStatelessBlock(
		ids.Empty,
		0,
		genesisTime,
		block.Content{GenesisContent: genesisBytes},
	)
}
