#!/bin/bash
# SPDX-License-Identifier: MIT
#
# Copyright (c) 2024 Berachain Foundation
#
# Permission is hereby granted, free of charge, to any person
# obtaining a copy of this software and associated documentation
# files (the "Software"), to deal in the Software without
# restriction, including without limitation the rights to use,
# copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following
# conditions:
#
# The above copyright notice and this permission notice shall be
# included in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
# EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
# OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
# NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
# HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
# WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
# FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.

set -x

GENESIS=$BEACOND_HOME/config/genesis.json
TMP_GENESIS=$BEACOND_HOME/config/tmp_genesis.json

# Create beacond config directory
if [ ! -d "$BEACOND_HOME/config" ]; then
    # # Init the chain
    /usr/bin/beacond init --chain-id "$BEACOND_CHAIN_ID" "$BEACOND_MONIKER" --home "$BEACOND_HOME" --beacon-kit.accept-tos
	
	# Set client config
	/usr/bin/beacond config set client keyring-backend $BEACOND_KEYRING_BACKEND --home "$BEACOND_HOME"
	/usr/bin/beacond config set client chain-id "$BEACOND_CHAIN_ID" --home "$BEACOND_HOME"

	# If keys exist they should be deleted
	/usr/bin/beacond keys add "$BEACOND_MONIKER" --keyring-backend $BEACOND_KEYRING_BACKEND --home "$BEACOND_HOME" --indiscreet --output json > "$BEACOND_HOME/config/mnemonic.json"

	# Change parameter token denominations to abgt
	jq '.consensus["params"]["block"]["max_gas"]="30000000"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq '.consensus.params.validator.pub_key_types += ["bls12_381"] | .consensus.params.validator.pub_key_types -= ["ed25519"]' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

	# Add pubkey to the genesis file.
	/usr/bin/beacond genesis add-validator --home "$BEACOND_HOME"
fi
