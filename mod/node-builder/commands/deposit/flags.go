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
)

const (
	// broadcastDepositShorthand is the shorthand flag for the broadcastDeposit
	// flag.
	broadcastDepositShorthand = "b"
)

const (
	// defaultBroadcastDeposit is the default value for the broadcastDeposit
	// flag.
	defaultBroadcastDeposit = false

	// defaultPrivateKey is the default value for the privateKey flag.
	defaultPrivateKey = ""
)

const (
	// broadcastDepositFlagUsage is the usage description for the
	// broadcastDeposit flag.
	broadcastDepositMsg = "broadcast the deposit transaction"

	// privateKeyFlagUsage is the usage description for the privateKey flag.
	privateKeyMsg = "private key to sign the deposit message"
)
