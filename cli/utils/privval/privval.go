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

// Package privval initializes the CometBFT node key and consensus
// validator files with BLS key support. The cosmos-sdk genutil helper
// (InitializeNodeValidatorFilesFromMnemonicWithKeyType) is not used
// because its BLS path overwrites an existing priv_validator_key.json.
package privval

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/berachain/beacon-kit/primitives/crypto"
	cmtcfg "github.com/cometbft/cometbft/config"
	cmtcrypto "github.com/cometbft/cometbft/crypto"
	cmtbls12381 "github.com/cometbft/cometbft/crypto/bls12381"
	cmted25519 "github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/cosmos/go-bip39"
)

// InitializeNodeValidatorFiles creates the node key and private validator
// files if they do not exist, and returns the node ID and validator pubkey.
func InitializeNodeValidatorFiles(
	config *cmtcfg.Config, keyType string,
) (string, cmtcrypto.PubKey, error) {
	return InitializeNodeValidatorFilesFromMnemonic(config, "", keyType)
}

// InitializeNodeValidatorFilesFromMnemonic is like
// InitializeNodeValidatorFiles but derives the ed25519 consensus key from
// the mnemonic when one is provided. BLS keys do not support mnemonics.
func InitializeNodeValidatorFilesFromMnemonic(
	config *cmtcfg.Config, mnemonic, keyType string,
) (string, cmtcrypto.PubKey, error) {
	if len(mnemonic) > 0 && !bip39.IsMnemonicValid(mnemonic) {
		return "", nil, errors.New("invalid mnemonic")
	}
	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return "", nil, err
	}
	nodeID := string(nodeKey.ID())

	pvKeyFile := config.PrivValidatorKeyFile()
	if err = os.MkdirAll(filepath.Dir(pvKeyFile), 0o750); err != nil {
		return "", nil, fmt.Errorf(
			"could not create directory %q: %w", filepath.Dir(pvKeyFile), err,
		)
	}
	pvStateFile := config.PrivValidatorStateFile()
	if err = os.MkdirAll(filepath.Dir(pvStateFile), 0o750); err != nil {
		return "", nil, fmt.Errorf(
			"could not create directory %q: %w", filepath.Dir(pvStateFile), err,
		)
	}

	privKey, err := genPrivKey(mnemonic, keyType)
	if err != nil {
		return "", nil, err
	}
	filePV := loadOrGenFilePV(privKey, pvKeyFile, pvStateFile)

	pubKey, err := filePV.GetPubKey()
	if err != nil {
		return "", nil, err
	}
	return nodeID, pubKey, nil
}

func genPrivKey(mnemonic, keyType string) (cmtcrypto.PrivKey, error) {
	switch keyType {
	case crypto.CometBLSType:
		if len(mnemonic) > 0 {
			return nil, errors.New("BLS key type does not support mnemonic")
		}
		return cmtbls12381.GenPrivKey()
	case "", cmted25519.KeyType:
		if len(mnemonic) > 0 {
			return cmted25519.GenPrivKeyFromSecret([]byte(mnemonic)), nil
		}
		return cmted25519.GenPrivKey(), nil
	default:
		return nil, fmt.Errorf("unsupported consensus key type %q", keyType)
	}
}

// loadOrGenFilePV loads a FilePV from the given file paths or generates a
// new one from privKey and saves it.
func loadOrGenFilePV(
	privKey cmtcrypto.PrivKey, keyFilePath, stateFilePath string,
) *privval.FilePV {
	if _, err := os.Stat(keyFilePath); err == nil {
		return privval.LoadFilePV(keyFilePath, stateFilePath)
	}
	pv := privval.NewFilePV(privKey, keyFilePath, stateFilePath)
	pv.Save()
	return pv
}
