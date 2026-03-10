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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package preconf

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/berachain/beacon-kit/cli/utils/parser"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/net/jwt"
)

// LoadWhitelist loads validator pubkeys from a JSON file.
func LoadWhitelist(path string) ([]crypto.BLSPubkey, error) {
	if path == "" {
		return nil, errors.New("whitelist path is empty")
	}

	data, err := os.ReadFile(path) // #nosec G304 -- path from operator-controlled config
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read whitelist file: %s", path)
	}

	var pubkeyStrings []string
	if err = json.Unmarshal(data, &pubkeyStrings); err != nil {
		return nil, errors.Wrapf(err, "failed to parse whitelist JSON from: %s", path)
	}

	pubkeys := make([]crypto.BLSPubkey, 0, len(pubkeyStrings))
	for i, s := range pubkeyStrings {
		pk, convErr := parser.ConvertPubkey(s)
		if convErr != nil {
			return nil, errors.Wrapf(convErr, "invalid pubkey at index %d in whitelist", i)
		}
		pubkeys = append(pubkeys, pk)
	}

	return pubkeys, nil
}

// ValidatorJWTs maps validator pubkeys to their JWT secrets.
type ValidatorJWTs map[crypto.BLSPubkey]*jwt.Secret

// LoadValidatorJWTs loads validator pubkeys to JWT secret mappings from a JSON file.
// The JSON file format is: {"0xpubkey1": "jwtsecret1hex", "0xpubkey2": "jwtsecret2hex", ...}
func LoadValidatorJWTs(path string) (ValidatorJWTs, error) {
	if path == "" {
		return nil, errors.New("validator JWTs path is empty")
	}

	data, err := os.ReadFile(path) // #nosec G304 -- path from operator-controlled config
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read validator JWTs file: %s", path)
	}

	var rawMap map[string]string
	if err = json.Unmarshal(data, &rawMap); err != nil {
		return nil, errors.Wrapf(err, "failed to parse validator JWTs JSON from: %s", path)
	}

	result := make(ValidatorJWTs, len(rawMap))
	for pubkeyStr, jwtSecretHex := range rawMap {
		// Parse pubkey
		pubkey, convErr := parser.ConvertPubkey(pubkeyStr)
		if convErr != nil {
			return nil, errors.Wrapf(convErr, "invalid pubkey in validator JWTs: %s", pubkeyStr)
		}

		// Parse JWT secret
		jwtSecret, jwtErr := jwt.NewFromHex(jwtSecretHex)
		if jwtErr != nil {
			return nil, errors.Wrapf(jwtErr, "invalid JWT secret for pubkey: %s", pubkeyStr)
		}

		result[pubkey] = jwtSecret
	}

	return result, nil
}

// LoadJWTSecret loads a JWT secret from a file containing a hex-encoded secret.
func LoadJWTSecret(path string) (*jwt.Secret, error) {
	if path == "" {
		return nil, errors.New("JWT secret path is empty")
	}

	data, err := os.ReadFile(path) // #nosec G304 -- path from operator-controlled config
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read JWT secret file: %s", path)
	}

	// Trim whitespace and newlines
	hexStr := strings.TrimSpace(string(data))

	secret, err := jwt.NewFromHex(hexStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse JWT secret from: %s", path)
	}

	return secret, nil
}
