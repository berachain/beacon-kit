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

package gethprimitives

import (
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
)

type (
	// ExecutionAddress represents an address on the execution layer
	// which is derived via secp256k1 w/recovery bit.
	//
	// Related: https://eips.ethereum.org/EIPS/eip-55
	ExecutionAddress = common.Address
	// ExecutionHash represents a hash on the execution layer which is
	// currently a Keccak256 hash.
	ExecutionHash = common.Hash
	// DisplayBytes is an alias for common.PrettyBytes.
	DisplayBytes   = common.PrettyBytes
	ExecutableData = engine.ExecutableData
	Genesis        = core.Genesis
	Block          = coretypes.Block
	Body           = coretypes.Body
	Log            = coretypes.Log
	Header         = coretypes.Header
	Receipt        = coretypes.Receipt
	Transaction    = coretypes.Transaction
	Transactions   = coretypes.Transactions
	Withdrawal     = coretypes.Withdrawal
	Withdrawals    = coretypes.Withdrawals
)

//nolint:gochecknoglobals // alias.
var (
	HexToAddress           = common.HexToAddress
	HexToHash              = common.HexToHash
	Hex2BytesFixed         = common.Hex2BytesFixed
	FromHex                = common.FromHex
	ZeroAddress            = ExecutionAddress{}
	ZeroHash               = ExecutionHash{}
	BlockToExecutableData  = engine.BlockToExecutableData
	NewBlockWithHeader     = coretypes.NewBlockWithHeader
	DeriveSha              = coretypes.DeriveSha
	EmptyUncleHash         = coretypes.EmptyUncleHash
	NewStackTrie           = trie.NewStackTrie
	BytesToBloom           = coretypes.BytesToBloom
	SignTx                 = coretypes.SignTx
	LatestSignerForChainID = coretypes.LatestSignerForChainID
	ReceiptStatusFailed    = coretypes.ReceiptStatusFailed
)
