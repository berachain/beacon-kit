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
	"log"
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/stretchr/testify/require"
)

var (
	testGenesisValidators []*Validator
	testEthGenesisBytes   []byte
)

func init() {
	// init testEthGenesisBytes
	var err error
	testEthGenesisBytes, err = DefaultEthGenesisBytes()
	if err != nil {
		log.Fatal(err)
	}

	// init testValidators
	val0, err := NewValidator(ids.GenerateTestNodeID(), uint64(999))
	if err != nil {
		log.Fatal(err)
	}
	val1, err := NewValidator(ids.GenerateTestNodeID(), uint64(1001))
	if err != nil {
		log.Fatal(err)
	}
	testGenesisValidators = []*Validator{val0, val1}
}

func TestEthGenesisEncoding(t *testing.T) {
	r := require.New(t)

	// setup genesis
	genesisData := &Base64Genesis{
		Validators: []Base64GenesisValidator{
			{
				NodeID: testGenesisValidators[0].NodeID.String(),
				Weight: testGenesisValidators[0].Weight,
			},
			{
				NodeID: testGenesisValidators[1].NodeID.String(),
				Weight: testGenesisValidators[1].Weight,
			},
		},
		EthGenesis: string(testEthGenesisBytes),
	}

	// marshal genesis
	genContent, err := BuildBase64GenesisString(genesisData)
	r.NoError(err)

	// unmarshal genesis
	parsedGenesisData, err := ParseBase64StringToBytes(genContent)
	r.NoError(err)

	_, rValidators, rGenEthData, err := parseGenesis(parsedGenesisData)
	r.NoError(err)
	r.Equal(testGenesisValidators, rValidators)
	r.Equal(testEthGenesisBytes, rGenEthData)
}
