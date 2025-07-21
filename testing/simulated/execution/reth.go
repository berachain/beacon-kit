//go:build simulated

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

package execution

import (
	"fmt"

	"github.com/ory/dockertest/docker"
)

// NewRethNode creates a new execution node configured for Reth by using the default Reth command builder.
func NewRethNode(homeDir string, image docker.PullImageOptions) *ExecNode {
	return NewExecNode(homeDir, image, defaultRethCmdStrBuilder)
}

// ValidRethImage returns the default Docker image options for the Reth node.
func ValidRethImage() docker.PullImageOptions {
	return docker.PullImageOptions{
		Repository: "ghcr.io/berachain/bera-reth",
		Tag:        "v1.0.0-rc.1",
	}
}

// defaultRethCmdStrBuilder returns a command string tailored for running a Geth node.
func defaultRethCmdStrBuilder(genesisFile string) string {
	return fmt.Sprintf(`
		bera-reth node --http --http.addr 0.0.0.0 --http.api eth,net,web3,debug \
			 --chain=/testdata/%s \
			 --authrpc.addr 0.0.0.0 \
			 --authrpc.jwtsecret /testing/files/jwt.hex \
			 --datadir /tmp/rethdata \
			 --full \
			 --engine.persistence-threshold=0 \
			 --engine.memory-block-buffer-target=0 \
			 -vvvv \
	`, genesisFile)
}
