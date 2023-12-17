// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package eth

import (
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/itsdevbear/bolaris/beacon/log"
)

// loadJWTSecret reads the JWT secret from a file and returns it.
// It returns an error if the file cannot be read or if the JWT secret is not valid.
func LoadJWTSecret(filepath string, logger log.Logger) ([]byte, error) {
	// Read the file.
	data, err := os.ReadFile(filepath)
	if err != nil {
		// Return an error if the file cannot be read.
		return nil, err
	}

	// Convert the data to a JWT secret.
	jwtSecret := common.FromHex(strings.TrimSpace(string(data)))

	// Check if the JWT secret is valid.
	if len(jwtSecret) != jwtLength {
		// Return an error if the JWT secret is not valid.
		return nil, fmt.Errorf("failed to load jwt secret from %s", filepath)
	}

	logger.Info("loaded exeuction client jwt secret file", "path", filepath, "crc32")
	return jwtSecret, nil
}
