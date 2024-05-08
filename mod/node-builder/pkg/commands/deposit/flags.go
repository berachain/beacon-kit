// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package deposit

const (
	// broadcastDeposit is the flag for broadcasting a deposit transaction.
	broadcastDeposit = "broadcast"

	// privateKey is the flag for the private key to sign the deposit message.
	privateKey = "private-key"

	// overrideNodeKey is the flag for overriding the node key.
	overrideNodeKey = "override-node-key"

	// validatorPrivateKey is the flag for the validator private key.
	valPrivateKey = "validator-private-key"

	// jwtSecretPath is the flag for the path to the JWT secret file.
	jwtSecretPath = "jwt-secret"

	// engineRPCURL is the flag for the URL for the engine RPC.
	engineRPCURL = "engine-rpc-url"
)

const (
	// broadcastDepositShorthand is the shorthand flag for the broadcastDeposit
	// flag.
	broadcastDepositShorthand = "b"

	// overrideNodeKeyShorthand is the shorthand flag for the overrideNodeKey
	// flag.
	overrideNodeKeyShorthand = "o"
)

const (
	// defaultBroadcastDeposit is the default value for the broadcastDeposit
	// flag.
	defaultBroadcastDeposit = false

	// defaultPrivateKey is the default value for the privateKey flag.
	defaultPrivateKey = ""

	// defaultOverrideNodeKey is the default value for the overrideNodeKey flag.
	defaultOverrideNodeKey = false

	// defaultValidatorPrivateKey is the default value for the
	// validatorPrivateKey flag.
	defaultValidatorPrivateKey = ""

	// defaultJWTSecretPath is the default value for the jwtSecret flag.
	// #nosec G101 // This is a default path
	defaultJWTSecretPath = "../jwt.hex"

	// defaultEngineRPCURL is the default value for the engineRPCURL flag.
	defaultEngineRPCURL = "http://localhost:8551"
)

const (
	// broadcastDepositFlagMsg is the usage description for the
	// broadcastDeposit flag.
	broadcastDepositMsg = "broadcast the deposit transaction"

	// privateKeyFlagMsg is the usage description for the privateKey flag.
	privateKeyMsg = `private key to sign and pay for the deposit message. 
	This is required if the broadcast flag is set.`

	// overrideNodeKeyFlagMsg is the usage description for the overrideNodeKey
	// flag.
	overrideNodeKeyMsg = "override the node private key"

	// valPrivateKeyMsg is the usage description for the
	// valPrivateKey flag.
	valPrivateKeyMsg = `validator private key. This is required if the 
	override-node-key flag is set.`

	// jwtSecretPathMsg is the usage description for the jwtSecretPath flag.
	// #nosec G101 // This is a descriptor
	jwtSecretPathMsg = "path to the JWT secret file"

	// engineRPCURLMsg is the usage description for the engineRPCURL flag.
	engineRPCURLMsg = "URL for the engine RPC"
)
